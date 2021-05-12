// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package schemas

import (
	json "encoding/json"
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

func easyjson2189435aDecodeGithubComPolygonIoErrandsServerSchemas(in *jlexer.Lexer, out *PipelineDependency) {
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
		case "target":
			out.Target = string(in.String())
		case "dependsOn":
			out.DependsOn = string(in.String())
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
func easyjson2189435aEncodeGithubComPolygonIoErrandsServerSchemas(out *jwriter.Writer, in PipelineDependency) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"target\":"
		out.RawString(prefix[1:])
		out.String(string(in.Target))
	}
	{
		const prefix string = ",\"dependsOn\":"
		out.RawString(prefix)
		out.String(string(in.DependsOn))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v PipelineDependency) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson2189435aEncodeGithubComPolygonIoErrandsServerSchemas(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v PipelineDependency) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson2189435aEncodeGithubComPolygonIoErrandsServerSchemas(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *PipelineDependency) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson2189435aDecodeGithubComPolygonIoErrandsServerSchemas(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *PipelineDependency) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson2189435aDecodeGithubComPolygonIoErrandsServerSchemas(l, v)
}
func easyjson2189435aDecodeGithubComPolygonIoErrandsServerSchemas1(in *jlexer.Lexer, out *Pipeline) {
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
		case "name":
			out.Name = string(in.String())
		case "deleteOnCompleted":
			out.DeleteOnCompleted = bool(in.Bool())
		case "errands":
			if in.IsNull() {
				in.Skip()
				out.Errands = nil
			} else {
				in.Delim('[')
				if out.Errands == nil {
					if !in.IsDelim(']') {
						out.Errands = make([]*Errand, 0, 8)
					} else {
						out.Errands = []*Errand{}
					}
				} else {
					out.Errands = (out.Errands)[:0]
				}
				for !in.IsDelim(']') {
					var v1 *Errand
					if in.IsNull() {
						in.Skip()
						v1 = nil
					} else {
						if v1 == nil {
							v1 = new(Errand)
						}
						(*v1).UnmarshalEasyJSON(in)
					}
					out.Errands = append(out.Errands, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "dependencies":
			if in.IsNull() {
				in.Skip()
				out.Dependencies = nil
			} else {
				in.Delim('[')
				if out.Dependencies == nil {
					if !in.IsDelim(']') {
						out.Dependencies = make([]*PipelineDependency, 0, 8)
					} else {
						out.Dependencies = []*PipelineDependency{}
					}
				} else {
					out.Dependencies = (out.Dependencies)[:0]
				}
				for !in.IsDelim(']') {
					var v2 *PipelineDependency
					if in.IsNull() {
						in.Skip()
						v2 = nil
					} else {
						if v2 == nil {
							v2 = new(PipelineDependency)
						}
						(*v2).UnmarshalEasyJSON(in)
					}
					out.Dependencies = append(out.Dependencies, v2)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "id":
			out.ID = string(in.String())
		case "status":
			out.Status = string(in.String())
		case "startedMillis":
			out.StartedMillis = int64(in.Int64())
		case "endedMillis":
			out.EndedMillis = int64(in.Int64())
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
func easyjson2189435aEncodeGithubComPolygonIoErrandsServerSchemas1(out *jwriter.Writer, in Pipeline) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"name\":"
		out.RawString(prefix[1:])
		out.String(string(in.Name))
	}
	if in.DeleteOnCompleted {
		const prefix string = ",\"deleteOnCompleted\":"
		out.RawString(prefix)
		out.Bool(bool(in.DeleteOnCompleted))
	}
	if len(in.Errands) != 0 {
		const prefix string = ",\"errands\":"
		out.RawString(prefix)
		{
			out.RawByte('[')
			for v3, v4 := range in.Errands {
				if v3 > 0 {
					out.RawByte(',')
				}
				if v4 == nil {
					out.RawString("null")
				} else {
					(*v4).MarshalEasyJSON(out)
				}
			}
			out.RawByte(']')
		}
	}
	if len(in.Dependencies) != 0 {
		const prefix string = ",\"dependencies\":"
		out.RawString(prefix)
		{
			out.RawByte('[')
			for v5, v6 := range in.Dependencies {
				if v5 > 0 {
					out.RawByte(',')
				}
				if v6 == nil {
					out.RawString("null")
				} else {
					(*v6).MarshalEasyJSON(out)
				}
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix)
		out.String(string(in.ID))
	}
	if in.Status != "" {
		const prefix string = ",\"status\":"
		out.RawString(prefix)
		out.String(string(in.Status))
	}
	{
		const prefix string = ",\"startedMillis\":"
		out.RawString(prefix)
		out.Int64(int64(in.StartedMillis))
	}
	if in.EndedMillis != 0 {
		const prefix string = ",\"endedMillis\":"
		out.RawString(prefix)
		out.Int64(int64(in.EndedMillis))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Pipeline) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson2189435aEncodeGithubComPolygonIoErrandsServerSchemas1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Pipeline) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson2189435aEncodeGithubComPolygonIoErrandsServerSchemas1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Pipeline) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson2189435aDecodeGithubComPolygonIoErrandsServerSchemas1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Pipeline) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson2189435aDecodeGithubComPolygonIoErrandsServerSchemas1(l, v)
}
func easyjson2189435aDecodeGithubComPolygonIoErrandsServerSchemas2(in *jlexer.Lexer, out *Log) {
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
func easyjson2189435aEncodeGithubComPolygonIoErrandsServerSchemas2(out *jwriter.Writer, in Log) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"severity\":"
		out.RawString(prefix[1:])
		out.String(string(in.Severity))
	}
	{
		const prefix string = ",\"message\":"
		out.RawString(prefix)
		out.String(string(in.Message))
	}
	{
		const prefix string = ",\"timestamp\":"
		out.RawString(prefix)
		out.Int64(int64(in.Timestamp))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Log) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson2189435aEncodeGithubComPolygonIoErrandsServerSchemas2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Log) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson2189435aEncodeGithubComPolygonIoErrandsServerSchemas2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Log) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson2189435aDecodeGithubComPolygonIoErrandsServerSchemas2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Log) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson2189435aDecodeGithubComPolygonIoErrandsServerSchemas2(l, v)
}
func easyjson2189435aDecodeGithubComPolygonIoErrandsServerSchemas3(in *jlexer.Lexer, out *Errand) {
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
			easyjson2189435aDecode(in, &out.Options)
		case "data":
			if in.IsNull() {
				in.Skip()
			} else {
				in.Delim('{')
				if !in.IsDelim('}') {
					out.Data = make(map[string]interface{})
				} else {
					out.Data = nil
				}
				for !in.IsDelim('}') {
					key := string(in.String())
					in.WantColon()
					var v7 interface{}
					if m, ok := v7.(easyjson.Unmarshaler); ok {
						m.UnmarshalEasyJSON(in)
					} else if m, ok := v7.(json.Unmarshaler); ok {
						_ = m.UnmarshalJSON(in.Raw())
					} else {
						v7 = in.Interface()
					}
					(out.Data)[key] = v7
					in.WantComma()
				}
				in.Delim('}')
			}
		case "created":
			out.Created = int64(in.Int64())
		case "status":
			out.Status = Status(in.String())
		case "results":
			if in.IsNull() {
				in.Skip()
			} else {
				in.Delim('{')
				if !in.IsDelim('}') {
					out.Results = make(map[string]interface{})
				} else {
					out.Results = nil
				}
				for !in.IsDelim('}') {
					key := string(in.String())
					in.WantColon()
					var v8 interface{}
					if m, ok := v8.(easyjson.Unmarshaler); ok {
						m.UnmarshalEasyJSON(in)
					} else if m, ok := v8.(json.Unmarshaler); ok {
						_ = m.UnmarshalJSON(in.Raw())
					} else {
						v8 = in.Interface()
					}
					(out.Results)[key] = v8
					in.WantComma()
				}
				in.Delim('}')
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
					var v9 Log
					(v9).UnmarshalEasyJSON(in)
					out.Logs = append(out.Logs, v9)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "pipeline":
			out.PipelineID = string(in.String())
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
func easyjson2189435aEncodeGithubComPolygonIoErrandsServerSchemas3(out *jwriter.Writer, in Errand) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.String(string(in.ID))
	}
	{
		const prefix string = ",\"name\":"
		out.RawString(prefix)
		out.String(string(in.Name))
	}
	{
		const prefix string = ",\"type\":"
		out.RawString(prefix)
		out.String(string(in.Type))
	}
	{
		const prefix string = ",\"options\":"
		out.RawString(prefix)
		easyjson2189435aEncode(out, in.Options)
	}
	if len(in.Data) != 0 {
		const prefix string = ",\"data\":"
		out.RawString(prefix)
		{
			out.RawByte('{')
			v10First := true
			for v10Name, v10Value := range in.Data {
				if v10First {
					v10First = false
				} else {
					out.RawByte(',')
				}
				out.String(string(v10Name))
				out.RawByte(':')
				if m, ok := v10Value.(easyjson.Marshaler); ok {
					m.MarshalEasyJSON(out)
				} else if m, ok := v10Value.(json.Marshaler); ok {
					out.Raw(m.MarshalJSON())
				} else {
					out.Raw(json.Marshal(v10Value))
				}
			}
			out.RawByte('}')
		}
	}
	{
		const prefix string = ",\"created\":"
		out.RawString(prefix)
		out.Int64(int64(in.Created))
	}
	if in.Status != "" {
		const prefix string = ",\"status\":"
		out.RawString(prefix)
		out.String(string(in.Status))
	}
	if len(in.Results) != 0 {
		const prefix string = ",\"results\":"
		out.RawString(prefix)
		{
			out.RawByte('{')
			v11First := true
			for v11Name, v11Value := range in.Results {
				if v11First {
					v11First = false
				} else {
					out.RawByte(',')
				}
				out.String(string(v11Name))
				out.RawByte(':')
				if m, ok := v11Value.(easyjson.Marshaler); ok {
					m.MarshalEasyJSON(out)
				} else if m, ok := v11Value.(json.Marshaler); ok {
					out.Raw(m.MarshalJSON())
				} else {
					out.Raw(json.Marshal(v11Value))
				}
			}
			out.RawByte('}')
		}
	}
	{
		const prefix string = ",\"progress\":"
		out.RawString(prefix)
		out.Float64(float64(in.Progress))
	}
	{
		const prefix string = ",\"attempts\":"
		out.RawString(prefix)
		out.Int(int(in.Attempts))
	}
	if in.Started != 0 {
		const prefix string = ",\"started\":"
		out.RawString(prefix)
		out.Int64(int64(in.Started))
	}
	if in.Failed != 0 {
		const prefix string = ",\"failed\":"
		out.RawString(prefix)
		out.Int64(int64(in.Failed))
	}
	if in.Completed != 0 {
		const prefix string = ",\"compelted\":"
		out.RawString(prefix)
		out.Int64(int64(in.Completed))
	}
	if len(in.Logs) != 0 {
		const prefix string = ",\"logs\":"
		out.RawString(prefix)
		{
			out.RawByte('[')
			for v12, v13 := range in.Logs {
				if v12 > 0 {
					out.RawByte(',')
				}
				(v13).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	if in.PipelineID != "" {
		const prefix string = ",\"pipeline\":"
		out.RawString(prefix)
		out.String(string(in.PipelineID))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Errand) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson2189435aEncodeGithubComPolygonIoErrandsServerSchemas3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Errand) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson2189435aEncodeGithubComPolygonIoErrandsServerSchemas3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Errand) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson2189435aDecodeGithubComPolygonIoErrandsServerSchemas3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Errand) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson2189435aDecodeGithubComPolygonIoErrandsServerSchemas3(l, v)
}
func easyjson2189435aDecode(in *jlexer.Lexer, out *struct {
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
func easyjson2189435aEncode(out *jwriter.Writer, in struct {
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
		first = false
		out.RawString(prefix[1:])
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
