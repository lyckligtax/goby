package main

import (
	"bytes"
	"github.com/lyckligtax/argo"
	"log"
	"time"
)

type ArArchive struct {
	Content *bytes.Buffer
	w       *ar.Writer
	closed  bool
}

func NewArArchive() *ArArchive {
	contentBuffer := new(bytes.Buffer)

	w, _ := ar.NewWriter(contentBuffer)
	return &ArArchive{
		Content: contentBuffer,
		w:       w,
	}
}

func (a *ArArchive) AddFile(filename string, body []byte) error {
	if a.closed {
		return ErrArchiveClosed
	}

	hdr := ar.Header{
		ObjectName:    filename,
		MTime:         time.Now(),
		UID:           0,
		GID:           0,
		FileMode:      0600,
		ContentLength: int64(len(body)),
	}
	_, err := a.w.WriteHeader(hdr)
	if err != nil {
		log.Fatal(err)
		return err
	}

	_, err = a.w.Write(body)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func (a *ArArchive) Close() error {
	if a.closed {
		return nil
	}
	a.closed = true
	a.w.Close()
	return nil
}
