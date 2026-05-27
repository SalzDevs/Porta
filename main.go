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
	MsgFunctionCallResponse = 'V'
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

func data_row(values [][]byte)([]byte,error){
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
	if err:= binary.Write(&buf,binary.BigEndian,int32(length)); err!=nil {
		return nil,err
	}
	if err:= binary.Write(&buf,binary.BigEndian,int16(len(values))); err!=nil {
		return nil,err
	}
	for _, v := range values {
		if v == nil {
			if err:= binary.Write(&buf,binary.BigEndian,int32(-1)); err!=nil {
				return nil,err
			}
		} else {
			if err:= binary.Write(&buf,binary.BigEndian,int32(len(v))); err!=nil {
				return nil,err
			}
			buf.Write(v)
		}
	}
	return buf.Bytes(),nil
}

func empty_query_response()([]byte,error){
	var buf bytes.Buffer
	buf.WriteByte(MsgEmptyQueryResponse)
	if err:= binary.Write(&buf,binary.BigEndian,int32(4)); err!=nil {
		return nil,err
	}
	return buf.Bytes(),nil
}

func error_response(fields map[byte]string)([]byte,error){
	var buf bytes.Buffer
	payload_len := 1
	for _, value := range fields {
		payload_len += 1 + len(value) + 1 
	}
	length := 4 + payload_len

	buf.WriteByte(MsgErrorResponse)
	if err:= binary.Write(&buf,binary.BigEndian,int32(length)); err!=nil {
		return nil,err
	}
	for code, value := range fields {
		buf.WriteByte(code)
		buf.Write([]byte(value))
		buf.WriteByte(0)
	}
	buf.WriteByte(0) 
	return buf.Bytes(),nil
}

func function_call_response(result []byte)([]byte,error){
	var buf bytes.Buffer
	var result_len int32
	if result == nil {
		result_len = -1
	} else {
		result_len = int32(len(result))
	}
	length := 4 + 4
	if result != nil {
		length += len(result)
	}

	buf.WriteByte(MsgFunctionCallResponse)
	if err:= binary.Write(&buf,binary.BigEndian,int32(length)); err!=nil {
		return nil,err
	}
	if err:= binary.Write(&buf,binary.BigEndian,result_len); err!=nil {
		return nil,err
	}
	if result != nil {
		buf.Write(result)
	}
	return buf.Bytes(),nil
}

func main(){
	println("Hello seamen!")
}
