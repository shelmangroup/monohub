package hooks

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func setupTestCase(t *testing.T) func(t *testing.T) {
	t.Log("setup test case")
	oldStdin := os.Stdin

	return func(t *testing.T) {
		t.Log("teardown test case")
		os.Stdin = oldStdin
	}
}

func TestRunPreRecevie(t *testing.T) {
	cases := []struct {
		name     string
		content  []byte
		expected string
	}{
		{"should pass", []byte("0000000000000000000000000000000000000000 4ea06c7022cc6ad23d2a62361a1935f49f123456 refs/heads/master"), ""},
		{"prevent force push", []byte("4ea06c7022cc6ad23d2a62361a1935f49f5168b3 4ea06c7022cc6ad23d2a62361a1935f49f123456 refs/heads/master"), "exit status 128"},
		{"allow force push not master", []byte("4ea06c7022cc6ad23d2a62361a1935f49f5168b3 4ea06c7022cc6ad23d2a62361a1935f49f123456 refs/heads/foobar"), ""},
	}
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tmpfile, err := ioutil.TempFile("", "test")
			if err != nil {
				log.Fatal(err)
			}
			defer os.Remove(tmpfile.Name()) // clean up

			if _, err := tmpfile.Write(tc.content); err != nil {
				log.Fatal(err)
			}

			if _, err := tmpfile.Seek(0, 0); err != nil {
				log.Fatal(err)
			}

			os.Stdin = tmpfile
			if err := RunHookPreReceive(); err != nil {
				if err.Error() != tc.expected {
					t.Errorf("RunHookPreReceive failed: %v", err)
				}
			}

			if err := tmpfile.Close(); err != nil {
				log.Fatal(err)
			}
		})
	}
}
