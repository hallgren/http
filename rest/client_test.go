package rest

import (
	"fmt"
	"testing"

	"github.com/lunjon/http/logging"
	"github.com/stretchr/testify/require"
)

func setupClient(t *testing.T) (*Client, *URL) {
	logger := logging.NewLogger()
	client := NewClient(server.Client(), logger, logger)
	url, _ := ParseURL(server.URL, nil)
	return client, url
}

func TestBuildRequest(t *testing.T) {
	client, _ := setupClient(t)
	tests := []struct {
		method  string
		url     string
		body    string
		wantErr bool
	}{
		// Valid
		{"GET", "http://localhost", "", false},
		{"POST", "https://api.example.com:1234", "[]", false},
		{"post", "https://api.example.com:1234/path?query=something", `{"name": "lol"}`, false},
		{"DELETE", "https://api.example.com:1234/path?query=something", "", false},
		{"HEAD", "http://localhost/path", `{}`, false},
		{"Put", "http://localhost/path", `{"name": "lol"}`, false},
		{"Patch", "http://localhost/path", `{"name": "lol"}`, false},
		// Invalid
		{"", "", "", true},
		{"WHAT", "localhost/path", "", true},
	}

	var body []byte
	for _, test := range tests {
		if test.body != "" {
			body = []byte(test.body)
		}
		t.Run(test.method+" "+test.url, func(t *testing.T) {
			url, _ := ParseURL(test.url, nil)
			_, err := client.BuildRequest(test.method, url, body, nil)
			if (err != nil) != test.wantErr {
				t.Errorf("BuildRequest() error = %v, wantErr = %v", err, test.wantErr)
				return
			}
		})
	}
}

func TestClientGet(t *testing.T) {
	client, url := setupClient(t)
	req, err := client.BuildRequest("GET", url, nil, nil)
	if err != nil {
		t.Errorf("failed to build: %v", err)
		return
	}

	res := client.SendRequest(req)
	if res.Error() != nil {
		t.Errorf("failed to send: %v", err)
		return
	}
	if !res.Successful() {
		t.Errorf("failed to send: %v", err)
		return
	}
}

func TestClientPost(t *testing.T) {
	client, url := setupClient(t)
	tests := []string{
		"{}",
		`{"name": "test"}`,
		`{"array": [1,2,3,4]}`,
		`{"array": [1,2,3,4], "bool": true}`,
	}
	for _, body := range tests {
		name := fmt.Sprintf("POST %s", url)

		t.Run(name, func(t *testing.T) {
			req, err := client.BuildRequest("POST", url, []byte(body), nil)
			if err != nil {
				t.Errorf("%s failed to build: %v", name, err)
				return
			}

			res := client.SendRequest(req)
			if res.Error() != nil {
				t.Errorf("%s failed to send: %v", name, err)
				return
			}
			if !res.Successful() {
				t.Errorf("%s failed to send: %v", name, err)
				return
			}
		})
	}
}

func TestClientResult(t *testing.T) {
	client, url := setupClient(t)
	req, err := client.BuildRequest("GET", url, nil, nil)
	require.NoError(t, err)

	res := client.SendRequest(req)
	require.NoError(t, res.Error())

	require.False(t, res.HasError())
	require.True(t, res.Successful())

	require.NotNil(t, res.Request())
	require.NotZero(t, res.ElapsedMilliseconds())
	require.NotZero(t, res.Status())

	body, err := res.Body()
	require.NoError(t, err)
	require.NotEmpty(t, body)
}
