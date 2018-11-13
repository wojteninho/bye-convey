package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/wojteninho/bye-convey"
	"log"
	"os"
	"strings"

	"github.com/wojteninho/scanner/pkg/scanner"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Bye-Convey is a tool to get rid of GoConvey and replace it with native Golang structs for the tests organization + Gomega for assertions.\n\nUsage\n\n\t ./bye [files or directories]")
		os.Exit(1)
	}

	var (
		files []string
		directories []string
	)

	for _, item := range os.Args[1:] {
		info, err := os.Stat(item)
		if err != nil {
			log.Fatalf("Error: Unable to proceed file %s due to error %s", item, err)
		}

		if info.IsDir() {
			directories = append(directories, item)
			continue
		}

		if !strings.HasSuffix(info.Name(), "_test.go") {
			log.Fatalf("Error: Unable to proceed file %s. It is not not a test file.", item)
		}

		files = append(files, item)
	}

	var (
		ctx = context.Background()
		s   = scanner.NewBuilder().
			Recursive().
			Files().
			Match(scanner.ExtensionFilter("_test.go")).
			In(directories...).
			MustBuild()
	)

	for _, file := range files {
		do(file)
	}

	for file := range scanner.MustScan(s.Scan(ctx)) {
		do(file.String())
	}
}

func do(filename string) {
	log.Printf("Processing: %s", filename)

	fileHandle, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Unable to open %s due to %s", filename, err)
	}
	defer fileHandle.Close()
	fileScanner := bufio.NewScanner(fileHandle)
	var fileBuffer bytes.Buffer

	for fileScanner.Scan() {
		current := fileScanner.Bytes()
		transformed := bye_convey.Transform(current)
		fileBuffer.Write(transformed)
		fileBuffer.Write([]byte("\n"))
	}

	filenameBuffer := filename + ".buffer"
	fileBufferHandle, err := os.Create(filenameBuffer)
	if err != nil {
		log.Fatalf("Unable to create %s due to %s", filenameBuffer, err)
	}
	defer fileBufferHandle.Close()

	fileBufferHandle.Write(fileBuffer.Bytes())
	err = os.Rename(filenameBuffer, filename)
	if err != nil {
		log.Fatalf("Unable to rename %s -> %s due to %s", filenameBuffer, filename, err)
	}
}
