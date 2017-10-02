package main

import (
	"flag"
	"log"
	"os"
)

func main() {
	configFile := flag.String("c", "", "JSON formatted configuration")
	flag.Parse()

	if *configFile == "" {
		log.Fatal("-c configuration not given")
	}

	bd, err := NewBinaryDeb(*configFile)
	if err != nil {
		log.Fatal(err)
	}

	content, err := bd.Build()
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create(bd.Output)
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}

	_, err = f.Write(content)
	if err != nil {
		log.Fatal(err)
	}
}
