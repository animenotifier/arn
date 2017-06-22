package arn

import "testing"

func TestConnect(t *testing.T) {
	if !DB.Client.IsConnected() {
		t.Fail()
	}
}
