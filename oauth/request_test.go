package oauth

import (
	"testing"
)

func TestRequestCreation(t *testing.T) {
	// todo make it 4 different tests
	// as well as test errors
	// for now it's sufficient
	r, err := NewAuthRequest(
		"http://localhost:9090",
		"someID",
		"http://localhost:9090/redirect",
		[]string{"scope1", "scope2"})
	if err != nil {
		t.Fail()
		return
	}
	if r.AuthURL.String() != "http://localhost:9090" {
		t.Fail()
		return
	}
	if r.RedirectURL.String() != "http://localhost:9090/redirect" {
		t.Fail()
		return
	}
	if r.ClientID != "someID" {
		t.Fail()
		return
	}
	slice := []string{"scope1", "scope2"}
	if len(slice) != len(r.Scopes) {
		t.Fail()
		return
	}
	for i := range slice {
		if slice[i] != r.Scopes[i] {
			t.Fail()
			return
		}
	}

}

func TestURIParsing(t *testing.T) {
	r, err := NewAuthRequest(
		"http://localhost:9090",
		"someID",
		"http://localhost:9090/redirect",
		[]string{"scope1", "scope2"})
	if err != nil {
		t.Fail()
		return
	}
	correctURI :=
		"http://localhost:9090" +
			"?client_id=someID" +
			"&redirect_uri=http%3A%2F%2Flocalhost%3A9090%2Fredirect" +
			"&response_type=code" +
			"&scope=scope1" +
			"&scope=scope2"

	requestsURI := r.URL()
	if requestsURI.String() != correctURI {
		t.Errorf(
			"URI %s is not equal to the expected: %s",
			requestsURI.String(),
			correctURI)
	}
}
