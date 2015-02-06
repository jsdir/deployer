package deployer

import (
	"github.com/boltdb/bolt"
)

type Context struct {
	Db        *bolt.DB
	Config    *ServerConfig
	NewBuilds *Broadcaster
}
