package oauth

import (
	"encoding/base64"
	"fmt"
	"net/url"
)

// AuthRequest represents data needed to make authorization OAuth request.
type AuthRequest struct {
	AuthURL     *url.URL
	ClientID    string
	Scopes      []string
	RedirectURL *url.URL
}

// NewAuthRequest parses URLs passed and reurns Request struct ready to make
// authorization request with.
func NewAuthRequest(authURL, clientID, redirectURL string, scopes []string) (*AuthRequest, error) {
	auth, err := url.Parse(authURL)
	if err != nil {
		return nil, err
	}
	redirect, err := url.Parse(redirectURL)
	if err != nil {
		return nil, err
	}
	return &AuthRequest{
		AuthURL:     auth,
		ClientID:    clientID,
		Scopes:      scopes,
		RedirectURL: redirect,
	}, nil
}

// URL returns authorization url parsed from the request data.
func (r *AuthRequest) URL() *url.URL {
	v := url.Values{}
	v.Set("response_type", "code")
	v.Set("client_id", r.ClientID)
	v.Set("redirect_uri", r.RedirectURL.String())
	for _, scope := range r.Scopes {
		v.Add("scope", scope)
	}
	fullURL := *r.AuthURL
	fullURL.RawQuery = v.Encode()
	return &fullURL
}

// AccessRequest represents data needed to make
// a successful OAuth request for the access token.
type AccessRequest struct {
	AccessURL   *url.URL
	ClientID    string
	SecretID    string
	Code        string
	RedirectURL *url.URL
}

// NewAccessRequest returns pointer to AccessRequest
// ready to make access request with.
// Parses given accessURL and redirectURL
// to url.URL.
func NewAccessRequest(accessURL, clientID, secretID, code, redirectURL string) (*AccessRequest, error) {
	access, err := url.Parse(accessURL)
	if err != nil {
		return nil, err
	}
	redirect, err := url.Parse(redirectURL)
	if err != nil {
		return nil, err
	}
	return &AccessRequest{
		AccessURL:   access,
		ClientID:    clientID,
		SecretID:    secretID,
		Code:        code,
		RedirectURL: redirect,
	}, nil
}

// URL returns URL to send the access request to.
func (r *AccessRequest) URL() string {
	return r.AccessURL.String()
}

// AuthorizationHeaderValue returns string in format
//
// "Basic AuthCode" where "AuthCode" is base64 encoded
// string in format "CLIENT_ID:CLIENT_SECRET".
//
// Needed to be passed in the request header under
// "Authorization" key.
func (r *AccessRequest) AuthorizationHeaderValue() string {
	authStr := fmt.Sprintf("%s:%s", r.ClientID, r.SecretID)
	base64AuthStr := base64.StdEncoding.EncodeToString([]byte(authStr))
	return fmt.Sprintf("Basic %s", base64AuthStr)
}

// todo refresh request
