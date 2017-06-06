package arn

import (
	"errors"

	as "github.com/aerospike/aerospike-client-go"
)

var namespace = "arn"
var client *as.Client
var scanPolicy *as.ScanPolicy

// Get ...
func Get(set string, key interface{}) (as.BinMap, error) {
	pk, keyErr := as.NewKey(namespace, set, key)

	if keyErr != nil {
		return nil, keyErr
	}

	rec, err := client.Get(nil, pk)

	if err != nil {
		return nil, err
	}

	if rec == nil {
		return nil, errors.New("Record not found")
	}

	return rec.Bins, nil
}

// GetObject ...
func GetObject(set string, key interface{}, obj interface{}) error {
	pk, keyErr := as.NewKey(namespace, set, key)

	if keyErr != nil {
		return keyErr
	}

	return client.GetObject(nil, pk, obj)
}

// SetObject ...
func SetObject(set string, key interface{}, obj interface{}) error {
	pk, keyErr := as.NewKey(namespace, set, key)

	if keyErr != nil {
		return keyErr
	}

	return client.PutObject(nil, pk, obj)
}

// Delete ...
func Delete(set string, key interface{}) (bool, error) {
	pk, keyErr := as.NewKey(namespace, set, key)

	if keyErr != nil {
		return false, keyErr
	}

	return client.Delete(nil, pk)
}

// Truncate ...
func Truncate(set string) error {
	return client.Truncate(nil, namespace, set, nil)
}

// Scan ...
func Scan(set string, channel interface{}) error {
	_, err := client.ScanAllObjects(scanPolicy, channel, namespace, set)
	return err
}

// ForEach ...
func ForEach(set string, callback func(as.BinMap)) {
	recs, _ := client.ScanAll(scanPolicy, namespace, set)

	for res := range recs.Results() {
		if res.Err != nil {
			recs.Close()
			return
		}

		callback(res.Record.Bins)
	}

	recs.Close()
}

// GetDBHost ...
func GetDBHost() string {
	return "arn-db"
}

// init
func init() {
	as.SetAerospikeTag("json")

	scanPolicy = as.NewScanPolicy()
	scanPolicy.ConcurrentNodes = true
	scanPolicy.Priority = as.HIGH
	scanPolicy.IncludeBinData = true

	var err error
	client, err = as.NewClient(GetDBHost(), 3000)

	if err != nil {
		panic(err)
	}
}
