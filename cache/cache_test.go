package cache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type Call struct {
	method string
	args   []string
}

type DBMock struct {
	data  map[string]string
	calls []Call
}

func (db *DBMock) Get(key string) (string, error) {
	db.calls = append(db.calls, Call{
		method: "Get",
		args:   []string{key},
	})

	return db.data[key], nil
}

func (db *DBMock) Keys() ([]string, error) {
	db.calls = append(db.calls, Call{
		method: "MGet",
	})

	keys := make([]string, 0, len(db.data))

	for k := range db.data {
		keys = append(keys, k)
	}

	return keys, nil
}

func (db *DBMock) GetCalls() []Call {
	return db.calls
}

func (db *DBMock) ClearCalls() {
	db.calls = []Call{}
}

func NewDBMock(mockData map[string]string) DBMock {
	return DBMock{
		data: mockData,
	}
}

func TestCache(t *testing.T) {
	testData := map[string]string{
		"key123":  "value123",
		"zhopa":   "zhopaValue",
		"hui":     "huiValue",
		"gavno":   "gavnoValue",
		"muravei": "muraveiValue",
	}

	t.Run("Should get value from db if no in cache", func(t *testing.T) {
		dbMockInstance := NewDBMock(testData)
		cache := NewCache(&dbMockInstance)

		cache.Get("zhopa")
		value, _ := cache.Get("zhopa")
		calls := dbMockInstance.GetCalls()

		require.Equal(t, value, "zhopaValue")
		require.Equal(t, calls, []Call{{method: "Get", args: []string{"zhopa"}}})
	})

	t.Run("Should go in db only one time and then take from cache", func(t *testing.T) {
		dbMockInstance := NewDBMock(testData)
		cache := NewCache(&dbMockInstance)

		for i := 0; i < 10; i++ {
			cache.Get("zhopa")
		}

		calls := dbMockInstance.GetCalls()

		require.Equal(t, len(calls), 1)
		require.Equal(t, calls, []Call{
			{
				method: "Get",
				args:   []string{"zhopa"},
			},
		})
	})

	t.Run("Should go in db in MGet if it is absents in cache", func(t *testing.T) {
		dbMockInstance := NewDBMock(testData)
		cache := NewCache(&dbMockInstance)

		cache.Get("hui")
		cache.Get("gavno")

		dbMockInstance.ClearCalls()

		result, _ := cache.MGet([]string{"hui", "gavno", "muravei"})

		calls := dbMockInstance.GetCalls()

		require.Equal(t, result, []string{"huiValue", "gavnoValue", "muraveiValue"})
		require.Equal(t, calls, []Call{{
			method: "Get",
			args:   []string{"muravei"},
		}})
	})
}
