package main

import (
	"strings"
	"testing"
)

var m = []map[string]string{{}}

func TestGetSwdjFields(t *testing.T) {
	tb := "gt3_dj_nsrxx"
	got := getFields(tb)
	for _, v := range got {
		m[0][v.name] = v.name
	}
}
func TestDbSwdjInsert(t *testing.T) {
	tb := "gt3_dj_nsrxx"
	got := getFields(tb)
	for _, v := range got {
		m[0][v.name] = v.name
		if strings.Compare(v.name, "SWJG_DM") == 0 {
			m[0][v.name] = "烟台市税务局"
		}
	}
	dbInsert(m, "gt3_dj_nsrxx")
}
func TestDbJrkInsert(t *testing.T) {
	tb := "gt3_jrk"
	got := getFields(tb)
	for _, v := range got {
		m[0][v.name] = v.name
	}
	dbInsert(m, tb)
}
func TestDbMdtInsert(t *testing.T) {
	tb := "gt3_mdt"
	got := getFields(tb)
	for _, v := range got {
		m[0][v.name] = v.name
	}
	dbInsert(m, tb)
}
func TestDbSbqcInsert(t *testing.T) {
	tb := "gt3_ybnsr_sbqc"
	got := getFields(tb)
	for _, v := range got {
		m[0][v.name] = v.name
	}
	dbInsert(m, tb)
}
