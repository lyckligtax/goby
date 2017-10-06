package main

import (
	"encoding/json"
	"math"
)

type Data struct {
	DataFiles map[string]string `json:"dataFiles"`

	// Package Files
	Files []*File

	err error
}

func NewData(config []byte) (*Data, error) {
	data := &Data{}
	err := json.Unmarshal(config, data)

	if err != nil {
		return nil, err
	}

	return data, nil
}

// Build takes alles Files from bd.Files and writes them into a tar Archive
// It returns the either
// 		(good) the gzipped bytes content of this archive and an nil error object
// 		(bad) an
func (bd *Data) Build() ([]byte, error) {
	if bd.err != nil {
		return nil, bd.err
	}

	// create tar archive
	dataTar := NewTarArchive()

	// add all files to the archive
	for _, f := range bd.Files {
		if f != nil {
			dataTar.AddFile(f.Destination, f.Content)
		}
	}

	gzippedTar, err := dataTar.Gzip()
	if err != nil {
		bd.err = err

	}
	return gzippedTar, bd.err
}

func (bd *Data) Check() error {
	if bd.err != nil {
		return bd.err
	}

	if err := bd.checkDataFiles(); err != nil {
		bd.err = err
	}

	return bd.err
}

// CheckDataFiles populates bd.Files List from the json files map
// It returns an error if any of these files is not readable or if an error happened before
// It maps the json files map from "key": "value" to "dest": "src"
func (bd *Data) checkDataFiles() error {
	if bd.err != nil {
		return bd.err
	}

	for dest, src := range bd.DataFiles {
		f, err := NewFile(src, dest)
		if err != nil {
			return err
		}
		bd.Files = append(bd.Files, f)
	}

	return nil
}

// DataSize returns the Actual Size of all DataFiles in kb
// Size in bytes / 1024 rounded up
func (bd *Data) DataSizeKB() int {
	size := 0
	for _, f := range bd.Files {
		size += len(f.Content)
	}

	return int(math.Ceil(float64(size) / 1024))
}
