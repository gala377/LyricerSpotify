// Package spotify defines methods and structs for
// the Lyricer app to interact with the
// Spotify services api.
package spotify

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gala377/Lyricer/config"
	"github.com/gala377/Lyricer/oauth"
)

// ErrEmptySongData is an error returned by the
// CurrentlyPlayedSong method when the song
// data from the spotify couldn't be fetched.
var ErrEmptySongData = errors.New("Response returned empty value")

// Response is a json response
// returned by the spotify service
// upon succesful granting of the
// access token.
type Response struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    uint   `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

// CurrentlyPlayed represents data of the
// currently played song on the user spotify
// client.
type CurrentlyPlayed struct {
	// Artist of the song.
	Artist string
	// Title of the song.
	Title string
	// How much time is left till the song ends in milliseconds.
	Left time.Duration
	// Is the song playing or is it paused.
	IsPlaying bool
}

// Spotify handles OAuth communication with
// the Spotify service.
type Spotify struct {
	conf         config.OAuthData
	authCode     string
	accessToken  string
	refreshToken string
	expires      time.Time
}

// NewSpotify creates new Spotify struct
// from the given spotify configuration.
func NewSpotify(conf config.OAuthData) *Spotify {
	return &Spotify{
		conf: conf,
	}
}

// Authorize grants authorization from the user
// for the spotify service.
// Access request should follow successful authorization
// attempt.
func (s *Spotify) Authorize() (<-chan string, error) {
	log.Println("Creating authorization request")
	r, err := oauth.NewAuthRequest(
		s.conf.AuthURL,
		s.conf.ClientID,
		s.conf.CallbackURL,
		s.conf.Scopes,
	)
	if err != nil {
		return nil, err
	}
	log.Println("Creating channel to pass code through")
	codeChan := make(chan string)
	go func() {
		log.Println("Authorize() waiting for auth code")
		s.authCode = <-oauth.Authorize(r)
		log.Printf("Got: %s, passing\n", s.authCode)
		codeChan <- s.authCode
		log.Println("Passed")
	}()
	return codeChan, nil
}

// Access grants access token for the spotify service.
// Note that the authorization needs to be granted
// first.
func (s *Spotify) Access() (<-chan string, error) {
	log.Println("Getting spotify access token")
	r, err := oauth.NewAccessRequest(
		s.conf.AccessURL,
		s.conf.ClientID,
		s.conf.SecretID,
		s.authCode,
		s.conf.CallbackURL,
	)
	if err != nil {
		return nil, err
	}
	out := make(chan string)
	go func() {
		log.Println("Access() making spotify request")
		accessRespBody, err := oauth.Access(r)
		if err != nil {
			log.Fatalf("Could not retrieve access response body: %s", err)
			close(out)
			return
		}
		s.parseRespBody(accessRespBody)
		out <- s.accessToken
		close(out)
	}()
	return out, nil
}

func (s *Spotify) parseRespBody(respBody []byte) error {
	var response Response
	err := json.Unmarshal(respBody, &response)
	if err != nil {
		return err
	}
	s.accessToken = response.AccessToken
	s.refreshToken = response.RefreshToken
	s.expires = time.Now().Add(time.Second * time.Duration(response.ExpiresIn))
	return nil
}

// Refresh refreshes spotify services access token
// if possible.
func (s *Spotify) Refresh() (<-chan string, error) {
	panic("Spotify Refresh unimplemented")
}

// CurrentlyPlayedSong returns data of the song currently
// played in the users spotify client
func (s *Spotify) CurrentlyPlayedSong() (CurrentlyPlayed, error) {
	log.Println("Creating played song request")
	req, err := s.playedSongRequest()
	if err != nil {
		return CurrentlyPlayed{}, err
	}
	log.Println("Sending request")
	respBody, err := s.playedSongResponce(req)
	if err != nil {
		return CurrentlyPlayed{}, err
	}
	log.Println("Parsing to currently played object")
	return s.spotifyResponseToCurrentlyPlayed(respBody)

}

func (s *Spotify) playedSongRequest() (*http.Request, error) {
	req, err := http.NewRequest(
		"GET",
		"https://api.spotify.com/v1/me/player/currently-playing",
		nil,
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))
	return req, nil
}

func (s *Spotify) playedSongResponce(r *http.Request) ([]byte, error) {
	client := http.Client{}
	resp, err := client.Do(r)
	log.Println("Request send")
	if err != nil {
		log.Printf("Error returned: %s", err)
		return nil, err
	}
	defer resp.Body.Close()
	log.Println("Reading resp body")
	return ioutil.ReadAll(resp.Body)
}

func (s *Spotify) spotifyResponseToCurrentlyPlayed(respBody []byte) (CurrentlyPlayed, error) {
	var relevantRespInfo struct {
		ProgressMS int `json:"progress_ms"`
		Item       struct {
			Artists []struct {
				Name string `json:"name"`
			} `json:"artists"`
			SongTitle  string `json:"name"`
			DurationMS int    `json:"duration_ms"`
		}
		IsPlaying bool `json:"is_playing"`
	}
	err := json.Unmarshal(respBody, &relevantRespInfo)
	if err != nil {
		return CurrentlyPlayed{}, err
	}
	if len(relevantRespInfo.Item.Artists) == 0 && relevantRespInfo.Item.SongTitle == "" {
		return CurrentlyPlayed{}, ErrEmptySongData
	}

	timeLeft := time.Millisecond * time.Duration(relevantRespInfo.Item.DurationMS)
	log.Printf("Duration is: %d", timeLeft)
	log.Printf("Progress is: %d", relevantRespInfo.ProgressMS)
	timeLeft = timeLeft - time.Millisecond*time.Duration(relevantRespInfo.ProgressMS)
	log.Printf("Left to song end is: %d", timeLeft)
	return CurrentlyPlayed{
		Artist:    relevantRespInfo.Item.Artists[0].Name,
		Title:     relevantRespInfo.Item.SongTitle,
		Left:      timeLeft,
		IsPlaying: relevantRespInfo.IsPlaying,
	}, nil

}
