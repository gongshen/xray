package business

import (
	"testing"
)

func TestNewTrafficQueryStatsRequestDoesNotResetCounters(t *testing.T) {
	req := NewTrafficQueryStatsRequest()

	if req.Reset_ {
		t.Fatal("traffic query request must not reset xray counters")
	}
}
