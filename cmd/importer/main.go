package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	pgx "github.com/jackc/pgx/v5"

	"github.com/EmadMokhtar/BuddyFit/internal"
)

func main() {
	dsn := os.Getenv("BF_DB_URL")
	dataDir := os.Getenv("BF_DATA_DIR")
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		_, err := fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		if err != nil {
			return
		}
		os.Exit(1)
	}
	defer func(conn *pgx.Conn, ctx context.Context) {
		err = conn.Close(ctx)
		if err != nil {

		}
	}(conn, ctx)

	// Read all the files in the data directory and insert them into the database
	err = filepath.Walk(dataDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			if filepath.Ext(info.Name()) != ".txt" {
				return nil
			}
			// Example: 0FE_0zjjaQs||The Truth About Genetics and YOUR Muscle Growth.en.srt.txt
			// Result should be: Video ID: 0FE_0zjjaQs, Title: The Truth About Genetics and YOUR Muscle Growth, Author: file directory name
			video, err := internal.NewVideoFromFileName(info.Name(), path)
			if err != nil {
				_, err = fmt.Fprintf(os.Stderr, "Unable to create video from file name: %v\n", err)
				if err != nil {
					return err
				}
				return nil // Continue processing other files
			}
			_, err = conn.Exec(ctx, "INSERT INTO yt_videos (author, title, transcript, url) VALUES ($1, $2, $3, $4)", video.Author, video.Title, video.Transcript, video.UnescapedURL())
			if err != nil {
				_, err = fmt.Fprintf(os.Stderr, "Unable to insert data: %v\n", err)
				if err != nil {
					return err
				}
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
