package main

import (
"fmt"
"os"
"path/filepath"
"strings"
)

var paths = make([]string, 0)

func walkDir(path string) {
	errors := make(chan error)
	done := make(chan bool)

	// Error handler
	go func() {
		for err := range errors {
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err) // TODO: send to a file
			}
		}
		done <- true
	}()

	//fmt.Printf(" + Checking: %q...\n", path)
	filepath.Walk(path, search(errors))
	close(errors)
	<-done
	fmt.Println()
}

// Search files with ".go~" extension.
func search(errors chan error) filepath.WalkFunc {
	// Implements "filepath.WalkFunc".
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			errors <- err
			return nil
		}

		if info.IsDir() {
			fmt.Printf(" + checking: %q\n", path)
		}

		// Skip directories ang hidden files
		if info.IsDir() || filepath.HasPrefix(info.Name(), ".") {
			return nil
		}

		fileExt := strings.ToLower(filepath.Ext(info.Name()))
		if fileExt == ".go~" {
			paths = append(paths, path)
		}

		return nil
	}
}

func main() {
	walkDir("/home/neo/go/src/github.com/kless/GoWizard")
	fmt.Println(paths)
}
