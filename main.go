package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"go.senan.xyz/taglib"
)

type AudioFile struct {
	filePath     string
	songName     string
	artist       string
	album        string
	albumArtPath string
}

type AudioAlbum struct {
	album        string
	albumArtPath string
}

var AudioFiles = []AudioFile{}
var UniqueArtists = []string{}
var UniqueAlbums = []AudioAlbum{}

func main() {
	start := time.Now()

	rootPath := ""
	if len(os.Args) > 1 {
		rootPath = os.Args[1]
	} else {
		fmt.Println("No arguments passed. Need filepath for music directory. Exiting...")
	}

	collectAudioFiles(rootPath)

	elapsed := time.Since(start)
	fmt.Printf("Static site generation time: %s.\n", elapsed)
}

func collectAudioFiles(rootPath string) {
	count := 0
	err := filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if strings.HasSuffix(path, ".mp3") {
			count += 1
			songTitle := "Unknown"
			artist := "Unknown"
			album := "Unknown"

			tags, _ := taglib.ReadTags(path)
			//fmt.Printf("tags: %v\n", tags)
			//fmt.Printf("AlbumArtist: %q\n", tags[taglib.AlbumArtist])

			songNameSlice := tags[taglib.Title]
			if len(songNameSlice) != 0 {
				songTitle = songNameSlice[0]
			} else {
				// take file name
			}

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

			audioFile := AudioFile{
				filePath:     path,
				songName:     songTitle,
				artist:       artist,
				album:        album,
				albumArtPath: "assets/images/controlp_sq_jpg.jpg",
			}

			albumInfo := AudioAlbum{
				album:        album,
				albumArtPath: "assets/images/controlp_sq_jpg.jpg",
			}

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
				audioFile.albumArtPath = filePath
				albumInfo.albumArtPath = filePath
			}

			AudioFiles = append(AudioFiles, audioFile)
			fmt.Println("Added the following song name " + AudioFiles[0].songName)

			if !slices.Contains(UniqueArtists, artist) {
				UniqueArtists = append(UniqueArtists, artist)
			}

			for i := range UniqueAlbums {
				if UniqueAlbums[i].album == album {
					continue
				}
				UniqueAlbums = append(UniqueAlbums, albumInfo)
			}
		}

		generateAllAlbumsPage()
		generateAllArtistsPage()

		return nil
	})

	if err != nil {
		log.Fatalf("Error walking the directory: %v", err)
	}
	fmt.Printf("Found %f audio files.\n", count)
}

func generateAllAlbumsPage() {
	for i := range UniqueAlbums {
		albumName := UniqueAlbums[i]
		fmt.Println("album name" + albumName.album)
	}
}

func generateAllArtistsPage() {

}

func generateAllSongsPage() {

}
