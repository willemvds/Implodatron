package implodatron

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type PythonFile struct {
	Name string
	Dir  string
}

func NewPythonFile(path string) *PythonFile {
	pf := PythonFile{}
	sidx := strings.LastIndex(path, "/")
	pf.Name = path[sidx+1:]
	pf.Dir = path[0 : sidx+1]
	return &pf
}

type ImportNode struct {
	Parent     *ImportNode
	Children   []*ImportNode
	PyFile     *PythonFile
	ImportLine string
}

func FindImport(line string) string {
	if strings.Index(line, "import") == 0 {
		what := strings.TrimRight(line[7:], " \n")
		return strings.Replace(what, ".", "/", -1) + ".py"
	}
	if strings.Index(line, "from") == 0 {
		from := 5
		to := strings.Index(line, "import")
		what := line[from : to-1]
		return strings.Replace(what, ".", "/", -1) + ".py"
	}
	return ""
}

func Import(name string, paths []string) (string, []byte, error) {
	var err error
	for _, path := range paths {
		_, err = os.Stat(path + name)
		if err == nil {
			src, err := ioutil.ReadFile(path + name)
			return path, src, err
		}
	}
	return "", []byte{}, errors.New("import not found")
}

func PrintNode(n *ImportNode, level int) {
	if len(n.Children) == 0 {
		return
	}
	fmt.Printf("\n")
	level++
	fmt.Printf("%d:", level)
	for i := range n.Children {
		if n.Children[i].PyFile != nil {
			fmt.Printf(" %s (%s)", n.Children[i].PyFile.Name, n.Children[i].PyFile.Dir)
		}
		PrintNode(n.Children[i], level)
	}
}

func (n *ImportNode) Print() {
	level := 0
	if n.PyFile != nil {
		fmt.Printf("%d: %s", level, n.PyFile.Name)
	}
	PrintNode(n, level)
}

func (n *ImportNode) FindImport(name string, dir string) bool {
	for node := n.Parent; node != nil; node = node.Parent {
		if node.PyFile.Name == name && node.PyFile.Dir == dir {
			return true
		}
	}
	return false
}

func Slurp(fromName string, paths []string, intoNode *ImportNode) {
	path, src, err := Import(fromName, paths)
	if err != nil {
		//log.Printf("%s: %v\n", fromName, err)
		return
	}
	fromFile := NewPythonFile(path + fromName)
	intoNode.PyFile = fromFile
	intoNode.ImportLine = fromName
	if intoNode.FindImport(fromFile.Name, fromFile.Dir) {
		log.Println("\u001B[0;31mCIRCULAR IMPORT; STOP THE PLANET\u001B[0;m")
		for n := intoNode; n != nil; n = n.Parent {
			fmt.Printf("%s", n.ImportLine)
			if n.Parent != nil {
				fmt.Printf(" <- ")
			}
		}
		fmt.Println()
		return
	}
	paths = append([]string{fromFile.Dir}, paths...)

	//log.Printf("%s read: %d bytes\n", fromFile.Name, len(src))
	lines := strings.Split(string(src), "\n")
	for _, line := range lines {
		partial := FindImport(line)
		if len(partial) > 0 {
			//log.Println(line, "->", partial)
			child := &ImportNode{
				Parent: intoNode,
			}
			intoNode.Children = append(intoNode.Children, child)
			Slurp(partial, paths, child)
		}
	}
}

func BuildTree(name string, paths []string) *ImportNode {
	root := &ImportNode{}
	Slurp(name, paths, root)
	return root
}
