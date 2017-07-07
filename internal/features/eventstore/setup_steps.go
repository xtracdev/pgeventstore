package eventstore

import (
	. "github.com/gucumber/gucumber"
	"os"
	"github.com/xtracdev/envinject"
)

var testEnv *envinject.InjectedEnv
var configErrors []string

func init() {
	Given(`^some tests to run$`, func() {
	})

	Then(`^the database connection configuration is read from the environment$`, func() {
	})

	GlobalContext.BeforeAll(func() {
		os.Unsetenv(envinject.ParamPrefixEnvVar)
		var err error
		testEnv, err = envinject.NewInjectedEnv()
		if err != nil {
			configErrors = append(configErrors, err.Error())
		}
	})

}
