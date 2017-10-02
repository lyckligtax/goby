package main

// Goby resembles the content of a deb-package as found on Ubuntu or Debian derivates
// it consists of a control.tar and a data.tar both optionally gzipped
// and a debian-binary file containing the package semantic version
// all three files are packaged into an ar-Archive
type Goby struct {
	ctrl Tar
	data Tar
	gzip bool
	err  error
}

func NewGoby(opt Config) (g *Goby) {
	g = &Goby{}

	if opt.packageType == "binary" {
		/*
			TODO: implement source type debian packages
		 */
		g.ctrl = NewBinCtrl(opt)
		g.data = NewBinData(opt)
	} else {
		g.ctrl = NewBinCtrl(opt)
		g.data = NewBinData(opt)
	}
	return
}

func (g *Goby) Check() (err error) {
	if g.err != nil {
		return g.err
	}

	if err = g.ctrl.Check(); err != nil {
		g.err = err
		return
	}

	if err = g.data.Check(); err != nil {
		g.err = err
		return
	}
	return
}

func (g *Goby) Build() ([]byte, error) {
	if g.err != nil {
		return nil, g.err
	}
	deb := NewArArchive()

	// Ubuntu is picky about debian-binary being first in the package
	deb.AddFile("debian-binary", []byte("2.0\n"))

	// Build gzipped control tar archive
	ctrl, err := g.ctrl.Build().Gzip()
	if err != nil {
		g.err = err
		return nil, err
	}

	// Add control archive to ar-Archive
	if err = deb.AddFile("control.tar.gz", ctrl);  err != nil {
		g.err = err
		return nil, err
	}

	// Build gzipped data tar archive
	data, err := g.data.Build().Gzip()
	if err != nil {
		g.err = err
		return nil, err
	}


	// Add data archive to ar-Archive
	if err = deb.AddFile("data.tar.gz", data); err != nil {
		g.err = err
		return nil, err
	}

	//TODO: should be deferred
	if err = deb.Close(); err != nil {
		return nil, err
	}

	return deb.Content.Bytes(), nil
}
