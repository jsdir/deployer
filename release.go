package deployer

import (
	//"encoding/json"

	"github.com/boltdb/bolt"
)

type Release struct {
	Id       int
	Name     string
	Config   interface{}
	Services map[string]string
}

func NewRelease(db *bolt.DB, build *Build) (*Release, error) {
	release := new(Release)
	err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("releases"))
		if err != nil {
			return err
		}
		return b.Put([]byte(), []byte("42"))
	})

	if err != nil {
		return nil, err
	}

	return release, nil
	// Create the new release by extending the last release with the new service build.
	/*
		c := b.Cursor()
		c.Last()
		k, v := c.Next()

		// Generate metadata
		release = Release{
			Id:       k + 1,
			Name:     0,
			Services: unmarshalled,
		}

		// Update to the new build.
		release.services[build.Service] = build.Tag

		db[id] = json.serialize(release)

		// Save the new release to the database
		return &release, nil
	*/
}
