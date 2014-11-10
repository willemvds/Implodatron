package implodatron

import (
	"io/ioutil"
	"log"
	"strings"
)

type PythonPackage struct {
	Path string
}

type PythonFile struct {
	Pypkg *PythonPackage
	Path string
}

type ImportNode struct {
	Parent *ImportNode
	Children []*ImportNode
	PyFile *PythonFile
	Line string
	StartPos int
	EndPos int
}

func Slurp(pyfile PythonFile) *ImportNode {
	root := ImportNode{}
	root.PyFile = &pyfile
	for node := &root; node != nil; {
		src, err := ioutil.ReadFile(node.PyFile.Path)
		if err != nil {
			log.Fatalf("%s: %v\n", node.PyFile.Path, err)	
		}
		log.Printf("%s read: %d lines\n", node.PyFile.Path, len(src))
		lines := strings.Split(string(src), "\n")
		for _, line := range lines {
			if strings.Index(line, "import") == 0 {
				log.Println(line)
				what := strings.TrimRight(line[7:], "\n")
				childfile := PythonFile{Path:what+".py"}
				child := Slurp(childfile)
				node.Children = append(node.Children, child)
				log.Println(what)
			}
		}

		for i := range node.Children {
			// do stuff
			log.Println(i)
		}
		node = nil
	}
	return &root
}
