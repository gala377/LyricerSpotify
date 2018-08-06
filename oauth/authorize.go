package oauth

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"runtime"
	"strings"
)

// AccessResponseError is returned by Access
// function if there was a non 200 response
// code value.
type AccessResponseError struct {
	// Code is the original responses code.
	Code int
	// Body is the original responses body.
	Body []byte
}

func (err *AccessResponseError) Error() string {
	return fmt.Sprintf(
		"Response with Code: %d Body: %s", err.Code, err.Body)
}

// Open opens the specified URL in the default browser of the user.
func openBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

// Authorize opens user browser and waits for authorization.
// It listens for callback on localhost + callbackURL request uri
// so callbackURL should always be on localhost but the port and
// the path can vary.
//
// Function is nonblocking and returns channel to get an authorization code from.
func Authorize(r *AuthRequest) <-chan string {
	log.Println("Starting authorization server")
	codeChan := ServCallback(
		r.RedirectURL.RequestURI(),
		":"+r.RedirectURL.Port())
	log.Println("Opening authorizarion URI")
	openBrowser(r.URL().String())
	log.Println("Returning communication channel")
	return codeChan
}

// Access makes request for the access token.
// Returns response body if the request was successful, error otherwise.
//
// Note that responses with error codes different than 200
// will be treated as an error. As they don't grant access.
// In this case response code will be returned in error as
// AccessResponseError with response code value in code field.
// Response body will be returned in the AccessResponseError
// Body field for the user to analyze.
func Access(r *AccessRequest) ([]byte, error) {
	req, err := accessRequest(r)
	if err != nil {
		return nil, err
	}
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, &AccessResponseError{
			Code: resp.StatusCode,
			Body: body,
		}
	}
	log.Printf("Got resp: %v", resp)
	return body, nil
}

func accessRequest(r *AccessRequest) (*http.Request, error) {
	postData := accessPostData(r)
	req, err := http.NewRequest("POST", r.URL(), strings.NewReader(postData.Encode()))
	if err != nil {
		return nil, err
	}
	fillInHeaders(req, r)
	// requestData := strings.Builder{}
	// req.Write(&requestData)
	// log.Printf("Request written data is: %s", requestData.String())
	return req, nil
}

func accessPostData(r *AccessRequest) url.Values {
	postData := url.Values{
		"grant_type":   {"authorization_code"},
		"code":         {r.Code},
		"redirect_uri": {r.RedirectURL.String()},
	}
	return postData
}

func fillInHeaders(request *http.Request, access *AccessRequest) {
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Authorization", access.AuthorizationHeaderValue())
	log.Printf("Filled request headers: %s", request.Header)
}
