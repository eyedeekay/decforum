package main

import (
	"io/ioutil"
	"strings"
)

func DeduplicateLinesInFile(filename string) error {
	// Open the file
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	// Split the file into lines
	lines := strings.Split(string(bytes), "\n")
	// Create a map of lines
	lineMap := make(map[string]bool)
	// Loop through the lines
	for _, line := range lines {
		// If the line is not in the map, add it
		if _, ok := lineMap[line]; !ok {
			lineMap[line] = true
		}
	}
	// Create a new string
	newString := ""
	// Loop through the map
	for line := range lineMap {
		// Add the line to the new string
		newString += line + "\n"
	}
	newString = strings.TrimLeft(newString, "\n")
	newString = strings.TrimRight(newString, "\n")
	// Write the new string to the file
	err = ioutil.WriteFile(filename, []byte(newString), 0644)
	if err != nil {
		return err
	}
	return nil
}
