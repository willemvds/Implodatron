package main

import (
	"flag"
	//"log"
	"os"
	"strings"

	"github.com/willemvds/Implodatron"
)

//var path string
var verbose bool

func init() {
	flag.BoolVar(&verbose, "v", false, "verbose")
}

func main() {
	flag.Parse()
	/*
		if path == "" {
			log.Fatalf("Usage: basic -path <myapp.py>")
		}
	*/
	if verbose {
		implodatron.LogLevel = 1
	}

	for _, path := range flag.Args() {
		var pathroot string
		sidx := strings.LastIndex(path, "/")
		if sidx != -1 {
			sidx++
			pathroot = path[0:sidx]
			path = path[sidx:]
		} else {
			pathroot, _ = os.Getwd()
			pathroot = pathroot + "/"
		}

		implodatron.BuildTree(path, []string{pathroot})
		//root.Print()
		print("\n")
	}
}
