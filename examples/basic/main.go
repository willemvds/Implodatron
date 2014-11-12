package main

import (
	"log"

	"github.com/willemvds/Implodatron"
)

func main() {
	mainpy := implodatron.PythonFile{}
	mainpy.Path = "app.py"
	log.Println(mainpy)
	root := implodatron.BuildTree(mainpy)
	root.Print()
	print("\n")
}
