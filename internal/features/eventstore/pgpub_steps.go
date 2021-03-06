package eventstore

import (
	log "github.com/Sirupsen/logrus"
	. "github.com/gucumber/gucumber"
	"github.com/stretchr/testify/assert"
	. "github.com/xtracdev/goes/sample/testagg"
	"github.com/xtracdev/pgconn"
	"github.com/xtracdev/pgeventstore"
	"strings"
)

func init() {
	var eventStore *pgeventstore.PGEventStore
	var testAgg, testAgg2 *TestAgg
	var eventCount int

	Given(`^an evironment with event publishing disabled$`, func() {
		if len(configErrors) != 0 {
			assert.Fail(T, strings.Join(configErrors, "\n"))
			return
		}

	})

	When(`^I store an aggregate$`, func() {
		var err error

		pgdb, err := pgconn.OpenAndConnect(testEnv, 3)
		if !assert.Nil(T, err) {
			return
		}
		eventStore, err = pgeventstore.NewPGEventStore(pgdb.DB, false)
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

		pgdb, err := pgconn.OpenAndConnect(testEnv, 3)
		if !assert.Nil(T, err) {
			return
		}

		var count int = -1
		log.Infof("looking for publish of agg %s version %d", testAgg.AggregateID, testAgg.Version)
		err = pgdb.DB.QueryRow("select count(*) from t_aepb_publish where aggregate_id = $1 and version = $2", testAgg.AggregateID, testAgg.Version).Scan(&count)
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
	})

	When(`^I store a new aggregate$`, func() {
		var err error

		pgdb, err := pgconn.OpenAndConnect(testEnv, 3)
		if !assert.Nil(T, err) {
			return
		}
		eventStore, err = pgeventstore.NewPGEventStore(pgdb.DB, true)
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

		pgdb, err := pgconn.OpenAndConnect(testEnv, 3)
		if !assert.Nil(T, err) {
			return
		}

		var payload []byte
		var typecode string

		err = pgdb.DB.QueryRow("select  typecode, payload from t_aepb_publish where aggregate_id = $1 and version = $2",
			testAgg2.AggregateID, testAgg2.Version).Scan(&typecode, &payload)
		if assert.Nil(T, err) {
			assert.Equal(T, "TACRE", typecode)
			assert.True(T, len(payload) > 0)
		}
	})

	When(`^I republish the events$`, func() {
		var err error

		pgdb, err := pgconn.OpenAndConnect(testEnv, 3)
		if !assert.Nil(T, err) {
			return
		}

		var eventRecords int
		err = pgdb.DB.QueryRow("select  count(*) from t_aeev_events").Scan(&eventRecords)

		if !assert.Nil(T, err) {
			return
		}

		eventCount = eventRecords

		err = eventStore.RepublishAllEvents()
		if !assert.Nil(T, err) {
			return
		}

	})

	Then(`^all the events are written to the publish table$`, func() {
		var err error

		pgdb, err := pgconn.OpenAndConnect(testEnv, 3)
		if !assert.Nil(T, err) {
			return
		}

		var publishRecords int
		err = pgdb.DB.QueryRow("select  count(*) from t_aepb_publish").Scan(&publishRecords)

		if !assert.Nil(T, err) {
			return
		}

		assert.Equal(T, eventCount, publishRecords)
	})

}
