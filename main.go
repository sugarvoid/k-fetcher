package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const VERSION string = "1.0.1"

// Removes invalid characters from filenames
func SanitizeFilename(filename string) string {
	re := regexp.MustCompile(`[<>:"/\|?*$&%]`) // bad characters
	return re.ReplaceAllString(filename, "")
}

func DownloadFile(url string, filename string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error fetching the file: %v", err)
	}
	defer resp.Body.Close()

	if filename == "" {
		fmt.Errorf("missing file name")
		os.Exit(1)
	}

	// Sanitize the filename to remove bad charaters
	filename = SanitizeFilename(filename)

	// Create the file to save the downloaded content
	outFile, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("Failed to create file: %v", err)
	}
	defer outFile.Close()

	// Copy the response body (file content) to the new file
	_, err = outFile.ReadFrom(resp.Body)
	if err != nil {
		return fmt.Errorf("Failed to save the file: %v", err)
	}

	fmt.Printf("Downloaded %s successfully.\n", filename)
	return nil
}

func main() {
	fmt.Printf("Video Downloader\nVersion: %s\nAuthor: Dave\nLicense: MIT\n\n", VERSION)
	// Check if for file being dragged in
	if len(os.Args) < 2 {
		fmt.Println("Please drag and drop a CSV file onto this executable.")
		return
	}
	if len(os.Args) > 2 {
		fmt.Println("One file at a time.")
		return
	}

	// Loop through the arguments passed (which will be the file paths)
	var filePath string = os.Args[1]

	var extension string = filepath.Ext(filePath)

	if extension != ".csv" {
		fmt.Println("Please use a cvs file.")
		fmt.Println("Press any key to exit...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		return
	}

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

	// Get the headers (first row) to find the "DOWNLOAD", "ITEM_NAME", and "ENTRY_ID" columns
	var headers []string = records[0]
	var downloadColumnIndex int = -1
	var titleColumnIndex int = -1
	var entryIdColumnIndex int = -1

	for i, header := range headers {
		if strings.ToUpper(header) == "DOWNLOAD" {
			downloadColumnIndex = i
		}
		if strings.ToUpper(header) == "ITEM_NAME" {
			titleColumnIndex = i
		}
		if strings.ToUpper(header) == "ENTRY_ID" {
			entryIdColumnIndex = i
		}
	}

	// Check for the needed headers
	if downloadColumnIndex == -1 {
		fmt.Println("No 'DOWNLOAD' column found in the CSV file.")
		return
	}
	if titleColumnIndex == -1 {
		fmt.Println("No 'ITEM_NAME' column found in the CSV file.")
		return
	}
	if entryIdColumnIndex == -1 {
		fmt.Println("No 'ENTRY_ID' column found in the CSV file.")
		return
	}

	// Loop through each record and download the file from the "DOWNLOAD" column
	for _, record := range records[1:] { // Skip the header row
		if downloadColumnIndex < len(record) && titleColumnIndex < len(record) && entryIdColumnIndex < len(record) {
			var downloadURL string = record[downloadColumnIndex]
			var itemName string = record[titleColumnIndex]
			var entryId string = record[entryIdColumnIndex]

			//fmt.Println("%s", entryId)

			// Construct the filename in the format: [ITEM_NAME](ENTRY_ID).mp4
			var filename string = fmt.Sprintf("%s(%s).mp4", itemName, entryId)
			filename = SanitizeFilename(filename) // Sanitize the filename

			// Print the download information
			//fmt.Printf("Line %d: Downloading file from URL: %s\n", lineNumber+1, downloadURL)

			err := DownloadFile(downloadURL, filename)
			if err != nil {
				fmt.Printf("Error downloading file: %v\n", err)
			}

			// Add delay between downloads
			time.Sleep(5 * time.Second)
		}
	}

}
