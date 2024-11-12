package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// Check if any arguments are passed (i.e., the file path)
	if len(os.Args) < 2 {
		fmt.Println("Please drag and drop a video list text file onto the executable.")
		return
	}

	// Loop through the arguments passed (which will be the file paths)
	for _, filePath := range os.Args[1:] {
		fmt.Println("Processing file:", filePath)

		// Open the file
		file, err := os.Open(filePath)
		if err != nil {
			fmt.Printf("Error opening file: %v\n", err)
			return
		}
		defer file.Close() // Ensure the file is closed after we're done

		// Create a new scanner to read the file line by line
		scanner := bufio.NewScanner(file)
		lineNumber := 1
		for scanner.Scan() {
			// For each line, print the line number and the content
			fmt.Printf("Line %d: %s\n", lineNumber, scanner.Text())
			// Download the video from url in the current dir
			lineNumber++
		}

		// Check for errors during scanning
		if err := scanner.Err(); err != nil {
			fmt.Printf("Error reading file: %v\n", err)
		}
	}
}
