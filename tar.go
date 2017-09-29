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
	Content *bytes.Buffer
	Gzipped *bytes.Buffer
	w       *tar.Writer
	z       *gzip.Writer
	closed  bool
}

func NewTarArchive() *TarArchive {
	contentBuffer := new(bytes.Buffer)
	gzipBuffer := new(bytes.Buffer)

	return &TarArchive{
		Content: contentBuffer,
		Gzipped: gzipBuffer,
		w:       tar.NewWriter(contentBuffer),
		z:       gzip.NewWriter(gzipBuffer),
	}
}

func (t *TarArchive) AddFile(filename string, body []byte) error {
	if t.closed {
		return ErrArchiveClosed
	}
	hdr := &tar.Header{
		Name: filename,
		Mode: 0600,
		Size: int64(len(body)),
	}
	if err := t.w.WriteHeader(hdr); err != nil {
		return err
	}
	if _, err := t.w.Write(body); err != nil {
		return err
	}

	return nil
}

func (t *TarArchive) Close() error {
	if t.closed {
		return nil
	}
	t.closed = true
	return t.w.Close()
}

func (t *TarArchive) Gzip() error {
	if t.closed {
		return ErrArchiveClosed
	}

	err := t.Close()
	if err != nil {
		return err
	}

	_, err = t.z.Write(t.Content.Bytes())
	if err != nil {
		return err
	}

	if err := t.z.Close(); err != nil {
		return err
	}

	return nil
}
