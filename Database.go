package arn

import (
	"errors"
	"reflect"

	as "github.com/aerospike/aerospike-client-go"
)

func init() {
	// This will make Aerospike use json tags for the field names in  the database.
	as.SetAerospikeTag("json")
}

// DB is the main database client.
var DB = NewDatabase(
	"arn-db",
	3000,
	"arn",
	[]interface{}{
		new(Anime),
		new(AnimeList),
		new(Post),
		new(Settings),
		new(Thread),
		new(User),
	},
)

// Database represents the Aerospike database.
type Database struct {
	namespace string
	types     map[string]reflect.Type
	client    *as.Client
}

// NewDatabase creates a new database client.
func NewDatabase(host string, port int, namespace string, tables []interface{}) *Database {
	// Convert example objects to their respective types
	tableTypes := make(map[string]reflect.Type)
	for _, example := range tables {
		typeInfo := reflect.TypeOf(example).Elem()
		tableTypes[typeInfo.Name()] = typeInfo
	}

	// Create client
	client, err := as.NewClient(host, port)

	if err != nil {
		panic(err)
	}

	// Make Set() calls delete old fields instead of only updating new ones
	client.DefaultWritePolicy.RecordExistsAction = as.REPLACE

	// Make scans faster
	client.DefaultScanPolicy.Priority = as.HIGH
	client.DefaultScanPolicy.ConcurrentNodes = true
	client.DefaultScanPolicy.IncludeBinData = true

	return &Database{
		namespace: namespace,
		types:     tableTypes,
		client:    client,
	}
}

// Get retrieves an object from the table.
func (db *Database) Get(table string, id string) (interface{}, error) {
	pk, keyErr := as.NewKey(db.namespace, table, id)

	if keyErr != nil {
		return nil, keyErr
	}

	t, exists := db.types[table]

	if !exists {
		return nil, errors.New("Data type has not been defined for table " + table)
	}

	obj := reflect.New(t).Interface()
	err := db.client.GetObject(nil, pk, obj)

	return obj, err
}

// Set sets an object's data for the given ID and erases old fields.
func (db *Database) Set(table string, id string, obj interface{}) error {
	pk, keyErr := as.NewKey(db.namespace, table, id)

	if keyErr != nil {
		return keyErr
	}

	// TODO: Implement write policy with as.REPLACE
	return db.client.PutObject(nil, pk, obj)
}

// Delete deletes an object from the database and returns if it existed.
func (db *Database) Delete(table string, id string) (existed bool, err error) {
	pk, keyErr := as.NewKey(db.namespace, table, id)

	if keyErr != nil {
		return false, keyErr
	}

	return db.client.Delete(nil, pk)
}

// Scan writes all objects from a given table to the channel.
func (db *Database) Scan(table string, channel interface{}) error {
	_, err := db.client.ScanAllObjects(nil, channel, db.namespace, table)
	return err
}

// All returns a stream of all objects in the given table.
func (db *Database) All(table string) (interface{}, error) {
	channel := reflect.MakeChan(db.types[table], 0)
	err := db.Scan(table, channel)
	return channel, err
}

// GetObject retrieves data from the table and stores it in the provided object.
func (db *Database) GetObject(table string, id string, obj interface{}) error {
	pk, keyErr := as.NewKey(db.namespace, table, id)

	if keyErr != nil {
		return keyErr
	}

	return db.client.GetObject(nil, pk, obj)
}

// GetMap retrieves the data as a map[string]interface{}.
func (db *Database) GetMap(table string, id string) (as.BinMap, error) {
	pk, keyErr := as.NewKey(db.namespace, table, id)

	if keyErr != nil {
		return nil, keyErr
	}

	rec, err := db.client.Get(nil, pk)

	if err != nil {
		return nil, err
	}

	if rec == nil {
		return nil, errors.New("Record not found")
	}

	return rec.Bins, nil
}

// DeleteTable deletes a table.
func (db *Database) DeleteTable(table string) error {
	return db.client.Truncate(nil, db.namespace, table, nil)
}

// // ForEach ...
// func ForEach(set string, callback func(as.BinMap)) {
// 	recs, _ := client.ScanAll(scanPolicy, namespace, set)

// 	for res := range recs.Results() {
// 		if res.Err != nil {
// 			recs.Close()
// 			return
// 		}

// 		callback(res.Record.Bins)
// 	}

// 	recs.Close()
// }
