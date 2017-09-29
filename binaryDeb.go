package main

import (
	"errors"
	"fmt"
	"strings"
)

type BinaryControl struct {
	Package      string `json:"package"`
	Version      string `json:"version"`
	Architecture string `json:"architecture"`
	Maintainer   string `json:"maintainer"`
	Description  string `json:"description"`
}

type BinaryDeb struct {
	Deb
	BinaryControl
}

func NewBinaryDeb(filename string) (*BinaryDeb, error) {
	bd := &BinaryDeb{}
	err := readJSON(filename, bd)
	if err != nil {
		return nil, err
	}

	err = bd.checkMandatoryFields()
	if err != nil {
		bd.err = err
		return nil, err
	}

	return bd, nil
}

func (bd *BinaryDeb) Build() ([]byte, error) {
	err := bd.checkSourceFiles()
	if err != nil {
		return nil, err
	}

	err = bd.checkControlScripts()
	if err != nil {
		return nil, err
	}

	bd.fileControl()
	bd.fileConfFiles()
	bd.fileMD5Sums()
	bd.fileShLibs()

	control, err := bd.buildControlPackage()
	if err != nil {
		bd.err = err
	}

	source, err := bd.buildSourcePackage()
	if err != nil {
		bd.err = err
	}

	// ar
	deb := NewArArchive()
	deb.AddFile("debian-binary", []byte("2.0\n"))
	deb.AddFile("control.tar.gz", control)
	deb.AddFile("data.tar.gz", source)
	deb.Close()

	return deb.Content.Bytes(), bd.err
}

func (bd *BinaryDeb) checkMandatoryFields() error {
	if bd.err != nil {
		return bd.err
	}

	if bd.Package == "" {
		return errors.New("invalid package name")
	}
	if !isSemanticVersionRxp.MatchString(bd.Version) {
		return errors.New("invalid version")
	}
	if bd.Architecture == "" {
		return errors.New("invalid architecture")
	}
	if bd.Maintainer == "" {
		return errors.New("invalid maintainer")
	}
	if bd.Description == "" {
		return errors.New("invalid description")
	}

	return nil
}

func (bd *BinaryDeb) checkControlScripts() error {
	if bd.err != nil {
		return bd.err
	}

	bd.appendFile(bd.PreInst, "preinst", false)
	bd.appendFile(bd.PostInst, "postinst", false)
	bd.appendFile(bd.PreRm, "prerm", false)
	bd.appendFile(bd.PostRm, "postrm", false)
	bd.appendFile(bd.Config, "config", false)

	return bd.err
}

func (bd *BinaryDeb) checkSourceFiles() error {
	if bd.err != nil {
		return bd.err
	}

	for source, destination := range bd.Files {
		bd.appendFile(source, destination, true)
	}

	return bd.err
}

func (bd *BinaryDeb) buildControlPackage() ([]byte, error) {
	if bd.err != nil {
		return nil, bd.err
	}

	t := NewTarArchive()
	for _, f := range bd.ControlFiles {
		if f == nil {
			continue
		}
		t.AddFile(f.Destination, f.Content)
	}

	if err := t.Gzip(); err != nil {
		return nil, err
	}

	return t.Gzipped.Bytes(), nil
}

func (bd *BinaryDeb) fileControl() error {
	if bd.err != nil {
		return bd.err
	}
	controlFile := `Package: %s
Version: %s
Architecture: %s
Maintainer: %s
Description: %s
`
	f := &File{
		Destination: "control",
		Content:     []byte(fmt.Sprintf(controlFile, bd.Package, bd.Version, bd.Architecture, bd.Maintainer, bd.Description)),
	}
	bd.ControlFiles = append(bd.ControlFiles, f)
	return nil
}

func (bd *BinaryDeb) fileConfFiles() error {
	if bd.err != nil {
		return bd.err
	}
	if len(bd.ConfFiles) <= 0 {
		return nil
	}
	f := &File{
		Destination: "conffiles",
		Content:     []byte(strings.Join(bd.ConfFiles, "\n")),
	}
	bd.ControlFiles = append(bd.ControlFiles, f)
	return nil
}

func (bd *BinaryDeb) fileShLibs() error {
	if bd.err != nil {
		return bd.err
	}
	if len(bd.ShLibs) <= 0 {
		return nil
	}
	f := &File{
		Destination: "shlibs",
		Content:     []byte(strings.Join(bd.ShLibs, "\n")),
	}
	bd.ControlFiles = append(bd.ControlFiles, f)
	return nil
}

func (bd *BinaryDeb) fileMD5Sums() error {
	if bd.err != nil {
		return bd.err
	}
	sums := make([]string, len(bd.SourceFiles))
	for i, file := range bd.SourceFiles {
		sums[i] = file.MD5 + " " + file.Destination
	}

	f := &File{
		Destination: "md5sums",
		Content:     []byte(strings.Join(sums, "\n")),
	}
	bd.ControlFiles = append(bd.ControlFiles, f)
	return nil
}

func (bd *BinaryDeb) buildSourcePackage() ([]byte, error) {
	if bd.err != nil {
		return nil, bd.err
	}

	t := NewTarArchive()
	for _, f := range bd.SourceFiles {
		if f == nil {
			continue
		}
		t.AddFile(f.Destination, f.Content)
	}

	if err := t.Gzip(); err != nil {
		return nil, err
	}

	return t.Gzipped.Bytes(), nil
}
