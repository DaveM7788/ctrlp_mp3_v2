package main

import (
	"ctrlpmp3v2/pkg/utils"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"slices"
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
	albumName    string
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

		if utils.IsAudioFileType(path) {
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
				songTitle = filepath.Base(path)
			}

			artistSlice := tags[taglib.AlbumArtist]
			if len(artistSlice) != 0 {
				artist = artistSlice[0]
			}

			albumSlice := tags[taglib.Album]
			if len(albumSlice) != 0 {
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
				albumName:    album,
				albumArtPath: "assets/images/controlp_sq_jpg.jpg",
			}

			embeddedImg, _ := taglib.ReadImage(path)
			if len(embeddedImg) != 0 && artist != "Unknown" && album != "Unknown" {
				directory := "albumart/" + artist + "/"
				filePath := directory + album + ".png"
				dirErr := os.MkdirAll(directory, os.ModePerm)
				if dirErr != nil {
					fmt.Printf("Error creating directory for album art: %v\n", dirErr)
				}
				err := os.WriteFile(filePath, embeddedImg, os.ModePerm)
				if err != nil {
					fmt.Printf("Error writing embedded album image: %v\n", err)
				}
				audioFile.albumArtPath = filePath
				albumInfo.albumArtPath = filePath
			}

			AudioFiles = append(AudioFiles, audioFile)
			if !slices.Contains(UniqueArtists, artist) {
				UniqueArtists = append(UniqueArtists, artist)
			}
			for i := range UniqueAlbums {
				if UniqueAlbums[i].albumName == album {
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
	fmt.Printf("Found %d audio files.\n", count)
}

func generateAllAlbumsPage() {
	for i := range UniqueAlbums {
		album := UniqueAlbums[i]
		fmt.Println("album name" + album.albumName)
	}
}

func generateAllArtistsPage() {
	for i := range UniqueArtists {
		artist := UniqueArtists[i]
		fmt.Println("artist" + artist)
	}
}

func generateAllSongsPage() {

}
