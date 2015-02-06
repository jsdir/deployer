package resources

import (
	"encoding/binary"
	"encoding/json"

	"github.com/boltdb/bolt"
	"github.com/jsdir/deployer/pkg/names"
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
	release.Name = names.NewRandomName("-")
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
		if err != nil {
			return err
		}

		return b.Put(convertIdToBytes(release.Id), data)
	})

	if err != nil {
		return nil, err
	}

	return release, nil
}

func GetRelease(db *bolt.DB, id int) (*Release, error) {
	var release Release

	err = db.Update(func(tx *blot.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("releases"))
		if err != nil {
			return err
		}

		data = b.Get(convertIdToBytes(idInt))
		return json.Unmarshal(data, &release)
	})

	if err != nil {
		return nil, err
	}

	return release, nil
}

func convertIdToBytes(id integer) []bytes {
	result := make([]byte, 4)
	binary.LittleEndian.PutUint16(idBytes, uint16(id))
	return result
}
