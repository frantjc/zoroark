package vdf_test

import (
	"strings"
	"testing"

	"github.com/frantjc/zoroark/encoding/vdf"
)

func TestDecoder(t *testing.T) {
	var (
		actual   map[string]map[string]string
		expected = map[string]map[string]string{
			"common": {
				"name": "Counter-Strike Global Offensive - Dedicated Server",
				"type": "Tool",
			},
		}
	)

	if err := vdf.NewDecoder(strings.NewReader(`Steam>trash
	"740"
	{
		"common"
		{
			"name"      "Counter-Strike Global Offensive - Dedicated Server"
			"type"		"Tool"
		}
	}
	`)).Decode(&actual); err != nil {
		t.Error(err.Error())
		t.FailNow()
	}

	var (
		actualName   = actual["common"]["name"]
		expectedName = expected["common"]["name"]
	)
	if actualName != expectedName {
		t.Error("actual was", actualName, "but expected", expectedName)
		t.FailNow()
	}
}
