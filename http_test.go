package wat

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/anacrolix/missinggo/httptoo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/websocket"
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
	require.NoError(t, err)
	<-handlerRunning
	assert.NoError(t, r.Body.Close())
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
			assert.Equal(t, io.EOF, err)
			// Expect this to close when the websocket is Closed.
			<-r.Context().Done()
		}).ServeHTTP(w, r)
	}))
	u, err := url.Parse(s.URL)
	require.NoError(t, err)
	u = httptoo.AppendURL(u, &url.URL{
		Scheme: "ws",
	})
	ws, err := websocket.Dial(u.String(), "", s.URL)
	require.NoError(t, err)
	<-serverHasConn
	assert.NoError(t, ws.Close())
	<-handlerDone
}
