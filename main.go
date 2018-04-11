package main

import (
	"flag"
	"log"
)

const (
	defaultPackageName = "$name-$version.deb"
)

func main() {
	configFile := flag.String("c", "goby.json", "JSON formatted configuration")
	outputFile := flag.String("o", defaultPackageName, "output file")
	flag.Parse()

	g, err := NewGoby(*configFile)
	if err != nil {
		log.Fatal(err)
	}

	if err = g.Make(*outputFile); err != nil {
		log.Fatal(err)
	}
}
