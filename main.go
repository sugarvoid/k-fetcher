package main

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

// removes invalid characters
func sanitize_filename(filename string) string {
	// Replace spaces with underscores and remove invalid characters for filenames
	re := regexp.MustCompile(`[<>:"/\|?*]`) // characters not allowed in filenames
	return re.ReplaceAllString(filename, "_")
}

func download_file(url, filename string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error fetching the file: %v", err)
	}
	defer resp.Body.Close()

	// Check if Content-Disposition header is present to suggest a filename
	// If it's not present, we will fall back on generating a filename from the URL
	// contentDisposition := resp.Header.Get("Content-Disposition")
	// if contentDisposition != "" {
	// 	re := regexp.MustCompile(`filename="(.+?)"`)
	// 	matches := re.FindStringSubmatch(contentDisposition)
	// 	if len(matches) > 1 {
	// 		filename = matches[1]
	// 	}
	// }

	if filename == "" {
		fmt.Errorf("missing file name")
		os.Exit(1)
	}

	// Sanitize the filename to ensure it's safe to use on all filesystems
	filename = sanitize_filename(filename)

	// Create the file to save the downloaded content
	outFile, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer outFile.Close()

	// Copy the response body (file content) to the new file
	_, err = outFile.ReadFrom(resp.Body)
	if err != nil {
		return fmt.Errorf("error saving the file: %v", err)
	}

	fmt.Printf("Downloaded %s successfully.\n", filename)
	return nil
}

func main() {
	// Check if any arguments are passed (i.e., the file path)
	if len(os.Args) < 2 {
		fmt.Println("Please drag and drop a CSV file onto the executable.")
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

		// Create a new CSV reader
		reader := csv.NewReader(file)
		// Read all records (lines) from the CSV file
		records, err := reader.ReadAll()
		if err != nil {
			fmt.Printf("Error reading CSV file: %v\n", err)
			return
		}

		// Check if the file is empty
		if len(records) == 0 {
			fmt.Println("CSV file is empty.")
			return
		}

		// Get the headers (first row) to find the "DOWNLOAD" and "ITEM_NAME" column indices
		headers := records[0]
		downloadColumnIndex := -1
		titleColumnIndex := -1

		for i, header := range headers {
			if strings.ToUpper(header) == "DOWNLOAD" {
				downloadColumnIndex = i
			}
			if strings.ToUpper(header) == "ITEM_NAME" {
				titleColumnIndex = i
			}
		}

		// check for the needed headers
		if downloadColumnIndex == -1 {
			fmt.Println("No 'DOWNLOAD' column found in the CSV file.")
			return
		}
		if titleColumnIndex == -1 {
			fmt.Println("No 'ITEM_NAME' column found in the CSV file.")
			return
		}

		// Loop through each record and download the file from the "DOWNLOAD" column
		for lineNumber, record := range records[1:] { // Skip the header row
			if downloadColumnIndex < len(record) && titleColumnIndex < len(record) {
				downloadURL := record[downloadColumnIndex]
				filename := record[titleColumnIndex]
				filename = sanitize_filename(filename) + ".mp4"

				// print the download information
				fmt.Printf("Line %d: Downloading file from URL: %s\n", lineNumber+1, downloadURL)

				err := download_file(downloadURL, filename)
				if err != nil {
					fmt.Printf("Error downloading file: %v\n", err)
				}

				// adds delay between downloads
				time.Sleep(5 * time.Second)
			}
		}
	}
}
