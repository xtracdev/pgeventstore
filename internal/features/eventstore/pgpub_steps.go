package eventstore

import (
	. "github.com/gucumber/gucumber"
	"github.com/stretchr/testify/assert"
	"strings"
	"os"
	"github.com/xtracdev/pgeventstore"
	log "github.com/Sirupsen/logrus"
	"github.com/xtracdev/pgconn"
	. "github.com/xtracdev/goes/sample/testagg"
)

func init() {
	var eventStore *pgeventstore.PGEventStore
	var testAgg, testAgg2 *TestAgg

	Given(`^an evironment with event publishing disabled$`, func() {
		if len(configErrors) != 0 {
			assert.Fail(T, strings.Join(configErrors, "\n"))
			return
		}

		os.Setenv(pgeventstore.EventPublishEnvVar, "0")
	})

	When(`^I store an aggregate$`, func() {
		var err error
		connectString := pgconn.BuildConnectString(DBUser,DBPassword,DBHost,DBPort,DBName)

		pgdb,err := pgconn.OpenAndConnect(connectString, 3)
		if !assert.Nil(T, err) {
			return
		}
		eventStore, err = pgeventstore.NewPGEventStore(pgdb.DB)
		if err != nil {
			log.Infof("Error connecting to oracle: %s", err.Error())
		}
		assert.NotNil(T, eventStore)
		assert.Nil(T, err)
		if assert.NotNil(T, eventStore) {
			var err error
			testAgg, err = NewTestAgg("new foo", "new bar", "new baz")
			assert.Nil(T, err)
			assert.NotNil(T, testAgg)

			err = testAgg.Store(eventStore)
			if err != nil {
				log.Infof("Error storing aggregate: %s", err.Error())
			}

			log.Infof("Stored aggregate %s", testAgg.AggregateID)
		}
	})

	Then(`^no events are written to the publish table$`, func() {
		var err error
		connectString := pgconn.BuildConnectString(DBUser,DBPassword,DBHost,DBPort,DBName)

		pgdb,err := pgconn.OpenAndConnect(connectString, 3)
		if !assert.Nil(T, err) {
			return
		}

		var count int = -1
		log.Infof("looking for publish of agg %s version %d", testAgg.AggregateID, testAgg.Version)
		err = pgdb.DB.QueryRow("select count(*) from es.t_aepb_publish where aggregate_id = $1 and version = $2", testAgg.AggregateID, testAgg.Version).Scan(&count)
		if err != nil {
			log.Infof("Error querying for published events: %s", err.Error())
		}

		assert.Nil(T, err)
		assert.Equal(T, 0, count)
	})

	Given(`^an environment with event publishing enabled$`, func() {
		if len(configErrors) != 0 {
			assert.Fail(T, strings.Join(configErrors, "\n"))
			return
		}

		os.Setenv(pgeventstore.EventPublishEnvVar, "1")
	})

	When(`^I store a new aggregate$`, func() {
		var err error
		connectString := pgconn.BuildConnectString(DBUser,DBPassword,DBHost,DBPort,DBName)

		pgdb,err := pgconn.OpenAndConnect(connectString, 3)
		if !assert.Nil(T, err) {
			return
		}
		eventStore, err = pgeventstore.NewPGEventStore(pgdb.DB)
		if err != nil {
			log.Infof("Error connecting to oracle: %s", err.Error())
		}

		if assert.NotNil(T, eventStore) {
			var err error
			testAgg2, err = NewTestAgg("new foo", "new bar", "new baz")
			assert.Nil(T, err)
			assert.NotNil(T, testAgg2)

			testAgg2.Store(eventStore)
		}
	})

	Then(`^the events are written to the publish table$`, func() {
		var err error
		connectString := pgconn.BuildConnectString(DBUser,DBPassword,DBHost,DBPort,DBName)

		pgdb,err := pgconn.OpenAndConnect(connectString, 3)
		if !assert.Nil(T, err) {
			return
		}

		var count int = -1
		err = pgdb.DB.QueryRow("select count(*) from es.t_aepb_publish where aggregate_id = $1 and version = $2", testAgg2.AggregateID, testAgg2.Version).Scan(&count)
		if assert.Nil(T, err) {
			assert.Equal(T, 1, count)
		}
	})

}

