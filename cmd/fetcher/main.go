package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

type ChannelDetails struct {
	ChannelID   string `json:"id"`
	ChannelName string `json:"name"`
	ChannelURL  string `json:"url"`
}

type Config struct {
	Channels      []ChannelDetails `json:"channels"`
	OutputDir     string           `json:"output_dir"`
	ArchiveFile   string           `json:"archive_file"`
	CheckInterval int              `json:"check_interval_hours"`
	Languages     []string         `json:"languages"`
	LastRunTime   int64            `json:"last_run_time"`
}

func main() {
	// Parse command-line flags
	configFile := flag.String("config", "config.json", "Path to configuration file")
	runOnce := flag.Bool("once", false, "Run once without scheduling")
	noDownload := flag.Bool("no-download", false, "Run without downloading subtitles")
	cookies := flag.String("cookies", "", "Path to cookies file")

	flag.Parse()

	// Load configuration
	config, err := loadConfig(*configFile)
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	// Ensure output directory exists
	if err := os.MkdirAll(config.OutputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Create archive file if it doesn't exist
	if _, err := os.Stat(config.ArchiveFile); os.IsNotExist(err) {
		file, err := os.Create(config.ArchiveFile)
		if err != nil {
			log.Fatalf("Failed to create archive file: %v", err)
		}
		err = file.Close()
		if err != nil {
			return
		}
	}

	if !*noDownload {
		// Run the download now
		downloadSubtitles(config, *cookies)
	}

	// Update last run time
	config.LastRunTime = time.Now().Unix()
	err = saveConfig(config, *configFile)
	if err != nil {
		log.Fatalf("Error saving configuration: %v", err)
	}

	archiveFile, err := os.OpenFile(config.ArchiveFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error opening archive file: %v", err)
	}
	if archiveFile == nil {
		log.Fatalf("Archive file is nil")
	}
	defer func(archiveFile *os.File) {
		err := archiveFile.Close()
		if err != nil {
			return
		}
	}(archiveFile)
	// Update the archive file with the downloaded subtitles
	// Clean up subtitles
	var wgConv sync.WaitGroup

	for _, channel := range config.Channels {
		channelDir := filepath.Join(config.OutputDir, channel.ChannelName)

		channelSubtitleFiles, err := os.ReadDir(channelDir)
		if err != nil {
			log.Fatalf("Error reading directory %s: %v", channelDir, err)
		}

		for _, channelSubtitleFile := range channelSubtitleFiles {
			if channelSubtitleFile.IsDir() {
				continue
			}

			wgConv.Add(1)
			go func(fileInfo os.DirEntry, channelDir string) {
				defer wgConv.Done()
				subtitleFilePath := filepath.Join(channelDir, fileInfo.Name())
				fmt.Printf("Cleaning subtitle file: %s\n", subtitleFilePath)
				outputPath := subtitleFilePath
				if !strings.HasSuffix(outputPath, ".txt") {
					outputPath += ".txt"
				}
				processSubtitleFile(subtitleFilePath, outputPath)
				fmt.Printf("Done cleaning subtitle file: %s\n", subtitleFilePath)
				err := os.Remove(subtitleFilePath)
				fmt.Printf("Removing subtitle file: %s\n", subtitleFilePath)
				if err != nil {
					fmt.Printf("Removing File Error: %v\n", err)
				}
			}(channelSubtitleFile, channelDir)

			videoID := extractVideoID(channelSubtitleFile.Name())
			fmt.Printf("Appending video ID to archive file: %s\n", videoID)
			archiveLine := fmt.Sprintf("youtube %s\n", videoID)
			if _, err := archiveFile.WriteString(archiveLine); err != nil {
				log.Fatalf("Error writing to archive file: %v", err)
			}
			fmt.Printf("Done appending video ID to archive file: %s\n", videoID)

		}
	}

	wgConv.Wait()

	// If not run once, schedule periodic runs
	if !*runOnce {
		fmt.Printf("Scheduled to check for new videos every %d hours\n", config.CheckInterval)
		ticker := time.NewTicker(time.Duration(config.CheckInterval) * time.Hour)
		defer ticker.Stop()

		for range ticker.C {
			downloadSubtitles(config, "")
			config.LastRunTime = time.Now().Unix()
			err = saveConfig(config, *configFile)
			if err != nil {
				log.Fatalf("Error saving configuration: %v", err)
			}
		}
	}
}

// loadConfig reads the configuration from a file
func loadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			return
		}
	}(file)

	config := &Config{}
	if err := json.NewDecoder(file).Decode(config); err != nil {
		return nil, fmt.Errorf("error parsing config: %w", err)
	}
	return config, nil
}

