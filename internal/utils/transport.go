package utils

import "net/http"

type UserAgentTransport struct {
	UserAgent string
	Transport http.RoundTripper
}

func (u *UserAgentTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", u.UserAgent)
	return u.Transport.RoundTrip(req)
}
