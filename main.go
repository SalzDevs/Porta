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
	MsgCopyData					 = 'd'
	MsgCopyDone					 = 'c'
	MsgCopyInResponse    = 'G'
	MsgCopyOutResponse   = 'H'
	MsgCopyBothResponse  = 'W'
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

func authentication_sasl_continue(sasl_data []byte)([]byte,error){
	var buf bytes.Buffer  
	length := len(sasl_data) + 8 
	buf.WriteByte(MsgAuthentication)
	if err:= binary.Write(&buf,binary.BigEndian,int32(length)); err!=nil {
		return nil,err
	}
	if err:= binary.Write(&buf,binary.BigEndian,int32(AuthSASLContinue)); err!=nil {
		return nil,err
	}
	buf.Write([]byte(sasl_data))
	return buf.Bytes(),nil
}

func authentication_sasl_final(sasl_additional_data []byte)([]byte,error){
	var buf bytes.Buffer  
	length := len(sasl_additional_data) + 8 
	buf.WriteByte(MsgAuthentication)
	if err:= binary.Write(&buf,binary.BigEndian,int32(length)); err!=nil {
		return nil,err
	}
	if err:= binary.Write(&buf,binary.BigEndian,int32(AuthSASLFinal)); err!=nil {
		return nil,err
	}
	buf.Write([]byte(sasl_additional_data))
	return buf.Bytes(),nil
}

func backend_key_data(secret_key []byte, process_id int32)([]byte,error){
	var buf bytes.Buffer  
	length := len(secret_key) + 8 
	buf.WriteByte(MsgBackendKeyData)
	if err:= binary.Write(&buf,binary.BigEndian,int32(length)); err!=nil {
		return nil,err
	}
	if err:= binary.Write(&buf,binary.BigEndian,process_id); err!=nil {
		return nil,err
	}
	buf.Write([]byte(secret_key))
	return buf.Bytes(),nil
}

func bind_complete()([]byte,error){
	var buf bytes.Buffer  
	buf.WriteByte(MsgBindComplete)
	if err:= binary.Write(&buf,binary.BigEndian,int32(4)); err!=nil {
		return nil,err
	}
	return buf.Bytes(),nil
}

func close_complete()([]byte,error){
	var buf bytes.Buffer  
	buf.WriteByte(MsgCloseComplete)
	if err:= binary.Write(&buf,binary.BigEndian,int32(4)); err!=nil {
		return nil,err
	}
	return buf.Bytes(),nil
}

func command_complete(command_tag string)([]byte,error){
	var buf bytes.Buffer 
	length := len(command_tag) + 5 
	buf.WriteByte(MsgCommandComplete)
	if err:= binary.Write(&buf,binary.BigEndian,int32(length)); err!=nil {
		return nil,err
	}
	buf.Write([]byte(command_tag))
	buf.WriteByte(0)
	return buf.Bytes(),nil
}

func copy_data(data []byte)([]byte,error){
	var buf bytes.Buffer 
	length := len(data) + 4 
	buf.WriteByte(MsgCopyData)
	if err:= binary.Write(&buf,binary.BigEndian,int32(length)); err!=nil {
		return nil,err
	}
	buf.Write(data)
	return buf.Bytes(),nil
}

func copy_done()([]byte,error){
	var buf bytes.Buffer 
	buf.WriteByte(MsgCopyDone)
	if err:= binary.Write(&buf,binary.BigEndian,int32(4)); err!=nil {
		return nil,err
	}
	return buf.Bytes(),nil
}

func copy_in_response(overall_format int8, column_formats []int16)([]byte,error){
	var buf bytes.Buffer
	length := 7 + len(column_formats)*2 
	buf.WriteByte(MsgCopyInResponse)
	if err:= binary.Write(&buf,binary.BigEndian,int32(length)); err!=nil {
		return nil,err
	}
	if err:= binary.Write(&buf,binary.BigEndian,overall_format); err!=nil {
		return nil,err
	}
	if err:= binary.Write(&buf,binary.BigEndian,int16(len(column_formats))); err!=nil {
		return nil,err
	}

	for _, fc := range column_formats {
		if err:= binary.Write(&buf,binary.BigEndian,fc); err!=nil {
			return nil,err
		}
	}
	return buf.Bytes(),nil
}

func copy_out_response(overall_format int8, column_formats []int16)([]byte,error){
	var buf bytes.Buffer
	length := 7 + len(column_formats)*2 
	buf.WriteByte(MsgCopyOutResponse)
	if err:= binary.Write(&buf,binary.BigEndian,int32(length)); err!=nil {
		return nil,err
	}
	if err:= binary.Write(&buf,binary.BigEndian,overall_format); err!=nil {
		return nil,err
	}
	if err:= binary.Write(&buf,binary.BigEndian,int16(len(column_formats))); err!=nil {
		return nil,err
	}

	for _, fc := range column_formats {
		if err:= binary.Write(&buf,binary.BigEndian,fc); err!=nil {
			return nil,err
		}
	}

	return buf.Bytes(),nil
}

func copy_both_response(overall_format int8, column_formats []int16)([]byte,error){
	var buf bytes.Buffer
	length := 7 + len(column_formats)*2 
	buf.WriteByte(MsgCopyBothResponse)
	if err:= binary.Write(&buf,binary.BigEndian,int32(length)); err!=nil {
		return nil,err
	}
	if err:= binary.Write(&buf,binary.BigEndian,overall_format); err!=nil {
		return nil,err
	}
	if err:= binary.Write(&buf,binary.BigEndian,int16(len(column_formats))); err!=nil {
		return nil,err
	}

	for _, fc := range column_formats {
		if err:= binary.Write(&buf,binary.BigEndian,fc); err!=nil {
			return nil,err
		}
	}

	return buf.Bytes(),nil	
}

func main(){
	println("Hello seamen!")
}
