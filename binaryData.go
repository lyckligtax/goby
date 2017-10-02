package main

type BinData struct {}

func NewBinData(opt Config) *BinData {
	return nil
}

func (bd *BinData) Build() *TarArchive {
	return nil
}

func (bd *BinData) Check() error {
	return nil
}