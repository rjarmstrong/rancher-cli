package config

import (
	"time"

	"github.com/rancher/go-rancher/client"
)

type Validator interface {
	Validate(lc *client.LaunchConfig, opts UpgradeOpts) error
}

type NoopValidator struct{}

func (val *NoopValidator) Validate(lc *client.LaunchConfig, opts UpgradeOpts) error {
	return nil
}

type UpgradeOpts struct {
	Envs        []string
	Wait        bool
	ServiceLike string
	Service     string
	CodeTag     string
	RuntimeTag  string
	Interval    time.Duration
}
