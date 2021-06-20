package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

type entry struct {
	crc       uint32
	timestamp int32
	keySize   uint32
	valueSize uint64
	key       string
	value     string
}

const headerSize = 4 + 4 + 4 + 8

type cask struct {
	dataFile   *os.File
	index      map[string]keyRecord
	lastOffset int64
	sync.RWMutex
	nameFileMap map[string]*os.File
	storageDir  string
}

type keyRecord struct {
	fileID    string
	valueSize uint64
	valuePos  int64
	ts        int32
}

var storageFileRegexp = regexp.MustCompile(`storage_\d+\.dat`)

// newSCask should open a directory and if its empty create storage file
// and setup in memory index, if directory is not empty it should rebuild
// in memory index and make the latest file as open, open file is file that
// accepts all recent changes
func newSCask(storageDir string) (*cask, error) {
	var dataFiles []string
	err := filepath.WalkDir(storageDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			panic(err)
		}

		if d.Type().IsRegular() && storageFileRegexp.MatchString(d.Name()) {
			dataFiles = append(dataFiles, d.Name())
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	log.Println("files on disk: ", dataFiles)

	// create first dataFile in empty dir and cask object
	if len(dataFiles) == 0 {
		f, err := os.OpenFile(storageDir+"/storage_1.dat", os.O_CREATE|os.O_RDWR|os.O_EXCL, 0755)
		if err != nil {
			return nil, err
		}

		nameFileMap := make(map[string]*os.File)
		nameFileMap[f.Name()] = f

		return &cask{
			dataFile:    f,
			index:       make(map[string]keyRecord),
			lastOffset:  0,
			nameFileMap: nameFileMap,
			storageDir:  storageDir,
		}, nil
	}

	index := make(map[string]keyRecord)
	nameFileMap := make(map[string]*os.File)

	for _, file := range dataFiles {
		// TODO: we need this file handle to populate nameFileMap
		err = readDataFile(storageDir+"/"+file, index)
		if err != nil {
			return nil, err
		}
	}

	var lastesID int

	for _, file := range dataFiles {
		sp := strings.Split(file, "_")
		sp = strings.Split(sp[len(sp)-1], ".")
		i, err := strconv.Atoi(sp[0])
		if err != nil {
			panic(err)
		}
		if i > lastesID {
			lastesID = i
		}
		f, err := os.Open(storageDir + "/" + file)
		nameFileMap[f.Name()] = f
	}

	of := nameFileMap[storageDir+"/"+fmt.Sprintf("storage_%d.dat", lastesID)]
	lastestName := of.Name()
	_ = of.Close()

	of, err = os.OpenFile(lastestName, os.O_APPEND|os.O_RDWR|os.O_EXCL, 0755)
	if err != nil {
		return nil, err
	}

	nameFileMap[of.Name()] = of

	fileInfo, err := of.Stat()
	if err != nil {
		return nil, err
	}

	return &cask{
		dataFile:    of,
		index:       index,
		lastOffset:  fileInfo.Size(),
		storageDir:  storageDir,
		nameFileMap: nameFileMap,
	}, nil
}

// nextDataFile should be called in the end of put when offset is bigger
// than threshold, it's not concurrently save to call this function
// otherwise
func (c *cask) nextDataFile() error {
	files, err := ioutil.ReadDir(c.storageDir)
	if err != nil {
		return err
	}

	var lastesID int

	for _, file := range files {
		if file.Mode().IsRegular() && storageFileRegexp.MatchString(file.Name()) {
			sp := strings.Split(file.Name(), "_")
			sp = strings.Split(sp[len(sp)-1], ".")
			cand, err := strconv.Atoi(sp[0])

			if err != nil {
				return err
			}

			if cand > lastesID {
				lastesID = cand
			}

		}
	}

	f, err := os.OpenFile(fmt.Sprintf("%s/storage_%d.dat", c.storageDir, lastesID+1),
		os.O_CREATE|os.O_RDWR|os.O_EXCL, 0755)
	if err != nil {
		return err
	}

	fmt.Println(c.nameFileMap)
	c.nameFileMap[f.Name()] = f
	c.lastOffset = 0
	c.dataFile = f

	return nil
}

func readDataFile(file string, index map[string]keyRecord) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	var offset int64

	for {
		e, err := readRecord(f, offset)
		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}

		index[e.key] = keyRecord{
			fileID:    file,
			valueSize: e.valueSize,
			valuePos:  offset,
			ts:        e.timestamp,
		}

		offset = offset + headerSize + int64(e.keySize) + int64(e.valueSize)
	}
}