// saveConfig writes the configuration to a file
func saveConfig(config *Config, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			return
		}
	}(file)

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(config)
}

// downloadSubtitles downloads subtitles from all configured channels
func downloadSubtitles(config *Config, cookies string) {
	fmt.Println("Starting subtitle download process...")
	currentTime := time.Now().UTC().Format("2006-01-02 15:04:05")
	fmt.Printf("Current UTC time: %s\n", currentTime)

	// Process each channel
	for _, channel := range config.Channels {
		fmt.Printf("Processing channel: %s\n", channel)

		// Create channel-specific directory
		channelDir := filepath.Join(config.OutputDir, channel.ChannelName)
		if err := os.MkdirAll(channelDir, 0755); err != nil {
			log.Printf("Error creating directory for channel %s: %v", channel.ChannelName, err)
			continue
		}

		// Prepare languages
		languageOpts := strings.Join(config.Languages, ",")

		// Using the specific format requested: %(title)s-%(id)s.%(ext)s
		outputTemplate := filepath.Join(channelDir, "%(id)s||%(title)s.%(ext)s")

		// Prepare command: yt-dlp with options for subtitle download only
		args := []string{
			"--skip-download",          // Don't download the video
			"--write-auto-sub",         // Download auto-generated subtitles
			"--sub-format", "srt/best", // Prefer SRT format
			"--convert-subs", "srt", // Convert subtitles to SRT
			"--sub-langs", languageOpts, // Languages to download
			"--download-archive", config.ArchiveFile, // Track downloaded videos
			"--output", outputTemplate, // User-specified output file naming
		}

		// Add cookies option only if cookies string is not empty
		if cookies != "" {
			args = append(args, "--cookies", cookies)
		}

		// Add channel URL at the end
		args = append(args, channel.ChannelURL)

		// Execute yt-dlp command
		cmd := exec.Command("yt-dlp", args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		fmt.Printf("Running: yt-dlp %s\n", strings.Join(args, " "))
		if err := cmd.Run(); err != nil {
			log.Printf("Error downloading subtitles from %s: %v", channel, err)
			continue
		}

		fmt.Printf("Finished processing channel: %s\n", channel)
	}

	fmt.Println("Subtitle download process completed.")
}

func extractVideoID(fileName string) string {
	videoFilename := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	parts := strings.Split(videoFilename, "||")
	if len(parts) > 0 {
		return strings.TrimSpace(parts[0])
	}
	return videoFilename
}

func processSubtitleFile(inputFilename, outputFilename string) {
	if inputFilename == "" {
		return
	}
	inputFile, err := os.Open(inputFilename)
	if err != nil {
		fmt.Printf("Error: Input file '%s' not found.\n", inputFilename)
		return
	}
	defer func(inputFile *os.File) {
		err := inputFile.Close()
		if err != nil {
			fmt.Printf("An error occurred: %v\n", err)
		}
	}(inputFile)

	var processedLines []string
	scanner := bufio.NewScanner(inputFile)
	reNumber := regexp.MustCompile(`^[0-9]+$`)
	reTimestamp := regexp.MustCompile(`^[0-9]{2}:[0-9]{2}:[0-9]{2}`)
	for scanner.Scan() {
		line := scanner.Text()
		if !reNumber.MatchString(line) && !reTimestamp.MatchString(line) && line != "" {
			if !contains(processedLines, line) {
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
	defer func(outputFile *os.File) {
		err := outputFile.Close()
		if err != nil {
			fmt.Printf("An error occurred: %v\n", err)
		}
	}(outputFile)

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
