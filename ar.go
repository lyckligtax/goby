package main

import (
	"bytes"
	"github.com/lyckligtax/argo"
	"log"
	"time"
)

type ArArchive struct {
	content *bytes.Buffer
	w       *ar.Writer

	closed bool
	err    error
}

func NewArArchive() *ArArchive {
	contentBuffer := new(bytes.Buffer)

	w, _ := ar.NewWriter(contentBuffer)
	return &ArArchive{
		content: contentBuffer,
		w:       w,
	}
}

func (a *ArArchive) AddFile(filename string, body []byte) error {
	if a.closed {
		return ErrArchiveClosed
	}

	if a.err != nil {
		return a.err
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

func (a *ArArchive) Data() ([]byte, error) {
	if a.err != nil {
		return nil, a.err
	}

	a.closed = true
	a.w.Close()
	return a.content.Bytes(), nil
}
