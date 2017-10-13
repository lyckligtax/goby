package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"errors"
)

// Goby resembles the content of a deb-package as found on Ubuntu or Debian derivates
// it consists of gzipped control.tar.gz and data.tar.gz tar archives
// and a debian-binary file containing the package semantic version
// all three files are packaged into an ar-Archive
type Goby struct {
	ctrl *Ctrl
	data *Data

	preGoby    string
	postGoby   string
	outputFile string

	err error
}

func NewGoby(cfgFilepath string) (*Goby, error) {
	config, err := ioutil.ReadFile(cfgFilepath)
	if err != nil {
		return nil, err
	}

	ctrl, err := NewCtrl(config)
	if err != nil {
		return nil, err
	}

	data, err := NewData(config)
	if err != nil {
		return nil, err
	}

	actions := struct {
		PreGoby    string `json:"preGoby"`
		PostGoby   string `json:"postGoby"`
		OutputFile string `json:"outputFile"`
	}{}

	err = json.Unmarshal(config, &actions)
	if err != nil {
		return nil, err
	}

	return &Goby{
		ctrl:       ctrl,
		data:       data,
		preGoby:    actions.PreGoby,
		postGoby:   actions.PostGoby,
		outputFile: actions.OutputFile,
	}, nil
}

func (g *Goby) Make(filePath string) error {
	if g.err != nil {
		return g.err
	}

	fmt.Print("Running pre goby commands ... ")
	if err := g.run(g.preGoby); err != nil {
		g.err = err
		return err
	}
	fmt.Print("done\n")

	/*
		Checks
	*/
	fmt.Print("Running package checks ... ")
	if err := g.check(); err != nil {
		g.err = err
		return err
	}
	fmt.Print("done\n")

	/*
		Builds
	*/
	fmt.Print("Building package ... ")
	content, err := g.build()
	if err != nil {
		g.err = err
		return err
	}
	fmt.Print("done\n")

	/*
		Writing Data
	*/
	if filePath == defaultPackageName {
		if g.outputFile != "" {
			filePath = g.outputFile
		} else {
			filePath = string(g.ctrl.Package) + "-" + string(g.ctrl.Version) + ".deb"
		}
	}

	fmt.Print("Writing contents to " + filePath + " ... ")
	if err := g.write(content, filePath); err != nil {
		g.err = err
		return err
	}
	fmt.Print("done\n")

	fmt.Print("Running post goby commands ... ")
	if err := g.run(g.postGoby); err != nil {
		g.err = err
		return err
	}
	fmt.Print("done\n")

	return nil
}

func (g *Goby) run(command string) error {
	if g.err != nil {
		return g.err
	}

	command = whiteTrim(command)
	if command == "" {
		return nil
	}

	
	shell, err := exec.LookPath("bash")
	if err != nil {
		shell, err = exec.LookPath("sh")
		if err != nil {
			return errors.New("Neither bash nor sh found to execute: " + command)
		}
	}
	
	cmd := exec.Command(shell,"-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (g *Goby) check() error {
	if g.err != nil {
		return g.err
	}

	if err := g.data.Check(); err != nil {
		g.err = err
		return err
	}

	if err := g.ctrl.Check(); err != nil {
		g.err = err
		return err
	}

	return nil
}

func (g *Goby) build() ([]byte, error) {
	if g.err != nil {
		return nil, g.err
	}

	deb := NewArArchive()

	// Ubuntu is picky about debian-binary being first in the package
	deb.AddFile("debian-binary", []byte("2.0\n"))
	
	// Build gzipped control tar archive
	ctrl, err := g.ctrl.Build(g.data)
	if err != nil {
		g.err = err
		return nil, err
	}
	
	// Add control archive to ar-Archive
	if err = deb.AddFile("control.tar.gz", ctrl); err != nil {
		g.err = err
		return nil, err
	}
	
	// Build gzipped data tar archive
	data, err := g.data.Build()
	if err != nil {
		g.err = err
		return nil, err
	}

	// Add data archive to ar-Archive
	if err = deb.AddFile("data.tar.gz", data); err != nil {
		g.err = err
		return nil, err
	}

	return deb.Data()
}

func (g *Goby) write(content []byte, filePath string) error {
	if g.err != nil {
		return g.err
	}

	f, err := os.Create(filePath)
	defer f.Close()
	if err != nil {
		g.err = err
		return err
	}

	_, err = f.Write(content)
	if err != nil {
		g.err = err
		return err
	}

	return nil
}
