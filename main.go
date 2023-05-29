package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	pathPrefix     string
	mapIgnoreFiles = map[string][]string{}
)

func loadAllfolders(path string, folderList []string) []string {
	stat, err := os.Stat(path)
	if err != nil {
		log.Fatalf("Error reading path %q: %v\n", path, err)
	}

	if !stat.IsDir() {
		return folderList
	}

	if strings.HasPrefix(filepath.Base(path), ".") {
		return folderList
	}

	folderList = append(folderList, path)

	folders, err := filepath.Glob(path + "/*")
	if err != nil {
		log.Fatalf("Error reading path %q: %v\n", path, err)
	}

	for _, folder := range folders {
		folderList =
			loadAllfolders(folder, folderList)
	}

	return folderList
}

func filterFolderWithGitingore(folderList []string) []string {
	ret := []string{}
	for _, folder := range folderList {
		_, err := os.Stat(folder + "/.gitignore")
		if err != nil {
			continue
		}

		ret = append(ret, folder)
	}

	return ret
}

// função que carrega todos os gitignore
func loadAllGitignore(path string) []string {
	ret := []string{}
	b, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Error reading path %q: %v\n", path, err)
	}

	lines := strings.Split(string(b), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "#") {
			continue
		}

		ret = append(ret, line)
	}

	return ret
}

// removeDuplicate
func removeDuplicate(list []string) []string {
	ret := []string{}
	for _, item := range list {
		found := false
		for _, item2 := range ret {
			if item == item2 {
				found = true
				break
			}
		}

		if !found {
			ret = append(ret, item)
		}
	}

	return ret
}

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

	folderList := loadAllfolders(root, []string{})
	folderList = filterFolderWithGitingore(folderList)

	for _, folder := range folderList {
		mapIgnoreFiles[folder] = loadAllGitignore(folder + "/.gitignore")
	}

	for k, v := range mapIgnoreFiles {
		for kk, vv := range mapIgnoreFiles {
			if strings.HasPrefix(kk, k) {
				mapIgnoreFiles[kk] = append(v, vv...)
			}
		}
	}

	for k, v := range mapIgnoreFiles {
		mapIgnoreFiles[k] = removeDuplicate(v)
	}

	for k, v := range mapIgnoreFiles {
		fmt.Printf("%q\n", k)
		for _, vv := range v {
			fmt.Printf("\t%q: %q\n", k, vv)
		}
	}

	return

	pathPrefix = root + "/"

	err = filepath.Walk(root, visitor)
	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", root, err)
	}
}
