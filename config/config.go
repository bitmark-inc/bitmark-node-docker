package config

import (
	"fmt"
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
		db, err := bolt.Open(dbPath, 0600, nil)
		if err != nil {
			return err
		}

		err = db.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists([]byte(CONFIG_BUCKET_NAME))
			return err
		})
		if err != nil {
			return err
		}
		c.db = db
		c.initialised = true
	}

	_, err := nodeConfig.Get()
	if err != nil {
		return err
	}
	return nil
}

func (c *BitmarkNodeConfig) Set(newConfig map[string]string) error {
	if !c.initialised {
		return ErrNotInitialised
	}
	c.Lock()
	defer c.Unlock()
	return c.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(CONFIG_BUCKET_NAME))
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
	c.RLock()
	defer c.RUnlock()
	err := c.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(CONFIG_BUCKET_NAME))
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
