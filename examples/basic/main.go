package main

import (
	"flag"
	"log"

	"github.com/willemvds/Implodatron"
)

var path string

func init() {
	flag.StringVar(&path, "path", "", "path to 'main' file")
}

func main() {
	flag.Parse()
	if path == "" {
		log.Fatalf("Usage: basic -path <myapp.py>")
	}
	mainpy := implodatron.PythonFile{
		Path: path,
	}
	log.Println(mainpy)
	root := implodatron.BuildTree(mainpy)
	root.Print()
	print("\n")
}
