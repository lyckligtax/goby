package main

import (
	"errors"
	"io/ioutil"
)

type File struct {
	Source      string
	Destination string
	Content     []byte
	MD5         string
}

func NewFile(source string, destination string) (*File, error) {
	content, err := ioutil.ReadFile(source)
	if err != nil {
		return nil, err
	}

	if content == nil || len(content) == 0 {
		return nil, errors.New("no content read")
	}

	return &File{
		Source:      whiteTrim(source),
		Destination: whiteTrim(destination),
		Content:     content,
		MD5:         md5ToString(content),
	}, nil
}
