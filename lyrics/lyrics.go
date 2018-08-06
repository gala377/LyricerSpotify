// Package lyrics defines interfaces and strucs
// for handling lyrics fething.
package lyrics

import "errors"

// Fetcher should provide lyrics
// given the songs artist and it's title.
type Fetcher interface {
	FetchLyrics(author, title string) (string, error)
}

// ErrLyricsNotFound is a generic Fetcher error.
// It doesn't provide much information beside that
// the lyrics for the given song could not be fetched.
var ErrLyricsNotFound = errors.New("lyrics for the requested song couldn't be found")

// SongInfo represents basic song information
// needed for the Lyricer app to print to the
// console.
type SongInfo struct {
	Author string
	Title  string
	Lyrics string
}

// FetchLyrics sets the SongInfos Lyrics field
// to the return value of the Fether FetchLyrics call.
//
// FetchLyrics arguments are taken from the current
// values of the SongInfos Author and Title fields.
//
// Example:
//  s := SongInfo{Author: "Me", Title: "MySong"}
//  f := ExampleLyricsFetcher{}
//  err := s.FetchLyrics(&f)
//  if err != nil {
//   log.Fatalf("Could not fetch my song lyrics: %s", err)
//  }
//  log.Printf("My songs lyrics are: %s", s.Lyrics)
func (s *SongInfo) FetchLyrics(f Fetcher) error {
	var err error
	s.Lyrics, err = f.FetchLyrics(s.Author, s.Title)
	return err
}
