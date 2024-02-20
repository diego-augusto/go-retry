package goretry_test

import (
	"errors"
	"net/http"
	"testing"

	goretry "github.com/diego-augusto/go-retry"
	"github.com/stretchr/testify/assert"
)

type mockDefaultTransport struct {
	called int
	MockFn func(*http.Request) (*http.Response, error)
}

func (m *mockDefaultTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	m.called++
	return m.MockFn(request)
}

func Test_RoundTrip(t *testing.T) {
	testCases := []struct {
		name        string
		err         error
		resp        *http.Response
		wantErr     error
		wantRespose *http.Response
	}{
		{
			name:    "error",
			err:     errors.New("http error"),
			wantErr: errors.New("http error"),
		},
		{
			name:    "status bad request",
			resp:    &http.Response{StatusCode: http.StatusBadRequest},
			wantErr: errors.New("invalid status code"),
		},
		{
			name:        "status ok",
			resp:        &http.Response{StatusCode: http.StatusOK},
			wantRespose: &http.Response{StatusCode: http.StatusOK},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			mock := mockDefaultTransport{
				MockFn: func(r *http.Request) (*http.Response, error) {
					return tc.resp, tc.err
				},
			}

			client := http.Client{
				Transport: goretry.New(goretry.WithRoudnTriper(&mock)),
			}

			req, _ := http.NewRequest(http.MethodGet, "http://www.github.com", nil)
			assert.NotNil(t, req)

			resp, err := client.Do(req)

			if tc.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorAs(t, err, &tc.wantErr)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.wantRespose.StatusCode, resp.StatusCode)
		})
	}
}

func Test_Default(t *testing.T) {

	client := http.Client{
		Transport: goretry.New(),
	}

	req, _ := http.NewRequest(http.MethodGet, "http://www.github.com", nil)
	assert.NotNil(t, req)

	resp, err := client.Do(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

}

func Test_Times(t *testing.T) {

	wantTime := 50
	wantError := errors.New("random network error")

	mock := mockDefaultTransport{
		MockFn: func(r *http.Request) (*http.Response, error) {
			return nil, wantError
		},
	}

	client := http.Client{
		Transport: goretry.New(goretry.WithRoudnTriper(&mock), goretry.WithTime(wantTime)),
	}

	req, _ := http.NewRequest(http.MethodGet, "http://www.github.com", nil)
	assert.NotNil(t, req)

	resp, err := client.Do(req)

	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Equal(t, wantTime, mock.called)
}

func Test_WithStatusCode(t *testing.T) {

	testCases := []struct {
		name    string
		url     string
		wantErr error
	}{
		{
			name:    "invalid status code",
			url:     "https://mock.httpstatus.io/500",
			wantErr: goretry.ErrInvalidStatusCode,
		},
		{
			name: "valid status code",
			url:  "https://mock.httpstatus.io/200",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			client := http.Client{
				Transport: goretry.New(goretry.WithStatusCode(http.StatusInternalServerError)),
			}

			req, err := http.NewRequest(http.MethodGet, tc.url, nil)
			assert.NotNil(t, req)
			assert.Nil(t, err)

			resp, err := client.Do(req)

			if tc.wantErr != nil {
				assert.Error(t, err)
				assert.Nil(t, resp)
				assert.ErrorIs(t, err, goretry.ErrInvalidStatusCode)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, resp)
		})
	}
}
