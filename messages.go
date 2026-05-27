package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// Message type bytes
const (
	MsgQuery             = 'Q'
	MsgTerminate         = 'X'
	MsgPasswordMessage   = 'p'

	MsgAuthentication    = 'R'
	MsgReadyForQuery     = 'Z'
	MsgErrorResponse     = 'E'
	MsgBackendKeyData    = 'K'
	MsgParameterStatus   = 'S'
	MsgCommandComplete   = 'C'
	MsgDataRow           = 'D'
	MsgRowDescription    = 'T'
	MsgEmptyQueryResponse = 'I'
	MsgNoData            = 'n'
	MsgNoticeResponse    = 'N'
	MsgParseComplete     = '1'
	MsgBindComplete      = '2'
	MsgCloseComplete     = '3'
	MsgPortalSuspended   = 's'
)

// Auth result codes
const (
	AuthOK  = 0
	AuthMD5 = 5
)

// Transaction statuses for ReadyForQuery
const (
	TxnIdle          = 'I'
	TxnInTransaction = 'T'
	TxnFailed        = 'E'
)

func authentication_ok() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte(MsgAuthentication)
	if err := binary.Write(&buf, binary.BigEndian, int32(8)); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.BigEndian, int32(AuthOK)); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func authentication_md5_password(salt [4]byte) ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte(MsgAuthentication)
	if err := binary.Write(&buf, binary.BigEndian, int32(12)); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.BigEndian, int32(AuthMD5)); err != nil {
		return nil, err
	}
	buf.Write(salt[:])
	return buf.Bytes(), nil
}

func backend_key_data(secret_key []byte, process_id int32) ([]byte, error) {
	var buf bytes.Buffer
	length := len(secret_key) + 8
	buf.WriteByte(MsgBackendKeyData)
	if err := binary.Write(&buf, binary.BigEndian, int32(length)); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.BigEndian, process_id); err != nil {
		return nil, err
	}
	buf.Write([]byte(secret_key))
	return buf.Bytes(), nil
}

func ready_for_query(status byte) ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte(MsgReadyForQuery)
	if err := binary.Write(&buf, binary.BigEndian, int32(5)); err != nil {
		return nil, err
	}
	buf.WriteByte(status)
	return buf.Bytes(), nil
}

func command_complete(command_tag string) ([]byte, error) {
	var buf bytes.Buffer
	length := len(command_tag) + 5
	buf.WriteByte(MsgCommandComplete)
	if err := binary.Write(&buf, binary.BigEndian, int32(length)); err != nil {
		return nil, err
	}
	buf.Write([]byte(command_tag))
	buf.WriteByte(0)
	return buf.Bytes(), nil
}

func error_response(fields map[byte]string) ([]byte, error) {
	var buf bytes.Buffer
	payload_len := 1
	for _, value := range fields {
		payload_len += 1 + len(value) + 1
	}
	length := 4 + payload_len

	buf.WriteByte(MsgErrorResponse)
	if err := binary.Write(&buf, binary.BigEndian, int32(length)); err != nil {
		return nil, err
	}
	for code, value := range fields {
		buf.WriteByte(code)
		buf.Write([]byte(value))
		buf.WriteByte(0)
	}
	buf.WriteByte(0)
	return buf.Bytes(), nil
}

func notice_response(fields map[byte]string) ([]byte, error) {
	var buf bytes.Buffer
	payload_len := 1
	for _, value := range fields {
		payload_len += 1 + len(value) + 1
	}
	length := 4 + payload_len

	buf.WriteByte(MsgNoticeResponse)
	if err := binary.Write(&buf, binary.BigEndian, int32(length)); err != nil {
		return nil, err
	}
	for code, value := range fields {
		buf.WriteByte(code)
		buf.Write([]byte(value))
		buf.WriteByte(0)
	}
	buf.WriteByte(0)
	return buf.Bytes(), nil
}

func parameter_status(name string, value string) ([]byte, error) {
	var buf bytes.Buffer
	length := 6 + len(name) + len(value)

	buf.WriteByte(MsgParameterStatus)
	if err := binary.Write(&buf, binary.BigEndian, int32(length)); err != nil {
		return nil, err
	}
	buf.Write([]byte(name))
	buf.WriteByte(0)
	buf.Write([]byte(value))
	buf.WriteByte(0)
	return buf.Bytes(), nil
}

