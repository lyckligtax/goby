package main

import (
	"encoding/json"
	"math"
	"os"
	"io/ioutil"
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
// It checks folders recursively
// It returns an error if any of these files is not readable or if an error happened before
// It maps the json files map from "key": "value" to "dest": "src"
func (bd *Data) checkDataFiles() error {
	if bd.err != nil {
		return bd.err
	}

	dataFiles, err := expandDataFiles(bd.DataFiles)
	if err != nil {
		return err
	}
	for dest, src := range dataFiles {
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

func expandDataFiles(fileList map[string]string) (map[string]string, error) {
	expandedFiles := map[string]string{}
	
	for dest, src := range fileList {
		srcStat, err := os.Stat(src)
		if err != nil {
			return nil, err
		}
		
		if srcStat.IsDir() {
			filesInSrcDir, err := ioutil.ReadDir(src)
			if err != nil {
				return nil, err
			}
			
			srcDirFileList := map[string]string{}
			for _, fileInSrcDir := range filesInSrcDir {
				srcDirFileList[dest + "/" + fileInSrcDir.Name()] = src + "/" + fileInSrcDir.Name()
			}
			
			filesInSrcFileList, err := expandDataFiles(srcDirFileList)
			if err != nil {
				return nil, err
			}
			
			for dest, src := range filesInSrcFileList {
				expandedFiles[dest] = src
			}
			continue
		}
		
		expandedFiles[dest] = src
	}
	return expandedFiles, nil
}