package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func dirents(dir string, printFiles bool) []os.FileInfo {
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return nil
	}
	if !printFiles {
		var onlyDirList []os.FileInfo = []os.FileInfo{}
		for _, entry := range entries {
			if entry.IsDir() {
				onlyDirList = append(onlyDirList, entry)
			}
		}
		entries = onlyDirList
	}
	return entries
}

func walkDir(out io.Writer, dir string, printFiles bool, graphicSymbol string) {
	var indent string
	files := dirents(dir, printFiles)
	length := len(files)
	for i, file := range files {
		if file.IsDir() {
			if length > i+1 {
				fmt.Fprintf(out, graphicSymbol+"├───"+"%s\n", file.Name())
				indent = graphicSymbol + "│\t"
			} else {
				fmt.Fprintf(out, graphicSymbol+"└───"+"%s\n", file.Name())
				indent = graphicSymbol + "\t"
			}
			subDir := filepath.Join(dir, file.Name())
			walkDir(out, subDir, printFiles, indent)
		} else if printFiles {
			if file.Size() > 0 {
				if length > i+1 {
					fmt.Fprintf(out, graphicSymbol+"├───%s (%vb)\n", file.Name(), file.Size())
				} else {
					fmt.Fprintf(out, graphicSymbol+"└───%s (%vb)\n", file.Name(), file.Size())
				}
			} else {
				if length > i+1 {
					fmt.Fprintf(out, graphicSymbol+"├───%s (empty)\n", file.Name())
				} else {
					fmt.Fprintf(out, graphicSymbol+"└───%s (empty)\n", file.Name())
				}
			}
		}
	}
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	walkDir(out, path, printFiles, "")
	return nil
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

