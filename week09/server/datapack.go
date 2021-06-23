package server

import (
	"bytes"
	"encoding/binary"
	"io"
	"net"
)

const PackageHeaderLen uint16 = 4 + 2 + 2 + 4 + 4

type PackageHeader struct {
	PackageLen      uint32
	HeaderLen       uint16
	ProtocolVersion uint16
	Operation       uint32
	SequenceId      uint32
}

func NewPackageHeader(packageLen uint32, protocolVersion uint16, operation uint32, sequenceId uint32) *PackageHeader {
	return &PackageHeader{
		PackageLen:      packageLen,
		HeaderLen:       PackageHeaderLen,
		ProtocolVersion: protocolVersion,
		Operation:       operation,
		SequenceId:      sequenceId,
	}
}

type PackageBody struct {
	Body []byte
}

func NewPackageBody(body []byte) *PackageBody {
	return &PackageBody{
		Body: body,
	}
}

type Message struct {
	header *PackageHeader
	body   *PackageBody
}

func NewMessage(protocolVersion uint16, operation uint32, sequenceId uint32, body []byte) *Message {
	return &Message{
		header: NewPackageHeader(uint32(PackageHeaderLen)+uint32(len(body)), protocolVersion, operation, sequenceId),
		body:   NewPackageBody(body),
	}
}

func (m *Message) GetHeaderLen() uint16 {
	return PackageHeaderLen
}

func (m *Message) GetPackageLen() uint32 {
	return m.header.PackageLen
}

func (m *Message) GetBodyLen() int {
	return len(m.body.Body)
}

func Pack(msg *Message) ([]byte, error) {
	dataBuff := bytes.NewBuffer([]byte{})
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.header.PackageLen); err != nil {
		return nil, err
	}
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.header.HeaderLen); err != nil {
		return nil, err
	}
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.header.ProtocolVersion); err != nil {
		return nil, err
	}
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.header.Operation); err != nil {
		return nil, err
	}
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.header.SequenceId); err != nil {
		return nil, err
	}
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.body.Body); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

func Unpack(conn net.Conn) (*Message, error) {
	packageHeader := make([]byte, PackageHeaderLen)
	_, err := io.ReadFull(conn, packageHeader)
	if err != nil {
		return nil, err
	}

	dataBuff := bytes.NewReader(packageHeader)
	var (
		packageLen      uint32
		headerLen       uint16
		protocolVersion uint16
		operation       uint32
		sequenceId      uint32
		body            []byte
	)
	if err := binary.Read(dataBuff, binary.LittleEndian, &packageLen); err != nil {
		return nil, err
	}
	if err := binary.Read(dataBuff, binary.LittleEndian, &headerLen); err != nil {
		return nil, err
	}
	if err := binary.Read(dataBuff, binary.LittleEndian, &protocolVersion); err != nil {
		return nil, err
	}
	if err := binary.Read(dataBuff, binary.LittleEndian, &operation); err != nil {
		return nil, err
	}
	if err := binary.Read(dataBuff, binary.LittleEndian, &sequenceId); err != nil {
		return nil, err
	}
	bodyLen := packageLen - uint32(headerLen)
	if bodyLen > 0 {
		body = make([]byte, bodyLen)
		if _, err := io.ReadFull(conn, body); err != nil {
			return nil, err
		}
	}
	msg := NewMessage(protocolVersion, operation, sequenceId, body)

	return msg, nil
}
