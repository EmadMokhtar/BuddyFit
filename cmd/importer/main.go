package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/jackc/pgx/v5"
)

const sep = "=="

func main() {
	dsn := os.Getenv("BF_DB_URL")
	dataDir := os.Getenv("BF_DATA_DIR")
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(ctx)

	// Read all the files in the data directory and insert them into the database
	err = filepath.Walk(dataDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			if filepath.Ext(info.Name()) != ".txt" {
				return nil
			}
			fileName := strings.Split(info.Name(), sep)
			author, title, vidURL := fileName[0], fileName[1], fileName[2]
			// TODO: Remove _ from the title
			title = strings.ReplaceAll(title, "_", " ")
			// TODO: Remove the .en.srt.txt from the url
			vidURL = strings.TrimSuffix(vidURL, ".en.srt.txt")
			// Read the content of the file
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			transcript := string(content)
			// Insert the content into the table
			decodedVidURL, err := url.QueryUnescape(vidURL)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Unable to decode the URL: %v\n", err)
				os.Exit(1)
			}
			_, err = conn.Exec(ctx, "INSERT INTO yt_videos (author, title, transcript, url) VALUES ($1, $2, $3, $4)", author, title, transcript, decodedVidURL)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Unable to insert data: %v\n", err)
				os.Exit(1)
			}
			fmt.Println(path)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Error while looping over files in the directory, %s", err)
	}

}
