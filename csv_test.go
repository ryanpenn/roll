package main

import "testing"

type TestItem struct {
	TestId       int      `json:"TestId"`
	TestType     int      `json:"TestType"`
	TestName     string   `json:"TestName"`
	TestValues   []int    `json:"TestValues"`
	TestFeatures []string `json:"TestFeatures"`
}

func TestCsvUtil_LoadFile(t *testing.T) {
	file := "./data/_Test.csv"
	if items, err := LoadFile[TestItem](file, "json"); err != nil {
		t.Error(err)
	} else {
		t.Log("success", len(items))
	}

	if items, err := LoadFile[*TestItem](file, "json"); err != nil {
		t.Error(err)
	} else {
		t.Log("success", len(items))
	}
}
