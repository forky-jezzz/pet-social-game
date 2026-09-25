package security

import (
	"context"
	"fmt"

	"pet-social-game/server/internal/config"
)

// AlipayDetector is a skeleton for alipay.security.risk.content.detect.
// Wire real OpenAPI signing when app credentials are available.
type AlipayDetector struct {
	cfg config.AlipayConfig
}

func NewAlipayDetector(cfg config.AlipayConfig) *AlipayDetector {
	return &AlipayDetector{cfg: cfg}
}

func (a *AlipayDetector) DetectText(ctx context.Context, content string) (bool, string, error) {
	if a.cfg.AppID == "" || a.cfg.PrivateKey == "" {
		return false, "alipay_not_configured", fmt.Errorf("alipay content detect not configured")
	}
	// TODO: call alipay.security.risk.content.detect with RSA2 signed request.
	// Placeholder keeps interface stable; fall back behavior should be set by caller.
	_ = ctx
	_ = content
	return true, "", nil
}

func NewDetector(cfg config.SecurityConfig) Detector {
	switch cfg.Provider {
	case "alipay":
		return NewAlipayDetector(cfg.Alipay)
	default:
		return NewMockDetector()
	}
}
