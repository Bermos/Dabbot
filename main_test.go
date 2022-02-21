package main

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestGetToken(t *testing.T) {
	tests := []struct {
		tokenEnv  string
		tokenFile string
		wantToken string
		wantError bool
	}{
		{"", "", "", true},
		{"env-token", "", "env-token", false},
		{"", "./test_token", "file-token", false},
		{"env-token", "./test_token", "env-token", false},
		{"", "./not_a_file", "", true},
		{"", "./empty_test_token", "", true},
	}

	ioutil.WriteFile("./test_token", []byte("file-token"), 0644)
	ioutil.WriteFile("./empty_test_token", []byte(""), 0644)
	defer os.Remove("./test_token")
	defer os.Remove("./empty_test_token")

	for _, test := range tests {
		t.Setenv("TOKEN", test.tokenEnv)
		t.Setenv("TOKEN_FILE", test.tokenFile)

		ansToken, ansError := getToken()
		if ansToken != test.wantToken || (ansError != nil) != test.wantError {
			t.Errorf("Expected: (%q, %v), got: (%q, %v)", test.wantToken, test.wantError, ansToken, ansError != nil)
		} else {
			t.Logf("Success: (%q, %v)", ansToken, ansError != nil)
		}

		os.Unsetenv("TOKEN")
		os.Unsetenv("TOKEN_FILE")
	}
}
