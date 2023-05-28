package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var pathPrefix string

func visitor(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if info.IsDir() && strings.HasPrefix(info.Name(), ".") {
		return filepath.SkipDir
	}

	if info.IsDir() {
		return nil
	}

	base := filepath.Base(path)
	if strings.HasPrefix(base, ".") {
		return nil
	}

	ppath := strings.TrimPrefix(path, pathPrefix)

	fmt.Printf("%q\n", ppath)

	return nil
}

func main() {
	root := "./"

	// convert para fillpath
	root, err := filepath.Abs(root)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	pathPrefix = root + "/"

	err = filepath.Walk(root, visitor)
	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", root, err)
	}
}
