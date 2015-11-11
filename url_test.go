package curious

import (
	"net/url"
	"strings"
	"testing"
)

func TestURLNoScheme(t *testing.T) {
	var u url.URL
	u.Host = "blah"
	t.Log(u.String())
	if !strings.HasPrefix(u.String(), "//") {
		t.FailNow()
	}
}
