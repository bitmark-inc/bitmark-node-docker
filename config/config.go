package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	bolt "github.com/coreos/bbolt"
)

const CONFIG_BUCKET_NAME = "config"

var (
	ErrNotInitialised = fmt.Errorf("the node config is not initialised")
)

type BitmarkNodeConfig struct {
	sync.RWMutex
	initialised bool
	db          *bolt.DB
	config      map[string]string
}

var dbPath string
var nodeConfig *BitmarkNodeConfig

func New() *BitmarkNodeConfig {
	if nodeConfig == nil {
		nodeConfig = &BitmarkNodeConfig{
			config: map[string]string{
				"btcAddr": "",
				"ltcAddr": "",
			},
		}
	}
	return nodeConfig
}

func (c *BitmarkNodeConfig) Initialise(dbPath string) error {
	if !c.initialised {
		db, err := bolt.Open(filepath.Join(dbPath, "bitmark-node.db"), 0600, nil)
		if err != nil {
			return err
		}

		err = db.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists([]byte(CONFIG_BUCKET_NAME))
			if err != nil {
				return err
			}
			bkt := tx.Bucket([]byte(CONFIG_BUCKET_NAME))

			v := bkt.Get([]byte("network"))
			if v == nil {
				err := bkt.Put([]byte("network"), []byte("bitmark"))
				if err != nil {
					return err
				}
			}

			_, err = bkt.CreateBucketIfNotExists([]byte("testing"))
			if err != nil {
				return err
			}
			_, err = bkt.CreateBucketIfNotExists([]byte("bitmark"))
			return err
		})
		if err != nil {
			return err
		}
		c.db = db
		c.initialised = true
	}
	return nil
}

func (c *BitmarkNodeConfig) GetNetwork() string {
	n := os.Getenv("NETWORK")
	if n != "bitmark" && n != "testing" {
		return ""
	}
	return n
}

func (c *BitmarkNodeConfig) GetDB() (*bolt.DB, error) {
	if nil == c.db {
		return nil, errors.New("no db exist")
	}
	return c.db, nil
}

func (c *BitmarkNodeConfig) Set(newConfig map[string]string) error {
	if !c.initialised {
		return ErrNotInitialised
	}

	bucketName := c.GetNetwork()
	if bucketName == "" {
		return fmt.Errorf("empty network name")
	}

	c.Lock()
	defer c.Unlock()
	return c.db.Update(func(tx *bolt.Tx) error {
		rootBucket := tx.Bucket([]byte(CONFIG_BUCKET_NAME))
		bucket := rootBucket.Bucket([]byte(bucketName))
		for key, val := range newConfig {
			err := bucket.Put([]byte(key), []byte(val))
			c.config[key] = val
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (c *BitmarkNodeConfig) Get() (map[string]string, error) {
	if !c.initialised {
		return nil, ErrNotInitialised
	}

	bucketName := c.GetNetwork()
	if bucketName == "" {
		return nil, fmt.Errorf("empty network name")
	}

	c.RLock()
	defer c.RUnlock()
	err := c.db.View(func(tx *bolt.Tx) error {
		rootBucket := tx.Bucket([]byte(CONFIG_BUCKET_NAME))
		bucket := rootBucket.Bucket([]byte(bucketName))
		for key := range c.config {
			b := bucket.Get([]byte(key))
			if b != nil {
				c.config[key] = string(b)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return c.config, err
}
