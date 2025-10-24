package test

import (
	"testing"

	"github.com/shii-park/Metasugo-Backend/internal/hub"
)

// This test verifies that SendJSON can enqueue messages up to the buffer size
// and then returns an error when the buffer is full. It avoids accessing
// unexported fields by observing SendJSON's return behavior.
func TestClientSendJSONBuffer(t *testing.T) {
	c := hub.NewClient(nil, nil, "test-user")
	// attempt more than the channel buffer (256) sends and expect at least one failure
	succeeded := 0
	var lastErr error
	for i := 0; i < 300; i++ {
		err := c.SendJSON(map[string]interface{}{"i": i})
		if err == nil {
			succeeded++
		} else {
			lastErr = err
			break
		}
	}
	if succeeded == 0 {
		t.Fatalf("expected some successful sends, got %d", succeeded)
	}
	if lastErr == nil {
		t.Fatalf("expected send buffer to become full at some point; succeeded=%d", succeeded)
	}
	t.Logf("SendJSON succeeded %d times before error: %v", succeeded, lastErr)
}
