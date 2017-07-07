package main

import (
	"github.com/xtracdev/envinject"
	"github.com/xtracdev/pgconn"
	"github.com/xtracdev/pgeventstore"
	"log"
)

func main() {
	env, err := envinject.NewInjectedEnv()
	if err != nil {
		log.Fatal(err.Error())
	}

	pgdb, err := pgconn.OpenAndConnect(env, 3)
	if err != nil {
		log.Fatal(err.Error())
	}

	eventStore, err := pgeventstore.NewPGEventStore(pgdb.DB, true)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = eventStore.RepublishAllEvents()
	if err != nil {
		log.Printf("Warning: %s", err.Error())
	}
}
