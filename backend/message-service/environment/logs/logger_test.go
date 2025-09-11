package logs

import (
	"fmt"
	"testing"
	"time"
)

func TestLogger(t *testing.T) {
	logger := NewLogger("test")
	logger.Info("working")

	path := "/api/v1/new_login"
	logger.InfoMetrics(NewInfoMetrics(&path, 200, 5*time.Second))
	logger.InfoMetrics(NewInfoMetrics(nil, 200, 5*time.Second))

	logger.Error(fmt.Errorf("mock error"))

	logger.Warn("mock warn")
}
