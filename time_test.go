package wat

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Tests that time.AfterFunc with a negative duration still schedules the
// event.
func TestTimeAfterFunc(t *testing.T) {
	ch := make(chan struct{})
	f := func() {
		ch <- struct{}{}
	}
	tmr := time.AfterFunc(-1, f)
	<-ch
	assert.False(t, tmr.Stop())
	tmr.Reset(-1)
	<-ch
	assert.False(t, tmr.Stop())
}
