package implodatron

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

type PythonFile struct {
	Path string
}

type ImportNode struct {
	Parent   *ImportNode
	Children []*ImportNode
	PyFile   *PythonFile
}

func PrintNode(n *ImportNode, level int) {
	if len(n.Children) == 0 {
		return
	}
	fmt.Printf("\n")
	level++
	fmt.Printf("%d:", level)
	for i := range n.Children {
		fmt.Printf(" %s", n.Children[i].PyFile.Path)
		PrintNode(n.Children[i], level)
	}
}

func (n *ImportNode) Print() {
	level := 0
	fmt.Printf("%d: %s", level, n.PyFile.Path)
	PrintNode(n, level)
}

func (n *ImportNode) FindPath(p string) bool {
	for node := n.Parent; node != nil; node = node.Parent {
		if node.PyFile.Path == p {
			return true
		}
	}
	return false
}

func Slurp(fromFile PythonFile, intoNode *ImportNode) {
	src, err := ioutil.ReadFile(fromFile.Path)
	if err != nil {
		log.Fatalf("%s: %v\n", fromFile, err)
	}
	log.Printf("%s read: %d lines\n", fromFile.Path, len(src))
	lines := strings.Split(string(src), "\n")

	for _, line := range lines {
		if strings.Index(line, "import") == 0 {
			log.Println(line)
			what := strings.TrimRight(line[7:], "\n")
			path := what + ".py"
			pyfile := &PythonFile{
				Path: path,
			}
			child := &ImportNode{
				Parent: intoNode,
				PyFile: pyfile,
			}
			intoNode.Children = append(intoNode.Children, child)
			if !child.FindPath(path) {
				Slurp(*pyfile, child)
			}
		}
	}
}

func BuildTree(pyfile PythonFile) *ImportNode {
	root := &ImportNode{
		PyFile: &pyfile,
	}
	Slurp(pyfile, root)
	return root
}
