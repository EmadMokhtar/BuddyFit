package main

import (
	"context"
	"fmt"
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
			author, title := fileName[0], fileName[1]
			title = strings.TrimSuffix(title, ".en.srt.txt")
			// Read the content of the file
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			transcript := string(content)
			// Insert the content into the table
			_, err = conn.Exec(ctx, "INSERT INTO yt_videos (author, title, transcript) VALUES ($1, $2, $3)", author, title, transcript)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Unable to insert data: %v\n", err)
				os.Exit(1)
			}
			fmt.Println(path)
		}
		return nil
	})

}
