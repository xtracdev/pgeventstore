package pgeventstore

import (
	"database/sql"
	log "github.com/Sirupsen/logrus"
	"os"
)

const (
	EventPublishEnvVar = "ES_PUBLISH_EVENTS"
)

type PGEventStore struct {
	db      *sql.DB
	publish bool
}

func NewPGEventStore(db *sql.DB) (*PGEventStore, error) {
	log.Infof("Creating event store...")
	publishEvents := os.Getenv(EventPublishEnvVar)
	switch publishEvents {
	case "1":
		log.Info("Event store configured to write records to publish table")
	default:
		log.Info("Event store will not write records to publish table - export ",
			EventPublishEnvVar, "= 1 to enable writing to publish table")

	}

	return &PGEventStore{
		db:      db,
		publish: publishEvents == "1",
	}, nil
}


func (es *PGEventStore) GetMaxVersionForAggregate(aggId string) (*int, error) {
	row, err := es.db.Query("select max(version) from es.t_aeev_events where aggregate_id = $1", aggId)
	if err != nil {
		return nil, err
	}

	var max int
	row.Scan(&max)

	err = row.Err()
	if err != nil {
		return nil, err
	}

	return &max, nil
}