package arn

import "testing"

func TestConnect(t *testing.T) {
	if !DB.client.IsConnected() {
		t.Fail()
	}
}
