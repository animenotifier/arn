package arn

import (
	"fmt"

	"github.com/aerogo/aero"
	as "github.com/aerospike/aerospike-client-go"
)

// SessionStoreAerospike is a store saving sessions in an Aerospike database.
type SessionStoreAerospike struct {
	set string
}

// NewAerospikeStore creates a session store using an Aerospike database.
func NewAerospikeStore(set string) *SessionStoreAerospike {
	return &SessionStoreAerospike{
		set: set,
	}
}

// Get loads the initial session values from the database.
func (store *SessionStoreAerospike) Get(sid string) *aero.Session {
	key, _ := as.NewKey(namespace, store.set, sid)
	record, err := client.Get(nil, key)

	if err != nil || record == nil {
		fmt.Println(err)
		return nil
	}

	return aero.NewSession(sid, record.Bins)
}

// Set updates the session values in the database.
func (store *SessionStoreAerospike) Set(sid string, session *aero.Session) {
	sessionData := session.Data()
	key, _ := as.NewKey(namespace, store.set, sid)

	if sessionData == nil {
		client.Delete(nil, key)
	} else {
		client.Put(nil, key, sessionData)
	}
}
