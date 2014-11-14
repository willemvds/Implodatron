package main

import (
	"flag"
	"log"
	"os"
	"strings"

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
	var pathroot string
	sidx := strings.LastIndex(path, "/")
	if sidx != -1 {
		sidx++
		pathroot = path[0:sidx]
		path = path[sidx:]
	} else {
		pathroot, _ = os.Getwd()
	}

	mainpy := implodatron.PythonFile{
		Root: pathroot,
		Path: path,
	}
	log.Println(mainpy)
	root := implodatron.BuildTree(mainpy, []string{pathroot})
	root.Print()
	print("\n")
}
