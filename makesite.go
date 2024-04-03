package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"strings"
)

// Page holds all the information we need to generate a new
// HTML page from a text file on the filesystem.
type Page struct {
	Title				 string
	Body				 string	
}

func main() {
	fileFlag := flag.String("file", "first-post.txt", "String representing a .txt file to read from")
	flag.Parse()

	fmt.Println("File: ", *fileFlag)

	fileContents, err := readFile(*fileFlag)
	if err != nil {
		panic(err)
	}

	title, body := parseContent(fileContents)

	page := createPage(title, body)

	err = processTemplate("template.tmpl", page)
	if err != nil {
		panic(err)
	}

	err = createHTMLFile(*fileFlag, page)
	if err != nil {
		panic(err)
	}
}

func readFile(fileName string) ([]byte, error) {
	return ioutil.ReadFile(fileName)
}

func parseContent(contents []byte) (string, string) {
	title := string(contents[0:16])
	body := strings.TrimSpace(string(contents[16:]))
	fmt.Println("Title: ", title)
	fmt.Println("Body: ", body)
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
	newFileName := strings.TrimSuffix(fileName, ".txt") + ".html"
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