package main

type Deb struct {
	SourceFiles  []*File
	ControlFiles []*File
	Scripts
	Files     map[string]string `json:"files"`
	ConfFiles []string          `json:"confFiles"`
	ShLibs    []string          `json:"shLibs"`
	Output    string            `json:"output"`
	err       error
}

func (d *Deb) appendFile(source string, destination string, isSource bool) error {
	if d.err != nil {
		return d.err
	}

	f, err := NewFile(source, destination)
	if err != nil {
		d.err = err
		return err
	}

	if isSource {
		d.SourceFiles = append(d.SourceFiles, f)
	} else {
		d.ControlFiles = append(d.ControlFiles, f)
	}

	return nil
}
