package main

type Tar interface {
	Build() *TarArchive
	Check() error
}