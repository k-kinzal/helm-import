package main

import (
	"github.com/spf13/pflag"
	"k8s.io/helm/pkg/helm/environment"
)

var (
	Env = &struct {
		BaseUrl string
	}{
		BaseUrl: "http://127.0.0.1:8879/charts",
	}
	HelmEnv = initHelmEnv()
)

func initHelmEnv() environment.EnvSettings {
	env := environment.EnvSettings{}
	fs := &pflag.FlagSet{}
	env.AddFlags(fs)
	env.Init(fs)

	return env
}
