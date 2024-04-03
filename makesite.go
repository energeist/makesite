package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"html/template"
	"os"
)

// Page holds all the information we need to generate a new
// HTML page from a text file on the filesystem.
type Page struct {
	Title				 string
	Body				 string	
}

func main() {
	fmt.Println("Hello, World!")

	// Read in the contents of the provided first-post.txt file
	fileContents, err := ioutil.ReadFile("first-post.txt")
	if err != nil {
		panic(err)
	}
	fmt.Print(string(fileContents))

	// split the file contents into two parts: the title and the body using slices
	title := string(fileContents[0:16])
	body := string(fileContents[16:])

	// remove any leading or trailing white space from the title and body using the "strings" package
	body = strings.TrimSpace(body)

	// print the title and body to stdout
	fmt.Println("Title: ", title)
	fmt.Println("Body: ", body)

	// Create a new page struct and populate it with the title and body
	page := Page{
		Title:      string(title),
		Body: 		 	string(body),
	}

	// Create a new template in memory named "template.tmpl"
	// When the template is executed, it will parse template.tmpl,
	// looking for {{ }} where we can inject content.

	tmpl := template.Must(template.New("template.tmpl").ParseFiles("template.tmpl"))
	if err != nil {
		panic(err)
	}

	// Execute the template with the contents of first-post.txt and write to stdout
	err = tmpl.Execute(os.Stdout, page)
	if err != nil {
		panic(err)
	}

	// Create a new, blank HTML file named first-post.html.
	newFile, err := os.Create("first-post.html")
	if err != nil {
				panic(err)
	}

	// Executing the template injects the Page instance's data,
	// allowing us to render the content of our text file.
	// Furthermore, upon execution, the rendered template will be
	// saved inside the new file we created earlier.
	tmpl.Execute(newFile, page)
}
