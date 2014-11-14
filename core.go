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
	Root string
	Path string
}

func (pf PythonFile) GetWD() string {
	if idx := strings.LastIndex(pf.Path, "/"); idx != -1 {
		return pf.Path[0 : idx+1]
	}
	return ""
}

type ImportNode struct {
	Parent   *ImportNode
	Children []*ImportNode
	PyFile   *PythonFile
}

func FindImport(line string) string {
	if strings.Index(line, "import") == 0 {
		what := strings.TrimRight(line[7:], "\n")
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

func Import(partial string, paths []string) (string, []byte, error) {
	var err error
	for _, path := range paths {
		log.Println(path)
		log.Println(partial)
		_, err = os.Stat(path + partial)
		if err == nil {
			src, err := ioutil.ReadFile(path + partial)
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

func Slurp(fromFile PythonFile, paths []string, intoNode *ImportNode) {
	root, src, err := Import(fromFile.Path, paths) //ioutil.ReadFile(fromFile.Path)
	if err != nil {
		log.Printf("%s: %v\n", fromFile, err)
		return
	}
	fromFile.Root = root
	paths = append([]string{root}, paths...)
	log.Printf("%s read: %d bytes\n", fromFile.Path, len(src))
	lines := strings.Split(string(src), "\n")

	for _, line := range lines {
		path := FindImport(line)
		if len(path) > 0 {
			//path = fromFile.GetWD() + path
			//paths = append(paths, fromFile.GetWD() + path)
			log.Println(line, "->", path)
			pyfile := &PythonFile{
				Path: path,
			}
			child := &ImportNode{
				Parent: intoNode,
				PyFile: pyfile,
			}
			intoNode.Children = append(intoNode.Children, child)
			if !child.FindPath(path) {
				Slurp(*pyfile, paths, child)
			}
		}
	}
}

func BuildTree(pyfile PythonFile, paths []string) *ImportNode {
	root := &ImportNode{
		PyFile: &pyfile,
	}
	Slurp(pyfile, paths, root)
	return root
}
