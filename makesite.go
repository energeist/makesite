package main

import (
	"flag"
	"fmt"
	"html/template"
	// "io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Page holds all the information we need to generate a new
// HTML page from a text file on the filesystem.
type Page struct {
	Title				 string
	Body				 string	
}

func main() {
	startTime := time.Now()
	fileFlag := flag.String("file", "text/first-post.txt", "String representing a .txt file to read from")
	dirFlag := flag.String("dir", "text", "String representing a directory to read from")
	flag.Parse()

	fmt.Println("File: ", *fileFlag)
	fmt.Println("Dir: ", *dirFlag)

	textFiles, _ := getTextFiles(*dirFlag)
	fmt.Println(textFiles)

	for _, file := range textFiles {
		processTextFiles(file)
	}

	endTime := time.Now()

	fmt.Println("Time elapsed: ", endTime.Sub(startTime))
}

func getTextFiles(dir string) ([]string, error) {
	var files []string

	// Walk through the directory
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if the file is a .txt file
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".txt") {
			fmt.Println("Found text file: ", path)

			fmt.Println("File name: ", info.Name())

			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

func processTextFiles(file string) {
	fileContents, err := readFile(file)
	if err != nil {
		panic(err)
	}

	fmt.Println("File contents: ", string(fileContents))

	title, body := parseContent(fileContents)

	page := createPage(title, body)

	err = processTemplate("template.tmpl", page)
	if err != nil {
		panic(err)
	}

	err = createHTMLFile(file, page)
	if err != nil {
		panic(err)
	}
}

func readFile(fileName string) ([]byte, error) {
	return ioutil.ReadFile(fileName)
}

func parseContent(contents []byte) (string, string) {
	// TODO: fix the parsing so that the first line is the title and the rest is the body
	title := string(contents[0:16])
	body := strings.TrimSpace(string(contents[16:]))
	return title, body
}

func createPage(title, body string) Page {
	return Page{
		Title: title,
		Body:  body,
	}
}

func processTemplate(templateName string, page Page) error {
	tmpl, err := template.ParseFiles(templateName)
	if err != nil {
		return err
	}
	return tmpl.Execute(os.Stdout, page)
}

func createHTMLFile(fileName string, page Page) error {
	splitFileName := strings.Split(fileName, "/")
	fileName = splitFileName[len(splitFileName)-1]
	newFileName := strings.TrimSuffix(fileName, ".txt") + ".html"
	newFileName = "pages/" + newFileName
	newFile, err := os.Create(newFileName)
	if err != nil {
		return err
	}
	defer newFile.Close()

	tmpl, err := template.ParseFiles("template.tmpl")
	if err != nil {
		return err
	}
	return tmpl.Execute(newFile, page)
}

// func getTextFiles(directory string) ([]string, error) {
// 	fmt.Println("Parsing directory: ", directory)
// 	files, err := ioutil.ReadDir(directory)
// 	if err != nil {
// 		return nil, err
// 	}

// 	textFiles := []string{}

// 	for _, file := range files {
// 		if file.IsDir() {
// 			fmt.Println("New Directory: ", file.Name())
// 			return parseDir(file.Name())
// 		} else {
// 		fmt.Println(file.Name())
// 		textFiles = append(textFiles, file.Name())
// 		}
// 	}
// 	return textFiles, nil
// }

