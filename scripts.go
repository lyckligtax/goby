package main

type Scripts struct {
	PreInst  string `json:"preInst"`
	PostInst string `json:"postInst"`
	PreRm    string `json:"preRm"`
	PostRm   string `json:"postRm"`
	Config   string `json:"config"`
}
