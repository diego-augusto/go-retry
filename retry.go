package goretry

import (
	"errors"
	"net/http"
)

type retriableRoudnTriper struct {
	rt    http.RoundTripper
	times int
}

func New(times int, rt http.RoundTripper) http.RoundTripper {
	if rt == nil {
		rt = http.DefaultTransport
	}
	return retriableRoudnTriper{
		rt:    rt,
		times: times,
	}
}

func (r retriableRoudnTriper) RoundTrip(request *http.Request) (*http.Response, error) {

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
