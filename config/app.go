// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package config

import (
	"time"
)

const NamespaceHeaderName = "G-N-ID"
const ApiKeyHeaderName = "X-API-KEY"

type APPConfig struct {
	DeployLimit     int32         `toml:"deployLimit"`
	ShutdownTimeout time.Duration `toml:"shutdownTimeout"`
	RepositoryPath  string        `toml:"repositoryPath"`
	PasswordPeriod  int           `toml:"passwordPeriod"`
}

func (a *APPConfig) OnChange() error {
	a.SetDefault()
	return nil
}

func (a *APPConfig) SetDefault() {
	if a.ShutdownTimeout == 0 {
		a.ShutdownTimeout = 10
	}
}
