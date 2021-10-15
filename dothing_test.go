package main

import (
	"dothing/element"
	"testing"
)

func TestAddItem(t *testing.T) {

	dothing := new(DoThing)

	AddItem(dothing, "Testing 1..2..3", "Test")

	if dothing.Items[0].Note != "Testing 1..2..3" {
		t.Error("Expected Testing 1..2..3 item")
	}

}

func TestCreated(t *testing.T) {

	dothing := new(DoThing)

	AddItem(dothing, "Testing 1..2..3", "Test")
	if !element.Is(dothing.Items[0].Events, "Created") {
		t.Error("Expected Testing 1..2..3 to be Created")
	}

}

func TestDone(t *testing.T) {

	dothing := new(DoThing)

	AddItem(dothing, "Testing 1..2..3", "Test")

	dothing.ItemDone(0)

	if !element.Is(dothing.Done[0].Events, "Done") {
		t.Error("Expected Testing 1..2..3 to be done")
	}
}
