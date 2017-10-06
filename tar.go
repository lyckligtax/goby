package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"errors"
)

var (
	ErrArchiveClosed = errors.New("archive already closed")
)

type TarArchive struct {
	content *bytes.Buffer
	w       *tar.Writer

	closed bool
	err    error
}

func NewTarArchive() *TarArchive {
	contentBuffer := new(bytes.Buffer)

	return &TarArchive{
		content: contentBuffer,
		w:       tar.NewWriter(contentBuffer),
	}
}

func (t *TarArchive) AddFile(filename string, body []byte) error {
	if t.err != nil {
		return t.err
	}

	if t.closed {
		return ErrArchiveClosed
	}

	hdr := &tar.Header{
		Name: filename,
		Mode: 0600,
		Size: int64(len(body)),
	}

	if err := t.w.WriteHeader(hdr); err != nil {
		t.err = err
		return err
	}

	if _, err := t.w.Write(body); err != nil {
		t.err = err
		return err
	}

	return nil
}

func (t *TarArchive) Gzip() ([]byte, error) {
	if t.err != nil {
		return nil, t.err
	}

	t.closed = true
	if err := t.w.Close(); err != nil {
		t.err = err
		return nil, err
	}

	gzipBuffer := new(bytes.Buffer)
	z := gzip.NewWriter(gzipBuffer)
	defer z.Close()

	if _, err := z.Write(t.content.Bytes()); err != nil {
		t.err = err
		return nil, err
	}

	if err := z.Close(); err != nil {
		t.err = err
		return nil, err
	}

	return gzipBuffer.Bytes(), nil
}
