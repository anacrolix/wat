package wat

import "testing"

func TestSelectClosedChanAndDefault(t *testing.T) {
	ch := make(chan int)
	close(ch)
	select {
	case <-ch:
	default:
		t.Fatal("selected default over closed chan")
	}
}
