package deployer

import (
	"log"

	"github.com/boltdb/bolt"
)

func StartDb() *bolt.DB {
	db, err := bolt.Open("data.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	return db
}
