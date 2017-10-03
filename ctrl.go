package main

import (
	"errors"
	"encoding/json"
	"fmt"
	"strings"
)

type BinCtrl struct {
	// Mandatory Fields to build a valid control file
	Version      NEString `json:"version"`
	Package      NEString `json:"package"`
	Architecture NEString `json:"architecture"`
	Maintainer   NEString `json:"maintainer"`
	Description  NEString `json:"description"`

	// Optional Scriptfiles
	PreInst  NEString `json:"preInst"`
	PostInst NEString `json:"postInst"`
	PreRm    NEString `json:"preRm"`
	PostRm   NEString `json:"postRm"`
	Config   NEString `json:"config"`

	// Optional List of configurations and libraries
	ConfFiles []NEString `json:"confFiles"`
	ShLibs    []NEString `json:"shLibs"`

	// Package Files
	Files []*File

	// Checks
	md5sBuild bool
	err   error
}

func NewBinCtrl(config []byte) (*BinCtrl, error) {
	ctrl := &BinCtrl{}
	err := json.Unmarshal(config, ctrl)

	if err != nil {
		return nil, err
	}
	return ctrl, nil
}

func (bc *BinCtrl) Check() error {
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

func (bc *BinCtrl) checkMandatoryFields() error {
	if bc.err != nil {
		return bc.err
	}

	if bc.Package == "" {
		return errors.New("invalid package name")
	}
	if !isSemanticVersionRxp.MatchString(string(bc.Version)) {
		return errors.New("invalid version")
	}
	if bc.Architecture == "" {
		return errors.New("invalid architecture")
	}
	if bc.Maintainer == "" {
		return errors.New("invalid maintainer")
	}
	if bc.Description == "" {
		return errors.New("invalid description")
	}

	return nil
}

func (bc *BinCtrl) checkScriptFiles() error {
	if bc.err != nil {
		return bc.err
	}

	files := map[string]NEString{
		"peinst":   bc.PreInst,
		"postinst": bc.PreInst,
		"prerm":    bc.PreRm,
		"postrm":   bc.PostRm,
		"config":   bc.Config,
	}

	var f *File
	var err error

	for dest, src := range files {
		if src == "" {
			continue
		}

		f, err = NewFile(string(src), dest);
		if err != nil {
			return err
		}
		bc.Files = append(bc.Files, f)
	}

	return nil
}

func (bc *BinCtrl) Build() ([]byte, error) {
	if bc.err != nil {
		return nil, bc.err
	}

	if bc.md5sBuild != true {
		return nil, errors.New("no md5 hashes found")
	}

	if err := bc.buildControlFile(); err != nil {
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

func (bc *BinCtrl) buildControlFile() error {

	if bc.err != nil {
		return bc.err
	}

	content := `Package: %s
Version: %s
Architecture: %s
Maintainer: %s
Description: %s`

	f := &File{
		Destination: "control",
		Content:     []byte(fmt.Sprintf(content, bc.Package, bc.Version, bc.Architecture, bc.Maintainer, bc.Description)),
	}

	bc.Files = append(bc.Files, f)
	return nil
}

func (bc *BinCtrl) buildConfFilesFile() error {
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

func (bc *BinCtrl) buildShLibsFile() error {
	if bc.err != nil {
		return bc.err
	}

	if len(bc.ShLibs) <= 0 {
		return nil
	}

	files := []string{}
	for _, file := range bc.ShLibs {
		if file != "" {
			files = append(files, string(file))
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

func (bc *BinCtrl) BuildMD5SumsFile(dataFiles []*File) error {
	if bc.err != nil {
		return bc.err
	}
	sums := make([]string, len(dataFiles))
	for i, file := range dataFiles {
		sums[i] = file.MD5 + " " + file.Destination
	}

	f := &File{
		Destination: "md5sums",
		Content:     []byte(strings.Join(sums, "\n")),
	}
	bc.Files = append(bc.Files, f)
	return nil
}