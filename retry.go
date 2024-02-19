package goretry

import (
	"errors"
	"net/http"
)

type retriableRoundTripper struct {
	rt    http.RoundTripper
	times int
}

func New(options ...optFunc) *retriableRoundTripper {

	rrt := &retriableRoundTripper{}

	for _, o := range options {
		o(rrt)
	}

	if rrt.rt == nil {
		rrt.rt = http.DefaultTransport
	}

	return rrt
}

func (r retriableRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {

	var response *http.Response
	var err error

	for i := 0; i < r.times; i++ {
		response, err = r.rt.RoundTrip(request)
		if err != nil {
			continue
		}
		if response.StatusCode > 399 {
			err = errors.New("invalid status code")
			continue
		}
		break
	}

	return response, err
}
