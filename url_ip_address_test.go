package main

import (
	"fmt"
	"testing"
)

func TestGetRegexSubstring(t *testing.T) {
	want := "192.168.0.10"
	got, err := GetRegexSubstring(`"address":\s?"(.*)"`, `{"address": "192.168.0.10"}`)
	if err != nil {
		t.Errorf("got %s, wanted %s", err, want)
	}

	if *got != want {
		t.Errorf("got %s, wanted %s", *got, want)
	}
}

func TestGetRegexSubstringNoMatch(t *testing.T) {
	_, err := GetRegexSubstring(`"address":\s?"(.*)"`, `does not match`)
	if err == nil {
		t.Errorf("expected error, got %v", err)
	}

	e := "unexpcted result when applying regex to does not match"
	if err.Error() != e {
		t.Errorf("wrong error, got %v, wanted %s", err, e)
	}
}

func TestGetRegexSubstringMultipleMatches(t *testing.T) {
	_, err := GetRegexSubstring(`([0-9]+) ([0-9]+)`, `12 24`)
	if err == nil {
		t.Errorf("expected error, got %v", err)
	}

	e := "unexpcted result when applying regex to 12 24"
	if err.Error() != e {
		t.Errorf("wrong error, got %v, wanted %s", err, e)
	}
}

func TestGetRegexSubstringInvalidRegex(t *testing.T) {
	_, err := GetRegexSubstring(`(.*`, `does not match`)
	if err == nil {
		t.Errorf("expected error, got %v", err)
	}

	e := "error parsing regexp: missing closing ): `(.*`"
	if err.Error() != e {
		t.Errorf("wrong error, got %v, wanted %s", err, e)
	}
}

func ExampleGetRegexSubstring() {
	s, _ := GetRegexSubstring(`"address":\s?"(.*)"`, `{"address": "192.168.0.10"}`)
	fmt.Println(*s)
	// Output: 192.168.0.10
}
