package arn

import (
	"errors"

	"github.com/aerogo/session"
	as "github.com/aerospike/aerospike-client-go"
)

// SessionStoreAerospike is a store saving sessions in an Aerospike database.
type SessionStoreAerospike struct {
	set         string
	writePolicy *as.WritePolicy

	// Session duration in seconds (a.k.a. TTL).
	duration int
}

// NewAerospikeStore creates a session store using an Aerospike database.
func NewAerospikeStore(set string, duration int) *SessionStoreAerospike {
	writePolicy := as.NewWritePolicy(0, uint32(duration))
	writePolicy.RecordExistsAction = as.REPLACE

	return &SessionStoreAerospike{
		set:         set,
		duration:    duration,
		writePolicy: writePolicy,
	}
}

// Get loads the initial session values from the database.
func (store *SessionStoreAerospike) Get(sid string) (*session.Session, error) {
	key, _ := as.NewKey(DB.Namespace(), store.set, sid)
	record, err := DB.Client.Get(nil, key)

	if err != nil {
		return nil, err
	}

	if record == nil {
		return nil, errors.New("Record is nil (session ID: " + sid + ")")
	}

	return session.New(sid, record.Bins), nil
}

// Set updates the session values in the database.
func (store *SessionStoreAerospike) Set(sid string, session *session.Session) error {
	sessionData := session.Data()
	key, _ := as.NewKey(DB.Namespace(), store.set, sid)

	// Set with nil as data means we should delete the session.
	if sessionData == nil {
		_, err := DB.Client.Delete(nil, key)
		return err
	}

	return DB.Client.Put(store.writePolicy, key, sessionData)
}
