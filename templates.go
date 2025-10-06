package main

import (
	"io/ioutil"
	"log"
)

var indexHTML string

func loadTemplates() {
	// Load the HTML template
	htmlBytes, err := ioutil.ReadFile("templates/index.html")
	if err != nil {
		log.Fatal("Error loading HTML template:", err)
	}
	indexHTML = string(htmlBytes)
}
