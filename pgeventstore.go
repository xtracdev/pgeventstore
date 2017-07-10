package pgeventstore

import (
	"database/sql"
	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/lib/pq"
	"github.com/xtracdev/goes"
	"time"
)

var (
	// ErrConcurrency is produced when version conflicts arise from uncoordinated write occur
	ErrConcurrency = errors.New("Concurrency Exception")
	// ErrPayloadType is produced when an unsupported payload type is used
	ErrPayloadType = errors.New("Only payloads of type []byte are allowed")
	// ErrEventInsert is produced when there's an error inserting into the events table
	ErrEventInsert = errors.New("Error inserting record into events table")
	// ErrPubInsert is oridyced when there's an error inserting into the publish table
	ErrPubInsert = errors.New("Error inserting record into pub table")
)

// PGEventStore is the type encapsulating event store functionality implemented using Postgres
// as the durable event store.
type PGEventStore struct {
	db      *sql.DB
	publish bool
}

// NewPGEventStore instantiates an event store with the provided database interface
// instance. Passing true via enableEventPublishing writes events to the publish table
// when transactions commit.
func NewPGEventStore(db *sql.DB, enableEventPublishing bool) (*PGEventStore, error) {
	log.Infof("Creating event store...")
	switch enableEventPublishing {
	case true:
		log.Info("Event store configured to write records to publish table")
	default:
		log.Info("Event store will not write records to publish table.")

	}

	return &PGEventStore{
		db:      db,
		publish: enableEventPublishing,
	}, nil
}

// GetMaxVersionForAggregate returns the greatest version number for the given
// aggregate.
func (es *PGEventStore) GetMaxVersionForAggregate(aggID string) (*int, error) {
	row, err := es.db.Query("select max(version) from t_aeev_events where aggregate_id = $1", aggID)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	var max int
	row.Scan(&max)

	err = row.Err()
	if err != nil {
		return nil, err
	}

	return &max, nil
}

// StoreEvents stores the uncommitted events assocaited with the given aggregate.
func (es *PGEventStore) StoreEvents(agg *goes.Aggregate) error {
	//Select max for the aggregate id
	max, err := es.GetMaxVersionForAggregate(agg.AggregateID)
	if err != nil {
		return err
	}

	//If the stored version is not smaller than the agg version then
	//its a concurrency exception. Note we'll have a null max if no record
	//exists
	if !(*max < agg.Version) {
		return ErrConcurrency
	}

	//Store the events
	return es.writeEvents(agg)
}

func (es *PGEventStore) writeEvents(agg *goes.Aggregate) error {

	log.Println("start transaction")
	tx, err := es.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	log.Println("create statement")
	stmt, err := tx.Prepare("insert into t_aeev_events (aggregate_id, version, typecode, payload) values ($1, $2, $3, $4)")
	if err != nil {
		return err
	}

	var pubStmt *sql.Stmt
	if es.publish {
		log.Println("create publish statement")
		var pubstmtErr error
		pubStmt, pubstmtErr = tx.Prepare("insert into t_aepb_publish (aggregate_id, version, typecode, payload) values ($1, $2, $3, $4)")
		if pubstmtErr != nil {
			return pubstmtErr
		}
	}

	for _, e := range agg.Events {
		log.Printf("process event %v\n", e)
		eventBytes, ok := e.Payload.([]byte)
		if !ok {
			stmt.Close()
			return ErrPayloadType
		}

		log.Println("execute statement")
		_, execerr := stmt.Exec(agg.AggregateID, e.Version, e.TypeCode, eventBytes)
		if execerr != nil {
			stmt.Close()
			log.Println(execerr.Error())
			return ErrEventInsert
		}

		if es.publish {
			log.Println("execute publish statement")
			_, puberr := pubStmt.Exec(agg.AggregateID, e.Version, e.TypeCode, eventBytes)
			if puberr != nil {
				log.Println(puberr.Error())
				return ErrPubInsert
			}
		}
	}

	stmt.Close()
	if pubStmt != nil {
		pubStmt.Close()
	}

	log.Println("commit transaction")
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// RetrieveEvents returns the events associated with the given aggregate.
func (es *PGEventStore) RetrieveEvents(aggID string) ([]goes.Event, error) {
	var events []goes.Event

	//Select the events, ordered by version
	rows, err := es.db.Query(`select version, typecode, payload from t_aeev_events where aggregate_id = $1 order by version`, aggID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var version int
	var typecode string
	var payload []byte

	for rows.Next() {
		rows.Scan(&version, &typecode, &payload)
		event := goes.Event{
			Source:   aggID,
			Version:  version,
			TypeCode: typecode,
			Payload:  payload,
		}

		events = append(events, event)

	}

	err = rows.Err()

	return events, err
}

// RepublishAllEvents writes all the events in the event store to the publish table.
func (es *PGEventStore) RepublishAllEvents() error {

	var tx *sql.Tx

	log.Debug("execute query")
	rows, err := es.db.Query(`select event_time, aggregate_id, version, typecode, payload from t_aeev_events order by event_time`)
	if err != nil {
		return err
	}

	defer rows.Close()

	var version int
	var typecode string
	var payload []byte
	var eventTime time.Time
	var aggregateID string

	log.Debug("create transaction")

	for rows.Next() {
		tx, err = es.db.Begin()
		if err != nil {
			return err
		}

		log.Debug("scan row")
		rows.Scan(&eventTime, &aggregateID, &version, &typecode, &payload)

		log.Debug("insert row")
		log.Infof("Publishing %s %d - %v", aggregateID, version, eventTime)
		_, err := tx.Exec(`insert into t_aepb_publish (event_time, aggregate_id, version, typecode, payload)  values($1,$2,$3,$4,$5)`,
			eventTime, aggregateID, version, typecode, payload,
		)

		if err != nil {
			pqError := err.(*pq.Error)
			log.Debug("%v", pqError.Code.Name())
			tx.Rollback()
			if pqError.Code.Name() != "unique_violation" {
				log.Debug("rollback transaction")
				return err
			}

			continue

		}

		log.Debug("commit tx")
		err = tx.Commit()
		if err != nil {
			tx.Rollback()
			return err
		}

	}

	return nil
}
