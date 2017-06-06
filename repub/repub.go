package main

import (
	"github.com/xtracdev/pgconn"
	"log"
	"os"
	"github.com/xtracdev/pgeventstore"
)

func main() {
	eventConfig, err := pgconn.NewEnvConfig()
	if err != nil {
		log.Fatal(err.Error())
	}

	pgdb, err := pgconn.OpenAndConnect(eventConfig.ConnectString(), 3)
	if err != nil {
		log.Fatal(err.Error())
	}

	os.Setenv("ES_PUBLISH_EVENTS", "1")

	eventStore, err := pgeventstore.NewPGEventStore(pgdb.DB)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = eventStore.RepublishAllEvents()
	if err != nil {
		log.Printf("Warning: %s", err.Error())
	}
}
