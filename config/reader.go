// Package config defines Lyricer apps configuration
// format as well as structs ready to
// unmarshall config file to.
// For that pakcage config provides simple
// reading function `Read`.
package config

import (
	"encoding/json"
	"io/ioutil"
)

// OAuthData represents complete data
// needed to make successful Authorization
// and Access OAuth requests with for
// the selected service provider.
type OAuthData struct {
	AuthURL     string
	AccessURL   string
	ClientID    string
	SecretID    string
	CallbackURL string
	Scopes      []string
}

// LyricerConfig is the data needed
// for the Lyricer app to successfuly
// access services providers (for now only spotify)
// web api.
type LyricerConfig struct {
	Spotify OAuthData
}

// Read opens and reads the configuration
// from the file given under the
// configFilePath argument. Then returns
// it as *LyricerConfig.
//
// Configuration file should be in json format
// with the structure following Go's json.Unmarshall
// logic.
func Read(configFilePath string) (*LyricerConfig, error) {
	data, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}
	var config LyricerConfig

	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, err
}
