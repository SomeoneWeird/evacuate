package providers

import "github.com/sirupsen/logrus"

// ProviderContext TODO
type ProviderContext struct {
	Config map[string]string
	Logger *logrus.Entry
	Finish chan string
}

// Provider TODO
type Provider interface {
	Run(ProviderContext, string)
}

// List TODO
type List map[string]Provider

var _providers List

// RegisterProvider TODO
func RegisterProvider(name string, provider Provider) {
	if _providers == nil {
		_providers = make(List)
	}

	_providers[name] = provider
}

// GetProviders returns a list of registered providers
func GetProviders() List {
	return _providers
}
