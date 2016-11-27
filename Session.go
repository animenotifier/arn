package arn

import as "github.com/aerospike/aerospike-client-go"

// SessionDatabase ...
type SessionDatabase struct {
}

// Load loads the initial session values from the database.
func (db *SessionDatabase) Load(sid string) map[string]interface{} {
	key, _ := as.NewKey(namespace, "Sessions", sid)
	record, err := client.Get(nil, key)

	if err != nil {
		return nil
	}

	return record.Bins
}

// Update updates the session values in the database.
func (db *SessionDatabase) Update(sid string, newValues map[string]interface{}) {
	key, _ := as.NewKey(namespace, "Sessions", sid)

	if len(newValues) == 0 {
		go client.Delete(nil, key)
	} else {
		go client.PutObject(nil, key, newValues)
	}
}
