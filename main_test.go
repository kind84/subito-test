package main

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

func TestTraverse(t *testing.T) {
	var rooms Rooms
	start := 4
	objects := []string{"Knife", "Potted Plant", "Pillow"}

	raw, err := ioutil.ReadFile("./map2.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(raw, &rooms)
	if err != nil {
		panic(err)
	}

	_, steps := traverse(&rooms, start, objects)
	if len(*steps) != 9 {
		t.Error("Expected 9 steps, got", len(*steps))
	}
}

func TestGetNext(t *testing.T) {
	edges := []int{1, 2, 3, 4}
	visited := map[int]bool{
		1: true,
		2: true,
		3: true,
		4: false,
	}
	rm := map[int]*Room{
		1: &Room{1, "Pippo", 0, 0, 0, 0, []Object{}},
		2: &Room{2, "Pluto", 0, 0, 0, 0, []Object{}},
		3: &Room{3, "Paperino", 0, 0, 0, 0, []Object{}},
		4: &Room{4, "Topolino", 0, 0, 0, 0, []Object{}},
	}
	r := getNext(edges, visited, rm)
	if r == nil {
		t.Error("Expected room n.4, got nil")
	} else if r.ID != 4 {
		t.Error("Expected room n.4, got", r.ID)
	}
}
