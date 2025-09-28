package cache

import "sync"

type DB interface {
	Get(key string) (string, error)
}

type CacheI interface {
	Get(key string) (string, error)
	KeysInCache() ([]string, error)
	MGet(keys []string) ([]string, error)
}

type Cache struct {
	mu   sync.Mutex
	data map[string]string
	db   DB
}

func (c *Cache) Get(key string) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	value, exist := c.data[key]

	if !exist {
		dbValue, error := c.db.Get(key)

		if error != nil {
			return "", error
		}

		c.data[key] = dbValue

		return dbValue, nil
	}

	return value, nil
}

func (c *Cache) MGet(keys []string) ([]string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	values := make([]string, 0, len(keys))

	for _, key := range keys {
		valueFromCache, exist := c.data[key]

		if exist {
			values = append(values, valueFromCache)
			continue
		}

		valueFromDb, err := c.db.Get(key)

		if err != nil {
			return nil, err
		}

		values = append(values, valueFromDb)
	}

	return values, nil
}

func (c *Cache) KeysInCache() ([]string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	keys := make([]string, 0, len(c.data))

	for k := range c.data {
		keys = append(keys, k)
	}

	return keys, nil
}

func NewCache(db DB) Cache {
	return Cache{
		db:   db,
		data: make(map[string]string),
	}
}
