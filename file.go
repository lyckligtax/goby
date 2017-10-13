package main

import (
	"errors"
	"io/ioutil"
	"path"
)

type File struct {
	Source      string
	Destination string
	Content     []byte
	MD5         string
}

func NewFile(source string, destination string) (*File, error) {
	source = path.Clean(whiteTrim(source))
	destination = path.Clean(whiteTrim(destination))
	
	content, err := ioutil.ReadFile(source)
	if err != nil {
		return nil, err
	}

	if content == nil {
		return nil, errors.New("no content read: " + source + " -> " + destination)
	}

	return &File{
		Source:      source,
		Destination: destination,
		Content:     content,
		MD5:         md5ToString(content),
	}, nil
}
