package goretry

import "net/http"

type optFunc func(*retriableRoundTripper)

func WithTime(times int) optFunc {
	return func(rrt *retriableRoundTripper) {
		rrt.times = times
	}
}

func WithRoudnTriper(rt http.RoundTripper) optFunc {
	return func(rrt *retriableRoundTripper) {
		rrt.rt = rt
	}
}
