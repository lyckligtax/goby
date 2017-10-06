package main

import (
	"encoding/json"
	"errors"
	"strings"
	"strconv"
	"time"
)

type Ctrl struct {
	// Mandatory Fields to build a valid control file
	Version      string `json:"version"`
	Package      string `json:"package"`
	Architecture string `json:"architecture"`
	Maintainer   string `json:"maintainer"`
	Synopsis     string `json:"synopsis"`
	
	// Optional Fields to build a valid control file
	Dependencies map[string][]string `json:"dependencies"`
	Homepage     string              `json:"homepage"`
	Description  string              `json:"description"`
	Essential    string              `json:"essential"`
	Section      string              `json:"section"`
	Priority     string              `json:"priority"`
	
	// Optional Scriptfiles
	PreInst  string `json:"preInst"`
	PostInst string `json:"postInst"`
	PreRm    string `json:"preRm"`
	PostRm   string `json:"postRm"`
	Config   string `json:"config"`
	
	// Optional List of configurations and libraries
	ConfFiles []string `json:"confFiles"`
	ShLibs    []string `json:"shLibs"`
	
	// Package Files
	Files []*File
	
	// Checks
	md5sBuild bool
	err       error
}

func NewCtrl(config []byte) (*Ctrl, error) {
	ctrl := &Ctrl{}
	err := json.Unmarshal(config, ctrl)
	
	if err != nil {
		return nil, err
	}
	return ctrl, nil
}

func (bc *Ctrl) Check() error {
	if bc.err != nil {
		return bc.err
	}
	
	if err := bc.checkMandatoryFields(); err != nil {
		bc.err = err
	}
	
	if err := bc.checkScriptFiles(); err != nil {
		bc.err = err
	}
	
	return bc.err
}

// TODO: Implement whitespace checking
func (bc *Ctrl) checkMandatoryFields() error {
	if bc.err != nil {
		return bc.err
	}
	
	// TODO: Implement stricter checking: [a-z0-9][a-z0-9\.+-]+
	if bc.Package == "" {
		return errors.New("invalid package name")
	}
	if bc.Version == "" {
		return errors.New("invalid version")
	}
	// TODO: Implement stricter checking: single architecture only
	if bc.Architecture == "" {
		return errors.New("invalid architecture")
	}
	// TODO: Implement stricter checking: "fistname surname <email>" (in RFC822 format)
	if bc.Maintainer == "" {
		return errors.New("invalid maintainer")
	}
	// TODO: Implement stricter checking: has to be single line
	if bc.Synopsis == "" {
		return errors.New("invalid synopsis")
	}
	
	return nil
}

func (bc *Ctrl) checkScriptFiles() error {
	if bc.err != nil {
		return bc.err
	}
	
	files := map[string]string{
		"preinst":  bc.PreInst,
		"postinst": bc.PreInst,
		"prerm":    bc.PreRm,
		"postrm":   bc.PostRm,
		"config":   bc.Config,
	}
	
	for dest, src := range files {
		if src == "" {
			continue
		}
		
		f, err := NewFile(src, dest)
		if err != nil {
			return err
		}
		bc.Files = append(bc.Files, f)
	}
	
	return nil
}

func (bc *Ctrl) Build(data *Data) ([]byte, error) {
	if bc.err != nil {
		return nil, bc.err
	}
	
	if err := bc.buildControlFile(data); err != nil {
		bc.err = err
	}
	
	if err := bc.buildMD5SumsFile(data); err != nil {
		bc.err = err
	}
	
	if err := bc.buildConfFilesFile(); err != nil {
		bc.err = err
	}
	
	if err := bc.buildShLibsFile(); err != nil {
		bc.err = err
	}
	
	if bc.err != nil {
		return nil, bc.err
	}
	
	// create tar archive
	ctrlTar := NewTarArchive()
	for _, f := range bc.Files {
		if f == nil {
			continue
		}
		ctrlTar.AddFile(f.Destination, f.Content)
	}
	
	gzippedTar, err := ctrlTar.Gzip()
	if err != nil {
		bc.err = err
		
	}
	return gzippedTar, bc.err
}

