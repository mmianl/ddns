package main

import (
	"fmt"
	"testing"
)

// TestGetRegexSubstring tests that a substring can be extracted from a json string successfully
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

// TestGetRegexSubstringNoMatch tests that an error is returned when there are no matches
func TestGetRegexSubstringNoMatch(t *testing.T) {
	_, err := GetRegexSubstring(`"address":\s?"(.*)"`, `does not match`)
	if err == nil {
		t.Errorf("expected error, got %v", err)
	}

	e := "unexpected result when applying regex to does not match"
	if err.Error() != e {
		t.Errorf("wrong error, got %v, wanted %s", err, e)
	}
}

// TestGetRegexSubstringMultipleMatches tests that an error is returned when there are multiple matches
func TestGetRegexSubstringMultipleMatches(t *testing.T) {
	_, err := GetRegexSubstring(`([0-9]+) ([0-9]+)`, `12 24`)
	if err == nil {
		t.Errorf("expected error, got %v", err)
	}

	e := "unexpected result when applying regex to 12 24"
	if err.Error() != e {
		t.Errorf("wrong error, got %v, wanted %s", err, e)
	}
}

// TestGetRegexSubstringInvalidRegex tests that an invalid regex returns an error
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

// ExampleGetRegexSubstring demonstrates how to use GetRegexSubstring
func ExampleGetRegexSubstring() {
	s, _ := GetRegexSubstring(`"address":\s?"(.*)"`, `{"address": "192.168.0.10"}`)
	fmt.Println(*s)
	// Output: 192.168.0.10
}
