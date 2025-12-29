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
			artist := "Unknown"
			album := "Unknown"

			tags, _ := taglib.ReadTags(path)
			//fmt.Printf("tags: %v\n", tags)
			fmt.Printf("AlbumArtist: %q\n", tags[taglib.AlbumArtist])
			artistSlice := tags[taglib.AlbumArtist]

			if len(artistSlice) != 0 {
				fmt.Println("AlbumArtist: " + artistSlice[0])
				artist = artistSlice[0]
			}

			albumSlice := tags[taglib.Album]

			if len(albumSlice) != 0 {
				fmt.Println("Album: " + albumSlice[0])
				album = albumSlice[0]
			}

			//fmt.Printf("Album: %q\n", tags[taglib.Album])
			//albumSlice := tags[taglib.Album]
			//fmt.Printf("TrackNumber: %q\n", tags[taglib.TrackNumber])

			embeddedImg, _ := taglib.ReadImage(path)
			if len(embeddedImg) != 0 && artist != "Unknown" && album != "Unknown" {
				directory := "albumart/" + artist + "/"
				filePath := directory + album + ".png"
				dirErr := os.MkdirAll(directory, os.ModePerm)
				if dirErr != nil {
					fmt.Println("Error creating directory for album art")
				}
				err := os.WriteFile(filePath, embeddedImg, os.ModePerm)
				if err != nil {
					fmt.Printf("Error writing embedded album image: %v\n", err)
				}
			}
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
