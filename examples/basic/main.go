package main

import (
	"log"

	"github.com/willemvds/Implodatron"
)

func main() {
	mainpy := implodatron.PythonFile{}
	mainpy.Path = "car.py"
	log.Println(mainpy)
	root := implodatron.BuildTree(mainpy)
	root.Print()
	print("\n")
}
