package goretry

import (
	"errors"
	"net/http"
)

var ErrInvalidStatusCode = errors.New("invalid status code")

type retriableRoundTripper struct {
	rt         http.RoundTripper
	times      int
	statusCode int
}

func New(options ...optFunc) *retriableRoundTripper {

	rrt := &retriableRoundTripper{
		times:      1,
		statusCode: http.StatusBadRequest,
	}

	for _, o := range options {
		o(rrt)
	}

	if rrt.rt == nil {
		rrt.rt = http.DefaultTransport
	}

	return rrt
}

func (r retriableRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {

	var resp *http.Response
	var err error
	var errs []error

	for i := 0; i < r.times; i++ {
		resp, err = r.rt.RoundTrip(request)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if resp.StatusCode >= r.statusCode {
			err = ErrInvalidStatusCode
			errs = append(errs, err)
			continue
		}
		break
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return resp, nil
}
