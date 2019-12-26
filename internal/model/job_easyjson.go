// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package model

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

func easyjson8a33d6c7DecodeGithubComGoParkMailRu20192ComandusInternalModel(in *jlexer.Lexer, out *Job) {
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
			out.ID = int64(in.Int64())
		case "hireManagerId":
			out.HireManagerId = int64(in.Int64Str())
		case "title":
			out.Title = string(in.String())
		case "description":
			out.Description = string(in.String())
		case "files":
			out.Files = string(in.String())
		case "specialityId":
			out.SpecialityId = int64(in.Int64Str())
		case "experienceLevelId":
			out.ExperienceLevelId = int64(in.Int64Str())
		case "paymentAmount":
			out.PaymentAmount = float32(in.Float32Str())
		case "country":
			out.Country = int64(in.Int64())
		case "city":
			out.City = int64(in.Int64())
		case "jobTypeId":
			out.JobTypeId = int64(in.Int64Str())
		case "date":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.Date).UnmarshalJSON(data))
			}
		case "status":
			out.Status = string(in.String())
		case "tagLine":
			out.TagLine = string(in.String())
		case "proposals":
			out.Proposals = int64(in.Int64())
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
func easyjson8a33d6c7EncodeGithubComGoParkMailRu20192ComandusInternalModel(out *jwriter.Writer, in Job) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.Int64(int64(in.ID))
	}
	{
		const prefix string = ",\"hireManagerId\":"
		out.RawString(prefix)
		out.Int64Str(int64(in.HireManagerId))
	}
	{
		const prefix string = ",\"title\":"
		out.RawString(prefix)
		out.String(string(in.Title))
	}
	{
		const prefix string = ",\"description\":"
		out.RawString(prefix)
		out.String(string(in.Description))
	}
	{
		const prefix string = ",\"files\":"
		out.RawString(prefix)
		out.String(string(in.Files))
	}
	{
		const prefix string = ",\"specialityId\":"
		out.RawString(prefix)
		out.Int64Str(int64(in.SpecialityId))
	}
	{
		const prefix string = ",\"experienceLevelId\":"
		out.RawString(prefix)
		out.Int64Str(int64(in.ExperienceLevelId))
	}
	{
		const prefix string = ",\"paymentAmount\":"
		out.RawString(prefix)
		out.Float32Str(float32(in.PaymentAmount))
	}
	{
		const prefix string = ",\"country\":"
		out.RawString(prefix)
		out.Int64(int64(in.Country))
	}
	{
		const prefix string = ",\"city\":"
		out.RawString(prefix)
		out.Int64(int64(in.City))
	}
	{
		const prefix string = ",\"jobTypeId\":"
		out.RawString(prefix)
		out.Int64Str(int64(in.JobTypeId))
	}
	{
		const prefix string = ",\"date\":"
		out.RawString(prefix)
		out.Raw((in.Date).MarshalJSON())
	}
	{
		const prefix string = ",\"status\":"
		out.RawString(prefix)
		out.String(string(in.Status))
	}
	{
		const prefix string = ",\"tagLine\":"
		out.RawString(prefix)
		out.String(string(in.TagLine))
	}
	{
		const prefix string = ",\"proposals\":"
		out.RawString(prefix)
		out.Int64(int64(in.Proposals))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Job) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson8a33d6c7EncodeGithubComGoParkMailRu20192ComandusInternalModel(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Job) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson8a33d6c7EncodeGithubComGoParkMailRu20192ComandusInternalModel(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Job) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson8a33d6c7DecodeGithubComGoParkMailRu20192ComandusInternalModel(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Job) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson8a33d6c7DecodeGithubComGoParkMailRu20192ComandusInternalModel(l, v)
}
