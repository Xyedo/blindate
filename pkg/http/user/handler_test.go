package user

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/xyedo/blindate/mock"
)

func Test_CreateUser(t *testing.T) {
	app := New(mock.UserService{})

	resp := httptest.NewRecorder()
	httptest.NewRequest()
}

type testServer struct {
	*httptest.Server
}

func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewServer(h)
	return &testServer{
		ts,
	}
}
