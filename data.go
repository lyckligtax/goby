package main

type BinData struct {}

func NewBinData(config []byte) (*BinData, error) {
	return nil, nil
}

func (bd *BinData) Build() *TarArchive {
	return nil
}

func (bd *BinData) Check() error {
	return nil
}