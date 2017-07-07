package main

import (
	"github.com/xtracdev/pgconn"
	"log"
	"github.com/xtracdev/pgeventstore"
	"github.com/xtracdev/envinject"
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