func (bc *Ctrl) buildControlFile(data *Data) error {
	if bc.err != nil {
		return bc.err
	}
	
	content := []string{
		"Package: " + bc.Package,
		"Version: " + bc.Version,
		"Architecture: " + bc.Architecture,
		"Maintainer: " + bc.Maintainer,
		"Installed-Size: " + strconv.Itoa(data.DataSizeKB()),
		"Date: " + time.Now().UTC().Format(time.UnixDate),
		"Description: " + bc.Synopsis,
	}
	
	if d := whiteTrim(bc.Description); d != "" {
		for _, line := range strings.Split(d, "\n") {
			content = append(content, " "+line)
		}
	}
	
	if h := whiteTrim(bc.Homepage); h != "" {
		content = append(content, "Homepage: "+h)
	}
	
	if e := whiteTrim(bc.Essential); e == "yes" {
		content = append(content, "Essential: yes")
	}
	
	if s := whiteTrim(bc.Section); s != "" {
		content = append(content, "Section: "+s)
	}
	
	if p := whiteTrim(bc.Priority); p == "required" || p == "important" || p == "standard" || p == "optional" {
		content = append(content, "Priority: "+p)
	}
	
	if deps, ok := bc.Dependencies["depends"]; ok && len(deps) > 0 {
		if pkgs := multiWhiteTrim(deps); len(pkgs) > 0 {
			content = append(content, "Depends: "+strings.Join(pkgs, ", "))
		}
	}
	
	if deps, ok := bc.Dependencies["recommends"]; ok && len(deps) > 0 {
		if pkgs := multiWhiteTrim(deps); len(pkgs) > 0 {
			content = append(content, "Recommends: "+strings.Join(pkgs, ", "))
		}
	}
	
	if deps, ok := bc.Dependencies["suggests"]; ok && len(deps) > 0 {
		if pkgs := multiWhiteTrim(deps); len(pkgs) > 0 {
			content = append(content, "Suggests: "+strings.Join(pkgs, ", "))
		}
	}
	
	if deps, ok := bc.Dependencies["enhances"]; ok && len(deps) > 0 {
		if pkgs := multiWhiteTrim(deps); len(pkgs) > 0 {
			content = append(content, "Enhances: "+strings.Join(pkgs, ", "))
		}
	}
	
	if deps, ok := bc.Dependencies["preDepends"]; ok && len(deps) > 0 {
		if pkgs := multiWhiteTrim(deps); len(pkgs) > 0 {
			content = append(content, "Pre-Depends: "+strings.Join(pkgs, ", "))
		}
	}
	
	if deps, ok := bc.Dependencies["breaks"]; ok && len(deps) > 0 {
		if pkgs := multiWhiteTrim(deps); len(pkgs) > 0 {
			content = append(content, "Breaks: "+strings.Join(pkgs, ", "))
		}
	}
	
	if deps, ok := bc.Dependencies["conflicts"]; ok && len(deps) > 0 {
		if pkgs := multiWhiteTrim(deps); len(pkgs) > 0 {
			content = append(content, "Conflicts: "+strings.Join(pkgs, ", "))
		}
	}
	
	f := &File{
		Destination: "control",
		Content:     []byte(strings.Join(content, "\n")),
	}
	
	bc.Files = append(bc.Files, f)
	return nil
}

func (bc *Ctrl) buildConfFilesFile() error {
	if bc.err != nil {
		return bc.err
	}
	
	if len(bc.ConfFiles) <= 0 {
		return nil
	}
	
	files := []string{}
	for _, file := range bc.ConfFiles {
		if file != "" {
			files = append(files, string(file))
		}
	}
	
	if len(files) <= 0 {
		return nil
	}
	
	f := &File{
		Destination: "conffiles",
		Content:     []byte(strings.Join(files, "\n")),
	}
	bc.Files = append(bc.Files, f)
	return nil
}

func (bc *Ctrl) buildShLibsFile() error {
	if bc.err != nil {
		return bc.err
	}
	
	if len(bc.ShLibs) <= 0 {
		return nil
	}
	
	files := []string{}
	for _, file := range bc.ShLibs {
		if file != "" {
			files = append(files, file)
		}
	}
	
	if len(files) <= 0 {
		return nil
	}
	
	f := &File{
		Destination: "shlibs",
		Content:     []byte(strings.Join(files, "\n")),
	}
	bc.Files = append(bc.Files, f)
	return nil
}

func (bc *Ctrl) buildMD5SumsFile(data *Data) error {
	if bc.err != nil {
		return bc.err
	}
	
	sums := make([]string, len(data.Files))
	for i, file := range data.Files {
		sums[i] = file.MD5 + " " + strings.TrimLeft(file.Destination, "/")
	}
	
	f := &File{
		Destination: "md5sums",
		Content:     []byte(strings.Join(sums, "\n")),
	}
	bc.Files = append(bc.Files, f)
	return nil
}
