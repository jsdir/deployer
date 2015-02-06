package deployer

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/boltdb/bolt"
)

type Release struct {
	Id       int               `json:"id"`
	Name     string            `json:"name"`
	Config   interface{}       `json:"config"`
	Services map[string]string `json:"services"`
}

func NewRelease(db *bolt.DB, build *Build) (*Release, error) {
	// Create a release
	release := new(Release)
	release.Id = 0
	release.Name = "random-name"
	release.Config = false
	release.Services = make(map[string]string)

	err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("releases"))
		if err != nil {
			return err
		}

		// Extend the params of the previous release if it exists.
		c := b.Cursor()
		k, v := c.Last()
		if k != nil {
			log.Printf("%v %v", string(k[:]), string(v[:]))
			lastRelease := new(Release)
			err := json.Unmarshal(v, &lastRelease)
			if err != nil {
				return err
			}

			release.Id = lastRelease.Id + 1
			release.Services = lastRelease.Services
		}

		// Set the new build tag.
		release.Services[build.Service] = build.Tag

		// Save the release.
		data, err := json.Marshal(&release)
		log.Printf("%v", string(data[:]))
		if err != nil {
			return err
		}

		return b.Put([]byte(strconv.Itoa(release.Id)), data)
	})

	if err != nil {
		return nil, err
	}

	return release, nil
}
