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

func easyjson9e1087fdDecodeGithubComGoParkMailRu20192ComandusInternalModel(in *jlexer.Lexer, out *User) {
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
		case "firstName":
			out.FirstName = string(in.String())
		case "secondName":
			out.SecondName = string(in.String())
		case "username":
			out.UserName = string(in.String())
		case "email":
			out.Email = string(in.String())
		case "password":
			out.Password = string(in.String())
		case "type":
			out.UserType = string(in.String())
		case "registrationDate":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.RegistrationDate).UnmarshalJSON(data))
			}
		case "freelancerId":
			out.FreelancerId = int64(in.Int64())
		case "hireManagerId":
			out.HireManagerId = int64(in.Int64())
		case "companyId":
			out.CompanyId = int64(in.Int64())
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
func easyjson9e1087fdEncodeGithubComGoParkMailRu20192ComandusInternalModel(out *jwriter.Writer, in User) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"firstName\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.FirstName))
	}
	{
		const prefix string = ",\"secondName\":"
		out.RawString(prefix)
		out.String(string(in.SecondName))
	}
	{
		const prefix string = ",\"username\":"
		out.RawString(prefix)
		out.String(string(in.UserName))
	}
	{
		const prefix string = ",\"email\":"
		out.RawString(prefix)
		out.String(string(in.Email))
	}
	{
		const prefix string = ",\"password\":"
		out.RawString(prefix)
		out.String(string(in.Password))
	}
	{
		const prefix string = ",\"type\":"
		out.RawString(prefix)
		out.String(string(in.UserType))
	}
	{
		const prefix string = ",\"registrationDate\":"
		out.RawString(prefix)
		out.Raw((in.RegistrationDate).MarshalJSON())
	}
	{
		const prefix string = ",\"freelancerId\":"
		out.RawString(prefix)
		out.Int64(int64(in.FreelancerId))
	}
	{
		const prefix string = ",\"hireManagerId\":"
		out.RawString(prefix)
		out.Int64(int64(in.HireManagerId))
	}
	{
		const prefix string = ",\"companyId\":"
		out.RawString(prefix)
		out.Int64(int64(in.CompanyId))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v User) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson9e1087fdEncodeGithubComGoParkMailRu20192ComandusInternalModel(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v User) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson9e1087fdEncodeGithubComGoParkMailRu20192ComandusInternalModel(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *User) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson9e1087fdDecodeGithubComGoParkMailRu20192ComandusInternalModel(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *User) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson9e1087fdDecodeGithubComGoParkMailRu20192ComandusInternalModel(l, v)
}
