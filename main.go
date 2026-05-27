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
	AuthCleartext      	 = 3
	AuthMD5              = 5
	AuthGSS 						 = 7 
	AuthGSSContinue 		 = 8
	AuthSSPI						 = 9
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

func authentication_clear_text_password()([]byte,error){
	var buf bytes.Buffer 
	buf.WriteByte(MsgAuthentication)
	if err:= binary.Write(&buf,binary.BigEndian,int32(8)); err!=nil {
		return nil,err
	}
	if err:= binary.Write(&buf,binary.BigEndian,int32(AuthCleartext)); err!=nil {
		return nil,err
	}
	return buf.Bytes(),nil	
}

func authentication_md5_password(salt [4]byte)([]byte,error){
	var buf bytes.Buffer 
	buf.WriteByte(MsgAuthentication)
	if err:= binary.Write(&buf,binary.BigEndian,int32(12)); err!=nil {
		return nil,err
	}
	if err:= binary.Write(&buf,binary.BigEndian,int32(AuthMD5)); err!=nil {
		return nil,err
	}
	buf.Write(salt[:])	
	return buf.Bytes(),nil	
}

func authentication_gss()([]byte,error){
	var buf bytes.Buffer 
	buf.WriteByte(MsgAuthentication)
	if err:= binary.Write(&buf,binary.BigEndian,int32(8)); err!=nil {
		return nil,err
	}
	if err:= binary.Write(&buf,binary.BigEndian,int32(AuthGSS)); err!=nil {
		return nil,err
	}
	return buf.Bytes(),nil	
}

func authentication_gss_continue(gssapi_or_sspi_data []byte)([]byte,error){
	length := len(gssapi_or_sspi_data) + 8
	var buf bytes.Buffer   
	buf.WriteByte(MsgAuthentication)
	if err:= binary.Write(&buf,binary.BigEndian,int32(length)); err!=nil {
		return nil,err
	}
	if err:= binary.Write(&buf,binary.BigEndian,int32(AuthGSSContinue)); err!=nil {
		return nil,err
	}		
	buf.Write(gssapi_or_sspi_data[:])
	return buf.Bytes(),nil	
}

func authentication_sspi()([]byte,error){
	var buf bytes.Buffer   
	buf.WriteByte(MsgAuthentication)
	if err:= binary.Write(&buf,binary.BigEndian,int32(8)); err!=nil {
		return nil,err
	}
	if err:= binary.Write(&buf,binary.BigEndian,int32(AuthSSPI)); err!=nil {
		return nil,err
	}
	return buf.Bytes(),nil		
}

func authentication_sasl(name_of_sals_auth_mechanism string)([]byte,error){
	var buf bytes.Buffer  
	length := len(name_of_sals_auth_mechanism) + 10
	buf.WriteByte(MsgAuthentication)
	if err:= binary.Write(&buf,binary.BigEndian,int32(length)); err!=nil {
		return nil,err
	}
	if err:= binary.Write(&buf,binary.BigEndian,int32(AuthSASL)); err!=nil {
		return nil,err
	}
	buf.Write([]byte(name_of_sals_auth_mechanism))
	buf.WriteByte(0)
	buf.WriteByte(0)
	return buf.Bytes(),nil		
}

func main(){
	println("Hello seamen!")
}
