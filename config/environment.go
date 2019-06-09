package config

import (
	"os"
)

const (
	TestScope        = "test"
	DevelopmentScope = "development"
	ProductionScope  = "production"
)

var Scope string

// GetScope returns the scope that will be used by the application:
// This can be four possible values:
//	  - 'test' for testing
//    - 'development' for local development
//    - 'production' for production scope
// The default case is development for dev purposes.
func GetScope() string {
	if Scope != "" {
		return Scope
	} else {
		envScope := os.Getenv("SCOPE")
		if envScope == "" || (envScope != ProductionScope && envScope != TestScope) {
			envScope = DevelopmentScope
		}

		Scope = envScope
		return Scope
	}

}

func InDevelopment() bool {
	return GetScope() == DevelopmentScope
}
