// Package oauth defines basic structs and functions
// needed to handle OAuth access granting flow.
package oauth

// Authorizer handles OAuth apps authorization
// process with the user of the desired
// service provider.
type Authorizer interface {
	// Authorize returns string channel to which
	// authorization code should be send and the
	// closed.
	//
	// If authotization was unsuccesful returned
	// channel should be closed immediately.
	Authorize() (<-chan string, error)
}

// Accesser handles OAuth requests for access token
// granting. As well as refreshing the acess token
// if needed.
type Accesser interface {
	// Access returns a channel to which the
	// access token is written and then closed.
	// If an error occured during granting access
	// the channel should be closed immediately.
	Access() (<-chan string, error)
	// Refresh() should return a channel to which the
	// newly granted access token is written and then closed.
	// If an error occured during refreshing access
	// the channel should be closed immediately.
	Refresh() (<-chan string, error)
}
