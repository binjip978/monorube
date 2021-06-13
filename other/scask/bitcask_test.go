package main

import (
	"encoding/binary"
	"io/ioutil"
	"os"
	"testing"
)

func TestRecord(t *testing.T) {
	f, err := ioutil.TempFile("", "storage.dat")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())

	cs, b := prepareRecord("hello", "world")

	_ = binary.Write(f, binary.BigEndian, cs)
	_, _ = f.Write(b)
	_ = f.Sync()

	entry, err := readRecord(f, 0)
	if err != nil {
		t.Error(err)
	}
	if entry.key != "hello" || entry.value != "world" {
		t.Error("key or values is not correct")
	}
}
