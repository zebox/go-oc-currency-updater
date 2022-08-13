package main

import (
	"github.com/stretchr/testify/assert"
	"os"
	"strconv"
	"testing"
)

func Test_main(t *testing.T) {
	port := chooseRandomUnusedPort()
	ts := prepareTestServer(port, t)
	defer ts()

	var tests = []struct {
		name string
		want int
		args []string
	}{
		{"without_credentials" + strconv.Itoa(port), 1, []string{"test", "-b=http://127.0.0.1:"}},
		{"bad_creds", 1, []string{"test", "-b=http://127.0.0.1:" + strconv.Itoa(port), "-u=bad_login", "-p=test_password"}},
		{"good_test", 0, []string{"test", "-b=http://127.0.0.1:" + strconv.Itoa(port), "-u=test_login", "-p=test_password"}},
	}

	for _, test := range tests {
		os.Args = test.args
		code := run()
		assert.Equal(t, test.want, code)
	}

}
