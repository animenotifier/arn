package arn

import (
	"errors"

	as "github.com/aerospike/aerospike-client-go"
)

// Client ...
var client *as.Client

// Get ...
func Get(set string, key interface{}) (as.BinMap, error) {
	pk, keyErr := as.NewKey("arn", set, key)

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
	pk, keyErr := as.NewKey("arn", set, key)

	if keyErr != nil {
		return keyErr
	}

	return client.GetObject(nil, pk, obj)
}

// Scan ...
func Scan(set string, channel interface{}) {
	spolicy := as.NewScanPolicy()
	spolicy.ConcurrentNodes = true
	spolicy.Priority = as.HIGH
	spolicy.IncludeBinData = true

	client.ScanAllObjects(spolicy, channel, "arn", set)
}

// ForEach ...
func ForEach(set string, callback func(as.BinMap)) {
	spolicy := as.NewScanPolicy()
	spolicy.ConcurrentNodes = true
	spolicy.Priority = as.HIGH
	spolicy.IncludeBinData = true

	recs, _ := client.ScanAll(spolicy, "arn", set)

	for res := range recs.Results() {
		if res.Err != nil {
			recs.Close()
			return
		}

		callback(res.Record.Bins)
	}

	recs.Close()
}

// GetUser ...
func GetUser(id string) (*User, error) {
	user := new(User)
	err := GetObject("Users", id, user)
	return user, err
}

// GetUserByNick ...
func GetUserByNick(nick string) (*User, error) {
	rec, err := Get("NickToUser", nick)

	if err != nil {
		return nil, err
	}

	return GetUser(rec["userId"].(string))
}

// GetAnime ...
func GetAnime(id int) (*Anime, error) {
	anime := new(Anime)
	err := GetObject("Anime", id, anime)
	return anime, err
}

// GetDBHost ...
func GetDBHost() string {
	return "arn-db"
}

// init
func init() {
	as.SetAerospikeTag("json")

	var err error
	client, err = as.NewClient(GetDBHost(), 3000)

	if err != nil {
		panic(err)
	}
}
