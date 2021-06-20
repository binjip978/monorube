package main

import (
	"encoding/binary"
	"io/ioutil"
	"os"
	"sort"
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

func TestBitcaskFiles(t *testing.T) {
	dir := t.TempDir()
	f, err := os.Create(dir + "/storage_1.dat")
	if err != nil {
		panic(err)
	}
	f.Close()
	f, err = os.Create(dir + "/storage_12.dat")
	if err != nil {
		panic(err)
	}
	f.Close()
	f, err = os.Create(dir + "/storage.dat")
	if err != nil {
		panic(err)
	}
	f.Close()

	files, err := bitcaskFiles(dir)
	if err != nil {
		t.Fatal(err)
	}

	sort.Strings(files)

	if files[0] != dir+"/storage_1.dat" {
		t.Errorf("%s is not correct filename", files[0])
	}
	if files[1] != dir+"/storage_12.dat" {
		t.Errorf("%s is not correct filename", files[1])
	}
	if len(files) != 2 {
		t.Error("wrong number of files")
	}
}
