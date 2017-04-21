package eventstore

import (
	. "github.com/gucumber/gucumber"
	"github.com/stretchr/testify/assert"
	"strings"
	log "github.com/Sirupsen/logrus"
	"github.com/xtracdev/pgconn"
	"github.com/xtracdev/pgeventstore"
	. "github.com/xtracdev/goes/sample/testagg"
)

func init() {
	var testAgg *TestAgg
	//var anotherAgg *TestAgg

	Given(`^a new aggregate instance$`, func() {
		if len(configErrors) != 0 {
			assert.Fail(T, strings.Join(configErrors, "\n"))
			return
		}

		log.Info("create event store")
		connectString := pgconn.BuildConnectString(DBUser,DBPassword,DBHost,DBPort,DBName)

		pgdb,err := pgconn.OpenAndConnect(connectString, 3)
		if !assert.Nil(T, err) {
			return
		}

		eventStore,_ := pgeventstore.NewPGEventStore(pgdb.DB)
		if assert.NotNil(T, eventStore) {
			var err error
			testAgg, err = NewTestAgg("new foo", "new bar", "new baz")
			assert.Nil(T, err)
			assert.NotNil(T, testAgg)
		}

	})

	When(`^we check the max version in the event store$`, func() {
		T.Skip() // pending
	})

	Then(`^the max version is (\d+)$`, func(i1 int) {
		T.Skip() // pending
	})

}