func (c *cask) get(key string) (string, error) {
	c.RLock()
	defer c.RUnlock()

	kr, ok := c.index[key]
	if !ok {
		return "", errKeyNotFount
	}

	f, ok := c.nameFileMap[kr.fileID]
	if !ok {
		panic("something is wrong")
	}

	entry, err := readRecord(f, kr.valuePos)
	if err != nil {
		return "", err
	}

	return entry.value, nil
}

func (c *cask) put(key string, value string) error {
	crc, b := prepareRecord(key, value)
	c.Lock()
	defer c.Unlock()
	err := binary.Write(c.dataFile, binary.BigEndian, crc)
	if err != nil {
		return err
	}

	n, err := c.dataFile.Write(b)
	if err != nil {
		return err
	}

	err = c.dataFile.Sync()
	if err != nil {
		return err
	}

	// TODO: fix me
	c.index[key] = keyRecord{
		fileID:    c.dataFile.Name(),
		valueSize: 0,
		valuePos:  c.lastOffset,
		ts:        0,
	}
	c.lastOffset += int64(n) + 4

	log.Println("current offset is: ", c.lastOffset)
	if c.lastOffset > 100 {
		err := c.nextDataFile()
		if err != nil {
			log.Println(err)
		}
	}

	return nil
}

func (c *cask) del(key string) error {
	c.Lock()
	defer c.Unlock()
	delete(c.index, key)
	return nil
}

func prepareRecord(key string, value string) (uint32, []byte) {
	buf := new(bytes.Buffer)

	ts := int32(time.Now().Unix())
	err := binary.Write(buf, binary.BigEndian, ts)
	if err != nil {
		panic(err)
	}

	err = binary.Write(buf, binary.BigEndian, uint32(len(key)))
	if err != nil {
		panic(err)
	}

	err = binary.Write(buf, binary.BigEndian, uint64(len(value)))
	if err != nil {
		panic(err)
	}

	_, err = buf.WriteString(key)
	if err != nil {
		panic(err)
	}

	_, err = buf.WriteString(value)
	if err != nil {
		panic(err)
	}

	b := buf.Bytes()

	return crc32.ChecksumIEEE(b), b
}

func readRecord(f *os.File, offset int64) (*entry, error) {
	b := make([]byte, headerSize)
	reader := bytes.NewReader(b)

	_, err := f.ReadAt(b, offset)
	if err == io.EOF {
		return nil, io.EOF
	}

	var crc uint32
	var ts int32
	var keySize uint32
	var valueSize uint64

	binary.Read(reader, binary.BigEndian, &crc)
	binary.Read(reader, binary.BigEndian, &ts)
	binary.Read(reader, binary.BigEndian, &keySize)
	binary.Read(reader, binary.BigEndian, &valueSize)

	kv := make([]byte, uint64(keySize)+valueSize)
	f.ReadAt(kv, offset+headerSize)

	e := entry{
		crc:       crc,
		timestamp: ts,
		keySize:   keySize,
		valueSize: valueSize,
		key:       string(kv[0:keySize]),
		value:     string(kv[keySize:]),
	}

	check := make([]byte, 0, headerSize-4+len(kv))
	check = append(check, b[4:]...)
	check = append(check, kv...)

	crcRec := crc32.ChecksumIEEE(check)
	if crcRec != e.crc {
		return nil, errCRCNotMatch
	}

	return &e, nil
}
