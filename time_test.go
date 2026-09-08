package wat

import (
	"testing"
	"time"

	"github.com/go-quicktest/qt"
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
	qt.Check(t, qt.IsFalse(tmr.Stop()))
	tmr.Reset(-1)
	<-ch
	qt.Check(t, qt.IsFalse(tmr.Stop()))
}
