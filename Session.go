package arn

import as "github.com/aerospike/aerospike-client-go"
import "github.com/aerogo/aero"

// SessionStoreAerospike is a store saving sessions in an Aerospike database.
type SessionStoreAerospike struct {
	namespace string
}

// NewAerospikeStore creates a session store using an Aerospike database.
func NewAerospikeStore() *SessionStoreAerospike {
	return &SessionStoreAerospike{
		namespace: "Session",
	}
}

// Get loads the initial session values from the database.
func (store *SessionStoreAerospike) Get(sid string) *aero.Session {
	key, _ := as.NewKey(namespace, store.namespace, sid)
	record, err := client.Get(nil, key)

	if err != nil {
		return nil
	}

	return aero.NewSession(sid, record.Bins)
}

// Set updates the session values in the database.
func (store *SessionStoreAerospike) Set(sid string, session *aero.Session) {
	sessionData := session.Data()
	key, _ := as.NewKey(namespace, store.namespace, sid)

	if len(sessionData) == 0 {
		go client.Delete(nil, key)
	} else {
		go client.PutObject(nil, key, sessionData)
	}
}
