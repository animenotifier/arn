package arn

import (
	"fmt"

	"github.com/aerogo/session"
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
func (store *SessionStoreAerospike) Get(sid string) *session.Session {
	key, _ := as.NewKey(DB.Namespace(), store.set, sid)
	record, err := DB.Client.Get(nil, key)

	if err != nil || record == nil {
		fmt.Println(err)
		return nil
	}

	return session.New(sid, record.Bins)
}

// Set updates the session values in the database.
func (store *SessionStoreAerospike) Set(sid string, session *session.Session) {
	sessionData := session.Data()
	key, _ := as.NewKey(DB.Namespace(), store.set, sid)

	if sessionData == nil {
		DB.Client.Delete(nil, key)
	} else {
		DB.Client.Put(nil, key, sessionData)
	}
}
