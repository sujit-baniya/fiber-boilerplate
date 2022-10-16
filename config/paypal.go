package config

import (
	"github.com/plutov/paypal/v3"
)

type PayPalConfig struct {
	*paypal.Client
	ClientID string `mapstructure:"PAYPAL_CLIENT_ID" yaml:"client_id" env:"PAYPAL_CLIENT_ID"`
	Secret   string `mapstructure:"PAYPAL_SECRET" yaml:"secret" env:"PAYPAL_SECRET"`
	Account  string `mapstructure:"PAYPAL_ACCOUNT" yaml:"account" env:"PAYPAL_ACCOUNT"`
	Mode     string `mapstructure:"PAYPAL_MODE" yaml:"mode" env:"PAYPAL_MODE" env-default:"sandbox"`
}

func (p *PayPalConfig) Connect(env string) {
	var err error
	if p.Client != nil {
		return
	}
	if env == "prod" {
		p.Client, err = paypal.NewClient(p.ClientID, p.Secret, paypal.APIBaseLive)
	} else {
		p.Client, err = paypal.NewClient(p.ClientID, p.Secret, paypal.APIBaseSandBox)
	}
	if err != nil {
		panic(err)
	}
}
