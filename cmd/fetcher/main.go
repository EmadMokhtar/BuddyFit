package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
)

var DataDir = os.Getenv("BF_DATA_DIR")

func main() {
	// Lead the csv file and get the youtube video urls
	csvFilePath := flag.String("file", "", "Path to the CSV file containing YouTube video URLs")
	includingHeaders := flag.Bool("headers", true, "Indicate whether the CSV file include headers or not. Default: Yes")
	flag.Parse()

	if *csvFilePath == "" {
		fmt.Printf("CSV file path is required\n")
		return
	}

	file, err := os.Open(*csvFilePath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)

	if *includingHeaders {
		// Read and discard the first record (header)
		if _, err := reader.Read(); err != nil {
			fmt.Println("Error:", err)
			return
		}
	}

	// Read all records from the CSV file
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var filePaths []string
	var wgDl sync.WaitGroup
	// Download the subtitles from the youtube videos
	for _, record := range records {
		vidTitle := record[0]
		vidURL := record[1]
		vidAuthor := record[2]
		wgDl.Add(1)
		go func(vidURL, title, author string) {
			defer wgDl.Done()
			filePath := DownloadSubtitles(vidURL, title, author, DataDir)
			filePaths = append(filePaths, filePath)
		}(vidURL, vidTitle, vidAuthor)
	}

	wgDl.Wait()

	var wgConv sync.WaitGroup
	// Convert the subtitles to text
	for _, filePath := range filePaths {
		wgConv.Add(1)
		fpWithExt := fmt.Sprintf("%s/%s.en.srt", DataDir, filePath)
		go func(filePath string) {
			defer wgConv.Done()
			if filePath == "" {
				return
			}
			processSrtFile(filePath, filePath+".txt")
			err := os.Remove(filePath)
			if err != nil {
				fmt.Printf("Removing File Error: %v\n", err)
			}
		}(fpWithExt)
	}

	wgConv.Wait()

}

func DownloadSubtitles(vidURL, title, author, outputDir string) string {
	binPath := getytDlPath()
	outputFile := fmt.Sprintf("%s==%s==%s==", author, formatOutputFile(title), url.QueryEscape(vidURL))

	cmd := fmt.Sprintf("%s --write-auto-subs --convert-subs srt --skip-download --sub-lang en -o '%s' %s", binPath, outputFile, vidURL)
	fmt.Printf("Running command: %s\n", cmd)
	dlDlp := exec.Command("sh", "-c", cmd)
	dlDlp.Dir = outputDir
	err := dlDlp.Run()
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}
	return outputFile
}

func getytDlPath() string {
	ytDlPath, err := exec.LookPath("yt-dlp")
	if err != nil {
		fmt.Println("yt-dlp not found in PATH")
		return ""
	}
	return ytDlPath
}

func formatOutputFile(input string) string {
	// Convert to lowercase
	result := strings.ToLower(input)
	// Replace spaces with underscores
	result = strings.ReplaceAll(result, "_", "")
	// Replace spaces with underscores
	result = strings.ReplaceAll(result, " ", "_")
	// Remove slashes
	result = strings.ReplaceAll(result, "/", "")
	// Remove colons
	result = strings.ReplaceAll(result, ":", "")
	// Remove question marks
	result = strings.ReplaceAll(result, "?", "")
	// Remove exclamation marks
	result = strings.ReplaceAll(result, "!", "")
	// Remove commas
	result = strings.ReplaceAll(result, ",", "")
	// Remove periods
	result = strings.ReplaceAll(result, ".", "")
	// Remove parentheses
	result = strings.ReplaceAll(result, "(", "")
	result = strings.ReplaceAll(result, ")", "")
	// Remove brackets
	result = strings.ReplaceAll(result, "[", "")
	result = strings.ReplaceAll(result, "]", "")
	// Remove curly braces
	result = strings.ReplaceAll(result, "{", "")
	result = strings.ReplaceAll(result, "}", "")
	// Remove ampersands
	result = strings.ReplaceAll(result, "&", "and")
	// Remove | signs
	result = strings.ReplaceAll(result, "|", "")
	// Remove single quotes signs
	result = strings.ReplaceAll(result, "'", "")
	// Remove double quotes signs
	result = strings.ReplaceAll(result, "\"", "")
	// Remove backticks
	result = strings.ReplaceAll(result, "`", "")
	return result
}

func processSrtFile(inputFilename, outputFilename string) {
	if inputFilename == "" {
		return
	}
	inputFile, err := os.Open(inputFilename)
	if err != nil {
		fmt.Printf("Error: Input file '%s' not found.\n", inputFilename)
		return
	}
	defer inputFile.Close()

	var processedLines []string
	scanner := bufio.NewScanner(inputFile)
	reNumber := regexp.MustCompile(`^[0-9]+$`)
	reTimestamp := regexp.MustCompile(`^[0-9]{2}:[0-9]{2}:[0-9]{2}`)
	for scanner.Scan() {
		line := scanner.Text()
		if !reNumber.MatchString(line) && !reTimestamp.MatchString(line) && line != "" {
			if line != "" && !contains(processedLines, line) {
				processedLines = append(processedLines, line)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("An error occurred: %v\n", err)
		return
	}

	outputFile, err := os.Create(outputFilename)
	if err != nil {
		fmt.Printf("An error occurred: %v\n", err)
		return
	}
	defer outputFile.Close()

	combinedLines := strings.Join(processedLines, " ")
	_, err = outputFile.WriteString(combinedLines)
	if err != nil {
		fmt.Printf("An error occurred: %v\n", err)
		return
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
