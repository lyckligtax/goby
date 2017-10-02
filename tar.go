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
	closed  bool
}

func NewTarArchive() *TarArchive {
	contentBuffer := new(bytes.Buffer)

	return &TarArchive{
		content: contentBuffer,
		w:       tar.NewWriter(contentBuffer),
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

func (t *TarArchive) close() error {
	if t.closed {
		return nil
	}
	t.closed = true
	return t.w.Close()
}

func (t *TarArchive) Data() ([]byte, error) {
	if t.closed {
		return nil, ErrArchiveClosed
	}

	err := t.close()
	if err != nil {
		return nil, err
	}

	return t.content.Bytes(), nil
}

func (t *TarArchive) Gzip() ([]byte, error) {
	if t.closed {
		return nil, ErrArchiveClosed
	}

	err := t.close()
	if err != nil {
		return nil, err
	}

	gzipBuffer := new(bytes.Buffer)
	z := gzip.NewWriter(gzipBuffer)
	_, err = z.Write(t.content.Bytes())
	if err != nil {
		return nil, err
	}

	if err := z.Close(); err != nil {
		return nil, err
	}

	return gzipBuffer.Bytes(), nil
}