func empty_query_response() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte(MsgEmptyQueryResponse)
	if err := binary.Write(&buf, binary.BigEndian, int32(4)); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func no_data() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte(MsgNoData)
	if err := binary.Write(&buf, binary.BigEndian, int32(4)); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func parse_complete() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte(MsgParseComplete)
	if err := binary.Write(&buf, binary.BigEndian, int32(4)); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func bind_complete() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte(MsgBindComplete)
	if err := binary.Write(&buf, binary.BigEndian, int32(4)); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func close_complete() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte(MsgCloseComplete)
	if err := binary.Write(&buf, binary.BigEndian, int32(4)); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func portal_suspended() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte(MsgPortalSuspended)
	if err := binary.Write(&buf, binary.BigEndian, int32(4)); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

type ColumnInfo struct {
	Name         string
	TableOID     int32
	ColumnAttr   int16
	TypeOID      int32
	TypeSize     int16
	TypeModifier int32
	FormatCode   int16
}

func row_description(columns []ColumnInfo) ([]byte, error) {
	var buf bytes.Buffer
	payload_len := 2
	for _, c := range columns {
		payload_len += len(c.Name) + 1 + 4 + 2 + 4 + 2 + 4 + 2
	}
	length := 4 + payload_len

	buf.WriteByte(MsgRowDescription)
	if err := binary.Write(&buf, binary.BigEndian, int32(length)); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.BigEndian, int16(len(columns))); err != nil {
		return nil, err
	}
	for _, c := range columns {
		buf.Write([]byte(c.Name))
		buf.WriteByte(0)
		if err := binary.Write(&buf, binary.BigEndian, c.TableOID); err != nil {
			return nil, err
		}
		if err := binary.Write(&buf, binary.BigEndian, c.ColumnAttr); err != nil {
			return nil, err
		}
		if err := binary.Write(&buf, binary.BigEndian, c.TypeOID); err != nil {
			return nil, err
		}
		if err := binary.Write(&buf, binary.BigEndian, c.TypeSize); err != nil {
			return nil, err
		}
		if err := binary.Write(&buf, binary.BigEndian, c.TypeModifier); err != nil {
			return nil, err
		}
		if err := binary.Write(&buf, binary.BigEndian, c.FormatCode); err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

func data_row(values [][]byte) ([]byte, error) {
	var buf bytes.Buffer
	payload_len := 2
	for _, v := range values {
		if v == nil {
			payload_len += 4
		} else {
			payload_len += 4 + len(v)
		}
	}
	length := 4 + payload_len

	buf.WriteByte(MsgDataRow)
	if err := binary.Write(&buf, binary.BigEndian, int32(length)); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.BigEndian, int16(len(values))); err != nil {
		return nil, err
	}
	for _, v := range values {
		if v == nil {
			if err := binary.Write(&buf, binary.BigEndian, int32(-1)); err != nil {
				return nil, err
			}
		} else {
			if err := binary.Write(&buf, binary.BigEndian, int32(len(v))); err != nil {
				return nil, err
			}
			buf.Write(v)
		}
	}
	return buf.Bytes(), nil
}

func parse_startup(data []byte) (string, string, int, error) {
	if len(data) < 8 {
		return "", "", 0, fmt.Errorf("message too short")
	}
	length := int(binary.BigEndian.Uint32(data[0:4]))
	if len(data) < length {
		return "", "", 0, fmt.Errorf("message incomplete")
	}
	protocol_version := binary.BigEndian.Uint32(data[4:8])
	if protocol_version != 196608 {
		return "", "", 0, fmt.Errorf("unsupported protocol version")
	}

	var user, database string
	params := data[8:length]
	i := 0
	for i < len(params) && params[i] != 0 {
		key_end := i + bytes.IndexByte(params[i:], 0)
		key := string(params[i:key_end])
		i = key_end + 1

		val_end := i + bytes.IndexByte(params[i:], 0)
		value := string(params[i:val_end])
		i = val_end + 1

		switch key {
		case "user":
			user = value
		case "database":
			database = value
		}
	}

	return user, database, length, nil
}

func parse_query(data []byte) (string, error) {
	if len(data) < 5 {
		return "", fmt.Errorf("message too short")
	}
	if data[0] != 'Q' {
		return "", fmt.Errorf("wrong message type")
	}
	length := int(binary.BigEndian.Uint32(data[1:5]))
	if len(data) < 1+length {
		return "", fmt.Errorf("message incomplete")
	}
	payload := data[5 : 1+length]
	if len(payload) == 0 || payload[len(payload)-1] != 0 {
		return "", fmt.Errorf("invalid format")
	}
	return string(payload[:len(payload)-1]), nil
}
