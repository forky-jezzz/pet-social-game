package security

import (
	"context"
	"strings"
)

type Detector interface {
	DetectText(ctx context.Context, content string) (pass bool, reason string, err error)
}

type MockDetector struct{}

func NewMockDetector() *MockDetector { return &MockDetector{} }

func (m *MockDetector) DetectText(_ context.Context, content string) (bool, string, error) {
	blocked := []string{"违禁", "赌博", "色情"}
	for _, w := range blocked {
		if strings.Contains(content, w) {
			return false, "content_risk", nil
		}
	}
	return true, "", nil
}
