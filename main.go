package main

import (
	"bytes"
	"encoding/binary"
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
	MsgNoticeResponse    = 'N'
	MsgParseComplete     = '1'
	MsgBindComplete      = '2'
	MsgCloseComplete     = '3'
)

// Auth result codes
const (
	AuthOK               = 0
	AuthKerberosV5			 = 2
	AuthMD5              = 5
	AuthCleartext        = 3
	AuthSASL             = 10
	AuthSASLContinue     = 11
	AuthSASLFinal        = 12
)

// Transaction statuses for ReadyForQuery
const (
	TxnIdle             = 'I'
	TxnInTransaction    = 'T'
	TxnFailed           = 'E'
)

func authentication_ok()([]byte,error){
	var buf bytes.Buffer 
	buf.WriteByte(MsgAuthentication)
	if err:= binary.Write(&buf,binary.BigEndian,int32(8)); err!=nil {
		return nil,err
	}
	if err:= binary.Write(&buf,binary.BigEndian,int32(AuthOK)); err!=nil {
		return nil,err
	}
	return buf.Bytes(),nil
}

func authentication_kerberos_v5() ([]byte,error){
	var buf bytes.Buffer 
	buf.WriteByte(MsgAuthentication)
	if err:= binary.Write(&buf,binary.BigEndian,int32(8)); err!=nil {
		return nil,err
	}
	if err:= binary.Write(&buf,binary.BigEndian,int32(AuthKerberosV5)); err!=nil {
		return nil,err
	}
	return buf.Bytes(),nil	
}

func main(){
	println("Hello seamen!")
}
