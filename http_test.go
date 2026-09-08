package wat

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/anacrolix/missinggo/httptoo"
	"golang.org/x/net/websocket"

	"github.com/go-quicktest/qt"
)

func TestHttpServerHandlerRequestContextDone(t *testing.T) {
	handlerRunning := make(chan struct{})
	var handlerDone <-chan struct{}
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerDone = r.Context().Done()
		select {
		case <-r.Context().Done():
			t.Fatal("request done as soon as handler started")
		default:
		}
		close(handlerRunning)
	write:
		for {
			select {
			case <-r.Context().Done():
				break write
			default:
				_, err := w.Write(make([]byte, 1024))
				if err != nil {
					break write
				}
			}
		}
	}))
	r, err := http.Get(s.URL)
	qt.Assert(t, qt.IsNil(err))
	<-handlerRunning
	qt.Check(t, qt.IsNil(r.Body.Close()))
	<-handlerDone
}

func TestWebSocketRequestContextDone(t *testing.T) {
	serverHasConn := make(chan struct{})
	var handlerDone <-chan struct{}
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerDone = r.Context().Done()
		websocket.Handler(func(ws *websocket.Conn) {
			select {
			case <-r.Context().Done():
				t.Fatal("request done as soon as server got conn")
			default:
			}
			close(serverHasConn)
			_, err := ws.Read(nil)
			qt.Check(t, qt.Equals(err, io.EOF))
			// Expect this to close when the websocket is Closed.
			// <-r.Context().Done()
		}).ServeHTTP(w, r)
	}))
	u, err := url.Parse(s.URL)
	qt.Assert(t, qt.IsNil(err))
	u = httptoo.AppendURL(u, &url.URL{
		Scheme: "ws",
	})
	ws, err := websocket.Dial(u.String(), "", s.URL)
	qt.Assert(t, qt.IsNil(err))
	<-serverHasConn
	qt.Check(t, qt.IsNil(ws.Close()))
	<-handlerDone
}
