// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package main

import (
	json "encoding/json"
	gin "github.com/gin-gonic/gin"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson89aae3efDecodeGithubComPolygonIoErrandsServer(in *jlexer.Lexer, out *Notification) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "event":
			out.Event = string(in.String())
		case "errand":
			(out.Errand).UnmarshalEasyJSON(in)
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson89aae3efEncodeGithubComPolygonIoErrandsServer(out *jwriter.Writer, in Notification) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"event\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Event))
	}
	if true {
		const prefix string = ",\"errand\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		(in.Errand).MarshalEasyJSON(out)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Notification) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson89aae3efEncodeGithubComPolygonIoErrandsServer(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Notification) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson89aae3efEncodeGithubComPolygonIoErrandsServer(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Notification) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson89aae3efDecodeGithubComPolygonIoErrandsServer(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Notification) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson89aae3efDecodeGithubComPolygonIoErrandsServer(l, v)
}
func easyjson89aae3efDecodeGithubComPolygonIoErrandsServer1(in *jlexer.Lexer, out *Log) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "severity":
			out.Severity = string(in.String())
		case "message":
			out.Message = string(in.String())
		case "timestamp":
			out.Timestamp = int64(in.Int64())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson89aae3efEncodeGithubComPolygonIoErrandsServer1(out *jwriter.Writer, in Log) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"severity\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Severity))
	}
	{
		const prefix string = ",\"message\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Message))
	}
	{
		const prefix string = ",\"timestamp\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int64(int64(in.Timestamp))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Log) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson89aae3efEncodeGithubComPolygonIoErrandsServer1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Log) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson89aae3efEncodeGithubComPolygonIoErrandsServer1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Log) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson89aae3efDecodeGithubComPolygonIoErrandsServer1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Log) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson89aae3efDecodeGithubComPolygonIoErrandsServer1(l, v)
}
func easyjson89aae3efDecodeGithubComPolygonIoErrandsServer2(in *jlexer.Lexer, out *Errand) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			out.ID = string(in.String())
		case "name":
			out.Name = string(in.String())
		case "type":
			out.Type = string(in.String())
		case "options":
			easyjson89aae3efDecode(in, &out.Options)
		case "data":
			if in.IsNull() {
				in.Skip()
				out.Data = nil
			} else {
				if out.Data == nil {
					out.Data = new(gin.H)
				}
				if in.IsNull() {
					in.Skip()
				} else {
					in.Delim('{')
					if !in.IsDelim('}') {
						*out.Data = make(gin.H)
					} else {
						*out.Data = nil
					}
					for !in.IsDelim('}') {
						key := string(in.String())
						in.WantColon()
						var v1 interface{}
						if m, ok := v1.(easyjson.Unmarshaler); ok {
							m.UnmarshalEasyJSON(in)
						} else if m, ok := v1.(json.Unmarshaler); ok {
							_ = m.UnmarshalJSON(in.Raw())
						} else {
							v1 = in.Interface()
						}
						(*out.Data)[key] = v1
						in.WantComma()
					}
					in.Delim('}')
				}
			}
		case "created":
			out.Created = int64(in.Int64())
		case "status":
			out.Status = string(in.String())
		case "results":
			if in.IsNull() {
				in.Skip()
				out.Results = nil
			} else {
				if out.Results == nil {
					out.Results = new(gin.H)
				}
				if in.IsNull() {
					in.Skip()
				} else {
					in.Delim('{')
					if !in.IsDelim('}') {
						*out.Results = make(gin.H)
					} else {
						*out.Results = nil
					}
					for !in.IsDelim('}') {
						key := string(in.String())
						in.WantColon()
						var v2 interface{}
						if m, ok := v2.(easyjson.Unmarshaler); ok {
							m.UnmarshalEasyJSON(in)
						} else if m, ok := v2.(json.Unmarshaler); ok {
							_ = m.UnmarshalJSON(in.Raw())
						} else {
							v2 = in.Interface()
						}
						(*out.Results)[key] = v2
						in.WantComma()
					}
					in.Delim('}')
				}
			}
		case "progress":
			out.Progress = float64(in.Float64())
		case "attempts":
			out.Attempts = int(in.Int())
		case "started":
			out.Started = int64(in.Int64())
		case "failed":
			out.Failed = int64(in.Int64())
		case "compelted":
			out.Completed = int64(in.Int64())
		case "logs":
			if in.IsNull() {
				in.Skip()
				out.Logs = nil
			} else {
				in.Delim('[')
				if out.Logs == nil {
					if !in.IsDelim(']') {
						out.Logs = make([]Log, 0, 1)
					} else {
						out.Logs = []Log{}
					}
				} else {
					out.Logs = (out.Logs)[:0]
				}
				for !in.IsDelim(']') {
					var v3 Log
					(v3).UnmarshalEasyJSON(in)
					out.Logs = append(out.Logs, v3)
					in.WantComma()
				}
				in.Delim(']')
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson89aae3efEncodeGithubComPolygonIoErrandsServer2(out *jwriter.Writer, in Errand) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.ID))
	}
	{
		const prefix string = ",\"name\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Name))
	}
	{
		const prefix string = ",\"type\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Type))
	}
	{
		const prefix string = ",\"options\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		easyjson89aae3efEncode(out, in.Options)
	}
	if in.Data != nil {
		const prefix string = ",\"data\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if *in.Data == nil && (out.Flags&jwriter.NilMapAsEmpty) == 0 {
			out.RawString(`null`)
		} else {
			out.RawByte('{')
			v4First := true
			for v4Name, v4Value := range *in.Data {
				if v4First {
					v4First = false
				} else {
					out.RawByte(',')
				}
				out.String(string(v4Name))
				out.RawByte(':')
				if m, ok := v4Value.(easyjson.Marshaler); ok {
					m.MarshalEasyJSON(out)
				} else if m, ok := v4Value.(json.Marshaler); ok {
					out.Raw(m.MarshalJSON())
				} else {
					out.Raw(json.Marshal(v4Value))
				}
			}
			out.RawByte('}')
		}
	}
	{
		const prefix string = ",\"created\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int64(int64(in.Created))
	}
	if in.Status != "" {
		const prefix string = ",\"status\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Status))
	}
	if in.Results != nil {
		const prefix string = ",\"results\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if *in.Results == nil && (out.Flags&jwriter.NilMapAsEmpty) == 0 {
			out.RawString(`null`)
		} else {
			out.RawByte('{')
			v5First := true
			for v5Name, v5Value := range *in.Results {
				if v5First {
					v5First = false
				} else {
					out.RawByte(',')
				}
				out.String(string(v5Name))
				out.RawByte(':')
				if m, ok := v5Value.(easyjson.Marshaler); ok {
					m.MarshalEasyJSON(out)
				} else if m, ok := v5Value.(json.Marshaler); ok {
					out.Raw(m.MarshalJSON())
				} else {
					out.Raw(json.Marshal(v5Value))
				}
			}
			out.RawByte('}')
		}
	}
	{
		const prefix string = ",\"progress\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Float64(float64(in.Progress))
	}
	{
		const prefix string = ",\"attempts\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.Attempts))
	}
	if in.Started != 0 {
		const prefix string = ",\"started\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int64(int64(in.Started))
	}
	if in.Failed != 0 {
		const prefix string = ",\"failed\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int64(int64(in.Failed))
	}
	if in.Completed != 0 {
		const prefix string = ",\"compelted\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int64(int64(in.Completed))
	}
	if len(in.Logs) != 0 {
		const prefix string = ",\"logs\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		{
			out.RawByte('[')
			for v6, v7 := range in.Logs {
				if v6 > 0 {
					out.RawByte(',')
				}
				(v7).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Errand) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson89aae3efEncodeGithubComPolygonIoErrandsServer2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Errand) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson89aae3efEncodeGithubComPolygonIoErrandsServer2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Errand) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson89aae3efDecodeGithubComPolygonIoErrandsServer2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Errand) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson89aae3efDecodeGithubComPolygonIoErrandsServer2(l, v)
}
func easyjson89aae3efDecode(in *jlexer.Lexer, out *struct {
	TTL               int  `json:"ttl,omitempty"`
	Retries           int  `json:"retries,omitempty"`
	Priority          int  `json:"priority,omitempty"`
	DeleteOnCompleted bool `json:"deleteOnCompleted,omitempty"`
}) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "ttl":
			out.TTL = int(in.Int())
		case "retries":
			out.Retries = int(in.Int())
		case "priority":
			out.Priority = int(in.Int())
		case "deleteOnCompleted":
			out.DeleteOnCompleted = bool(in.Bool())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson89aae3efEncode(out *jwriter.Writer, in struct {
	TTL               int  `json:"ttl,omitempty"`
	Retries           int  `json:"retries,omitempty"`
	Priority          int  `json:"priority,omitempty"`
	DeleteOnCompleted bool `json:"deleteOnCompleted,omitempty"`
}) {
	out.RawByte('{')
	first := true
	_ = first
	if in.TTL != 0 {
		const prefix string = ",\"ttl\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.TTL))
	}
	if in.Retries != 0 {
		const prefix string = ",\"retries\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.Retries))
	}
	if in.Priority != 0 {
		const prefix string = ",\"priority\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.Priority))
	}
	if in.DeleteOnCompleted {
		const prefix string = ",\"deleteOnCompleted\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Bool(bool(in.DeleteOnCompleted))
	}
	out.RawByte('}')
}