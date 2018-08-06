package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gala377/Lyricer/lyrics"

	"github.com/gala377/Lyricer/config"
	"github.com/gala377/Lyricer/spotify"
)

// TODO refresh token after given time elapsed

func main() {
	// App init
	conf, err := config.Read("hidden_conf.json")
	if err != nil {
		log.Fatalf("Error while reading the config file: %s", err)
		return
	}
	spotify := spotify.NewSpotify(conf.Spotify)

	spotCodeChan, err := spotify.Authorize()
	if err != nil {
		log.Fatalf("Could not authorize spotify %s", err)
		return
	}
	// todo we can do some work while wainting for the authorization
	log.Println("Main Reading code")
	log.Printf("Authorized spotify. Code is: %s\n", <-spotCodeChan)

	log.Println("Main accessing spotify")
	accessChan, err := spotify.Access()
	if err != nil {
		log.Fatalf("Could not access spotify %s", err)
		return
	}
	code, ok := <-accessChan
	if !ok {
		log.Fatalf("Access channel closed early")
	}
	log.Printf("Accessed spotify. Access token is: %s\n", code)
	f := lyrics.TekstowoFetcher{}

	currPlaying, err := spotify.CurrentlyPlayedSong()
	if err != nil {
		log.Fatalf("Couldn't retrieve currently played song %s", err)
		return
	}
	song := lyrics.SongInfo{
		Author: currPlaying.Artist,
		Title:  currPlaying.Title,
	}
	err = song.FetchLyrics(f)
	if err != nil {
		log.Fatalf(
			"Could not fetch lyrics for the: %s, %s\n reason: %s\n",
			song.Author,
			song.Title,
			err,
		)
	}
	log.Printf("Song: %s, %s\n\n %s\n", song.Author, song.Title, song.Lyrics)

	refreshChannel := make(chan bool)
	closeChannel := make(chan bool)

	go func() {
	mainLoop:
		for {
			log.Printf("Waiting %d for the song end", currPlaying.Left)
			select {
			case <-time.After(currPlaying.Left):
			case <-refreshChannel:
			case <-closeChannel:
				break mainLoop
			}
			currPlaying, err = spotify.CurrentlyPlayedSong()
			if err != nil {
				log.Printf("Couldn't retrieve currently played song %s", err)
				log.Println("Trying again in 30 seconds")
				currPlaying.Left = time.Second * 30
			} else {
				song = lyrics.SongInfo{
					Author: currPlaying.Artist,
					Title:  currPlaying.Title,
				}
				err = song.FetchLyrics(f)
				if err != nil {
					log.Printf(
						"Could not fetch lyrics for the: %s, %s\n reason: %s\n",
						song.Author,
						song.Title,
						err,
					)
				}
				log.Printf("Song: %s, %s\n\n %s\n", song.Author, song.Title, song.Lyrics)
			}
		}
		closeChannel <- true
	}()

	for {
		fmt.Println("Q to quit, R to refresh")
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		if text == "r\n" {
			refreshChannel <- true
		} else if text == "q\n" {
			closeChannel <- true
			<-closeChannel
			return
		}
	}
}
