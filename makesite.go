package main

import (
	"flag"
	"fmt"
	"github.com/ttacon/chalk"
	"github.com/jbrodriguez/mlog"
	"html/template"
	"io/fs"
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
	// start mlog with a log level of info, log to makesite.log, rotate at 100kB, keep 5 logs
	mlog.StartEx(mlog.LevelInfo, "makesite.log", 100*1024, 5)
	mlog.Info("Starting makesite...")
	startTime := time.Now()

	// fileFlag := flag.String("file", "text/first-post.txt", "String representing a .txt file to read from")
	// dirFlag := flag.String("dir", "text", "String representing a directory to read from")
	fileFlag := flag.String("file", "", "String representing a .txt file to read from")
	dirFlag := flag.String("dir", "", "String representing a directory to read from")
	flag.Parse()

	if *fileFlag == "" && *dirFlag == "" {
		mlog.Warning("No file or directory specified. Defaulting to demo parse of /text/latest-post.txt.")
		*fileFlag = "text/latest-post/latest-post.txt"
	} else {
		mlog.Info("Parsed flags from CLI:")
		mlog.Info("File: %v", *fileFlag)
		mlog.Info("Dir: %v", *dirFlag)
	}

	if *dirFlag != "" {
		textFiles, _ := getTextFiles(*dirFlag)
		
		mlog.Info("Found %v text files in directory /%v/:", len(textFiles), *dirFlag)
		mlog.Info(strings.Join(textFiles, ", "))

		for _, file := range textFiles {
			processTextFiles(file)
		}

		endTime := time.Now()
		
		htmlSize, _ := calcFileSizeInDirectory("pages")

		greenAndBold := chalk.Green.NewStyle().
			WithTextStyle(chalk.Bold).
			Style
		
		// This style didn't work properly in my terminal so I just used vanilla escape codes instead.
		// bold := chalk.Bold.NewStyle().Style

		// This is supposed to be green but may appear differently in your terminal depending on theme settings, it seems.
		fmt.Printf(greenAndBold("\nSuccess!") + " Generated \033[1m%v\033[0m HTML files in the pages directory.  Wrote %.1f kB in %v.\n", len(textFiles), htmlSize, endTime.Sub(startTime))
		mlog.Info("Success! Generated %v HTML files in the pages directory. Wrote %.1f kB in %v.", len(textFiles), htmlSize, endTime.Sub(startTime))
	} else if *fileFlag != "" {
		processTextFiles(*fileFlag)

		endTime := time.Now()
		
		htmlSize, _ := calcFileSizeInDirectory("pages")

		greenAndBold := chalk.Green.NewStyle().
			WithTextStyle(chalk.Bold).
			Style
		
		// This style didn't work properly in my terminal so I just used vanilla escape codes instead.
		// bold := chalk.Bold.NewStyle().Style

		// This is supposed to be green but may appear differently in your terminal depending on theme settings, it seems.
		fmt.Printf(greenAndBold("\nSuccess!") + " Generated \033[1m1\033[0m HTML files in the pages directory.  Wrote %.1f kB in %v.\n", htmlSize, endTime.Sub(startTime))
		mlog.Info("Success! Generated 1 HTML file in the pages directory. Wrote %.1f kB in %v.", htmlSize, endTime.Sub(startTime))
	}
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
			mlog.Info("Found text file: %v", path)

			mlog.Info("File name: %v", info.Name())

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

	// fmt.Println("File contents: ", string(fileContents))

	title, body := parseContent(fileContents)

	page := createPage(title, body)

	err = processTemplate("template.tmpl", page)
	if err != nil {
		mlog.Warning("Error processing template template.tmpl")
		mlog.Error(err)
		panic(err)
	}

	err = createHTMLFile(file, page)
	if err != nil {
		mlog.Warning("Error creating HTML file: %v", file)
		mlog.Error(err)
		panic(err)
	}
}

func readFile(fileName string) ([]byte, error) {
	return ioutil.ReadFile(fileName)
}

func parseContent(contents []byte) (string, string) {
	// TODO: fix the parsing so that the first line is the title and the rest is the body
	splitContents := strings.Split(string(contents), "\n")
	title := string(splitContents[0])
	body := strings.TrimSpace(strings.Join(splitContents[1:], "\n"))
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
		mlog.Warning("Error parsing template: %v", templateName)
		mlog.Error(err)
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
		mlog.Warning("Error creating new HTML file %v", fileName)
		mlog.Error(err)
		return err
	}
	defer newFile.Close()

	tmpl, err := template.ParseFiles("template.tmpl")
	if err != nil {
		mlog.Warning("Error parsing template: template.tmpl")
		mlog.Error(err)
		return err
	}
	return tmpl.Execute(newFile, page)
}

func calcFileSizeInDirectory(directory string) (float64, error) {
	var size float64

	err := filepath.Walk(directory, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			mlog.Warning("Error walking directory %v to calculate file sizes: %v", directory)
			mlog.Error(err)
			return err
		}

		size += float64(info.Size())

		return nil
	})

	if err != nil {
		mlog.Warning("Error calculating file size in directory: %v", directory)
		mlog.Error(err)
		return 0, err
	}

	size = size / 1024 // convert from bits to KB

	return size, nil
}

