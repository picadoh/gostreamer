package main

import (
	"testing"
)

func TestCounterIncrement(t *testing.T) {

	counter := NewCounter()

	xCount1 := counter.Increment("x")
	xCount2 := counter.Increment("x")
	yCount := counter.Increment("y")

	if (xCount1 != 1) {
		t.Error("Expected 1, got ", xCount1)
	}

	if (xCount2 != 2) {
		t.Error("Expected 2, got ", xCount2)
	}

	if (yCount != 1) {
		t.Error("Expected 1, got ", yCount)
	}
}

func TestGetValue(t *testing.T) {

	counter := NewCounter()

	xCount := counter.GetValue("x")

	if (xCount != 0) {
		t.Error("Expected 0, got ", xCount)
	}

	counter.Increment("x")
	xCount = counter.GetValue("x")

	if (xCount != 1) {
		t.Error("Expected 1, got ", xCount)
	}
}
