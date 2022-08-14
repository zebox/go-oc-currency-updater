package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/rand"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"
)

func Test_getUserToken(t *testing.T) {
	port := chooseRandomUnusedPort()
	ts := prepareTestServer(port, t)
	defer ts()

	waitForServerStart(port)

	opts.BaseOpenCartUrl = fmt.Sprintf("http://127.0.0.1:%d", port)
	opts.Login = "test_login"
	opts.Password = "test_password"

	testReq, err := getUserToken()
	require.NoError(t, err)
	require.NotNil(t, testReq)

	assert.NotNil(t, testReq.cookies)
	assert.NotEmpty(t, testReq.userToken)

	// test with bad credentials
	opts.Login = "unknown"
	testReq, err = getUserToken()
	assert.Error(t, err)
	assert.Nil(t, testReq)

	// test with bad base URL
	opts.BaseOpenCartUrl = "bad_url"
	testReq, err = getUserToken()
	assert.Error(t, err)
	assert.Nil(t, testReq)

}

func Test_doUpdateCurrency(t *testing.T) {
	port := chooseRandomUnusedPort()
	ts := prepareTestServer(port, t)
	defer ts()

	waitForServerStart(port)
	opts.BaseOpenCartUrl = fmt.Sprintf("http://127.0.0.1:%d", port)
	opts.Login = "test_login"
	opts.Password = "test_password"

	testReq, err := getUserToken()
	require.NoError(t, err)
	require.NotNil(t, testReq)

	err = doUpdateCurrency(testReq)
	assert.NoError(t, err)

	// test with bad credentials
	testReq.userToken = "bad_token"
	err = doUpdateCurrency(testReq)
	assert.Error(t, err)

}
func prepareTestServer(port int, t *testing.T) func() {
	var userToken, sessionID string

	ts := &http.Server{
		Addr: fmt.Sprintf(":%d", port),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Logf("[TEST OPENCART] request %s %s %+v", r.Method, r.URL, r.Header)
			switch {
			case strings.HasPrefix(r.URL.RawQuery, "route=common/login"):
				err := r.ParseForm()
				assert.NoError(t, err)

				username := r.FormValue("username")
				password := r.FormValue("password")

				if username != "test_login" || password != "test_password" {
					w.WriteHeader(http.StatusOK)
					_, err := w.Write([]byte("bad login or password"))
					assert.NoError(t, err)
					return
				}

				userToken = randStringRunes(32)
				sessionID = randStringRunes(26)

				w.Header().Set("Content-Type", "text/html; charset=UTF-8")
				w.Header().Set("Location", fmt.Sprintf("%s%s?route=common/dashboard&user_token=%s", opts.BaseOpenCartUrl, baseResource, userToken))
				expiration := time.Now().Add(time.Minute)
				cookie := http.Cookie{Name: "OCSESSID", Value: sessionID, Expires: expiration, Path: "/"}
				http.SetCookie(w, &cookie)
				w.WriteHeader(http.StatusFound)

			case strings.HasPrefix(r.URL.RawQuery, "route=localisation/currency/refresh"):
				uToken := r.URL.Query()["user_token"]
				uCookie, err := r.Cookie("OCSESSID")
				assert.NoError(t, err)

				if uCookie.Value != sessionID || uToken[0] != userToken {
					w.WriteHeader(http.StatusOK)
					_, err := w.Write([]byte("session data"))
					assert.NoError(t, err)
					return
				}
				w.Header().Set("Location", fmt.Sprintf("%s%s?route=common/dashboard&user_token=%s", opts.BaseOpenCartUrl, baseResource, userToken))
				w.WriteHeader(http.StatusFound)
				_, err = w.Write([]byte("ok"))
				assert.NoError(t, err)
			default:
				t.Fatalf("unexpected oauth request %s %s", r.Method, r.URL)
			}
		}),
	}

	go func() { _ = ts.ListenAndServe() }()

	return func() {
		assert.NoError(t, ts.Close())
	}
}

func chooseRandomUnusedPort() (port int) {
	for i := 0; i < 10; i++ {
		port = 40000 + int(rand.Int31n(10000)) //nolint:gosec
		if ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port)); err == nil {
			_ = ln.Close()
			break
		}
	}
	return port
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func randStringRunes(n int) string {
	var charactersRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ01234567890")

	b := make([]rune, n)
	for i := range b {
		b[i] = charactersRunes[rand.Intn(len(charactersRunes))]
	}
	return string(b)
}

func waitForServerStart(port int) {
	// wait for up to 3 seconds for HTTPS server to start
	for i := 0; i < 300; i++ {
		time.Sleep(time.Millisecond * 10)
		conn, _ := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", port), time.Millisecond*10)
		if conn != nil {
			_ = conn.Close()
			break
		}
	}
}
