package tftp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"strings"
)

const (
	DatagramSize = 516
	BlockSize    = DatagramSize - 4
)

type OpCode uint16

const (
	OpRRQ OpCode = iota + 1
	_
	OpData
	OpAck
	OpErr
)

type ErrCode uint16

const (
	ErrUnknown ErrCode = iota
	ErrNotFount
	ErrAccessViolation
	ErrDiskFull
	ErrIllegalOp
	ErrUnknownID
	ErrNoUser
)

// ReadReq packet structure
// 2 byte		n byte		1 byte		n byte		1 byte
// OpCode		Filename	0			Mode		0
type ReadReq struct {
	Filename string
	Mode     string
}

func (req ReadReq) MarshalBinary() ([]byte, error) {
	mode := "octet"
	if req.Mode != "" {
		mode = req.Mode
	}

	b := new(bytes.Buffer)
	b.Grow(2 + 2 + len(req.Filename) + 1 + len(req.Mode) + 1)

	err := binary.Write(b, binary.BigEndian, OpRRQ)
	if err != nil {
		return nil, err
	}

	_, err = b.WriteString(req.Filename)
	if err != nil {
		return nil, err
	}

	err = b.WriteByte(0)
	if err != nil {
		return nil, err
	}

	_, err = b.WriteString(mode)
	if err != nil {
		return nil, err
	}

	err = b.WriteByte(0)
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func (req ReadReq) UnmarshalBinary(p []byte) error {
	r := bytes.NewBuffer(p)

	var code OpCode
	err := binary.Read(r, binary.BigEndian, &code)
	if err != nil {
		return err
	}

	if code != OpRRQ {
		return fmt.Errorf("invalid RRQ, code")
	}

	req.Filename, err = r.ReadString(0)
	if err != nil {
		return fmt.Errorf("ivalid RRQ, filename")
	}

	req.Filename = strings.TrimRight(req.Filename, "\x00")
	if len(req.Filename) == 0 {
		return fmt.Errorf("invalid RRQ, trim filename")
	}

	req.Mode, err = r.ReadString(0)
	if err != nil {
		return fmt.Errorf("ivalid RRQ, mode")
	}

	req.Mode = strings.TrimRight(req.Mode, "\x00")
	if len(req.Mode) == 0 {
		return fmt.Errorf("invalid RRQ, trim mode")
	}

	actual := strings.ToLower(req.Mode)
	if actual != "octet" {
		return fmt.Errorf("only binary transfers supported")
	}

	return nil
}

// Data packet structure
// 2 bytes		2 bytes		n bytes
// OpCode		Block #		Payload
type Data struct {
	Block   uint16
	Payload io.Reader
}

func (d *Data) MarshalBinary() ([]byte, error) {
	b := new(bytes.Buffer)
	b.Grow(DatagramSize)

	d.Block++

	err := binary.Write(b, binary.BigEndian, OpData)
	if err != nil {
		return nil, err
	}

	err = binary.Write(b, binary.BigEndian, d.Block)
	if err != nil {
		return nil, err
	}

	_, err = io.CopyN(b, d.Payload, BlockSize)
	if err != nil && err != io.EOF {
		return nil, err
	}

	return b.Bytes(), nil
}

func (d *Data) UnmarshallBinary(p []byte) error {
	if l := len(p); l < 4 || l > DatagramSize {
		return fmt.Errorf("invalid DATA")
	}

	var opcode OpCode

	err := binary.Read(bytes.NewReader(p[:2]), binary.BigEndian, &opcode)
	if err != nil || opcode != OpData {
		return fmt.Errorf("invalid DATA")
	}

	err = binary.Read(bytes.NewReader(p[2:4]), binary.BigEndian, &d.Block) // & needed?
	if err != nil {
		return fmt.Errorf("invalid DATA")
	}

	d.Payload = bytes.NewBuffer(p[4:])

	return nil
}

// Ack packet structure
// 2 bytes		2 bytes
// OpCode		Block #
type Ack uint16

func (a Ack) MarshalBinary() ([]byte, error) {
	b := new(bytes.Buffer)
	b.Grow(2 + 2)

	err := binary.Write(b, binary.BigEndian, OpAck)
	if err != nil {
		return nil, err
	}

	err = binary.Write(b, binary.BigEndian, a)
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func (a *Ack) UnmarshalBinary(p []byte) error {
	var code OpCode
	r := bytes.NewReader(p)

	err := binary.Read(r, binary.BigEndian, &code)
	if err != nil {
		return err
	}

	if code != OpAck {
		return fmt.Errorf("invalid Ack")
	}

	return binary.Read(r, binary.BigEndian, a)
}

// Err packet structure
// 2 bytes		2 bytes		n bytes		1 byte
// OpCode		ErrCode		Message		0
type Err struct {
	Error   ErrCode
	Message string
}

func (e Err) MarshalBinary() ([]byte, error) {
	b := new(bytes.Buffer)
	b.Grow(2 + 2 + len(e.Message) + 1)

	err := binary.Write(b, binary.BigEndian, OpErr)
	if err != nil {
		return nil, err
	}

	err = binary.Write(b, binary.BigEndian, e.Error)
	if err != nil {
		return nil, err
	}

	_, err = b.WriteString(e.Message)
	if err != nil {
		return nil, err
	}

	err = b.WriteByte(0)
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func (e *Err) UnmarshalBinary(p []byte) error {
	r := bytes.NewBuffer(p)

	var code OpCode

	err := binary.Read(r, binary.BigEndian, &code)
	if err != nil {
		return err
	}

	if code != OpErr {
		return fmt.Errorf("invalid ERROR")
	}

	err = binary.Read(r, binary.BigEndian, &e.Error)
	if err != nil {
		return err
	}

	e.Message, err = r.ReadString(0)
	e.Message = strings.TrimRight(e.Message, "\x00")

	return err
}
