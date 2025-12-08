package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go.senan.xyz/taglib"
)

func main() {
	start := time.Now()

	rootPath := ""
	if len(os.Args) > 1 {
		rootPath = os.Args[1]
	} else {
		fmt.Println("No arguments passed. Need filepath for music directory. Exiting...")
	}

	count := 0

	err := filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if strings.HasSuffix(path, ".mp3") {
			count += 1
			tags, _ := taglib.ReadTags(path)
			fmt.Printf("tags: %v\n", tags)
			fmt.Printf("AlbumArtist: %q\n", tags[taglib.AlbumArtist])
			fmt.Printf("Album: %q\n", tags[taglib.Album])
			fmt.Printf("TrackNumber: %q\n", tags[taglib.TrackNumber])
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Error walking the directory: %v", err)
	}

	elapsed := time.Since(start)
	fmt.Printf("directory audio checks took %s.\n", elapsed)

	fmt.Printf("count was %f.\n", count)
}
