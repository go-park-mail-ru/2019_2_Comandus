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

func easyjson61e0ab13DecodeGithubComGoParkMailRu20192ComandusInternalModel(in *jlexer.Lexer, out *Review) {
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
		case "companyName":
			out.CompanyName = string(in.String())
		case "jobTitle":
			out.JobTitle = string(in.String())
		case "clientGrade":
			out.ClientGrade = int(in.Int())
		case "clientComment":
			out.ClientComment = string(in.String())
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
func easyjson61e0ab13EncodeGithubComGoParkMailRu20192ComandusInternalModel(out *jwriter.Writer, in Review) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"companyName\":"
		out.RawString(prefix[1:])
		out.String(string(in.CompanyName))
	}
	{
		const prefix string = ",\"jobTitle\":"
		out.RawString(prefix)
		out.String(string(in.JobTitle))
	}
	{
		const prefix string = ",\"clientGrade\":"
		out.RawString(prefix)
		out.Int(int(in.ClientGrade))
	}
	{
		const prefix string = ",\"clientComment\":"
		out.RawString(prefix)
		out.String(string(in.ClientComment))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Review) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson61e0ab13EncodeGithubComGoParkMailRu20192ComandusInternalModel(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Review) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson61e0ab13EncodeGithubComGoParkMailRu20192ComandusInternalModel(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Review) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson61e0ab13DecodeGithubComGoParkMailRu20192ComandusInternalModel(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Review) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson61e0ab13DecodeGithubComGoParkMailRu20192ComandusInternalModel(l, v)
}
func easyjson61e0ab13DecodeGithubComGoParkMailRu20192ComandusInternalModel1(in *jlexer.Lexer, out *OutputResponse) {
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
			out.Id = int64(in.Int64())
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
func easyjson61e0ab13EncodeGithubComGoParkMailRu20192ComandusInternalModel1(out *jwriter.Writer, in OutputResponse) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.Int64(int64(in.Id))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v OutputResponse) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson61e0ab13EncodeGithubComGoParkMailRu20192ComandusInternalModel1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v OutputResponse) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson61e0ab13EncodeGithubComGoParkMailRu20192ComandusInternalModel1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *OutputResponse) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson61e0ab13DecodeGithubComGoParkMailRu20192ComandusInternalModel1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *OutputResponse) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson61e0ab13DecodeGithubComGoParkMailRu20192ComandusInternalModel1(l, v)
}
func easyjson61e0ab13DecodeGithubComGoParkMailRu20192ComandusInternalModel2(in *jlexer.Lexer, out *ExtendResponse) {
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
		case "Response":
			if in.IsNull() {
				in.Skip()
				out.R = nil
			} else {
				if out.R == nil {
					out.R = new(Response)
				}
				easyjson61e0ab13DecodeGithubComGoParkMailRu20192ComandusInternalModel3(in, out.R)
			}
		case "firstName":
			out.FirstName = string(in.String())
		case "lastName":
			out.SecondName = string(in.String())
		case "jobTitle":
			out.JobTitle = string(in.String())
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
func easyjson61e0ab13EncodeGithubComGoParkMailRu20192ComandusInternalModel2(out *jwriter.Writer, in ExtendResponse) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"Response\":"
		out.RawString(prefix[1:])
		if in.R == nil {
			out.RawString("null")
		} else {
			easyjson61e0ab13EncodeGithubComGoParkMailRu20192ComandusInternalModel3(out, *in.R)
		}
	}
	{
		const prefix string = ",\"firstName\":"
		out.RawString(prefix)
		out.String(string(in.FirstName))
	}
	{
		const prefix string = ",\"lastName\":"
		out.RawString(prefix)
		out.String(string(in.SecondName))
	}
	{
		const prefix string = ",\"jobTitle\":"
		out.RawString(prefix)
		out.String(string(in.JobTitle))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ExtendResponse) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson61e0ab13EncodeGithubComGoParkMailRu20192ComandusInternalModel2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ExtendResponse) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson61e0ab13EncodeGithubComGoParkMailRu20192ComandusInternalModel2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ExtendResponse) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson61e0ab13DecodeGithubComGoParkMailRu20192ComandusInternalModel2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ExtendResponse) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson61e0ab13DecodeGithubComGoParkMailRu20192ComandusInternalModel2(l, v)
}
func easyjson61e0ab13DecodeGithubComGoParkMailRu20192ComandusInternalModel3(in *jlexer.Lexer, out *Response) {
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
		case "freelancerId":
			out.FreelancerId = int64(in.Int64())
		case "jobId":
			out.JobId = int64(in.Int64())
		case "files":
			out.Files = string(in.String())
		case "date":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.Date).UnmarshalJSON(data))
			}
		case "statusManager":
			out.StatusManager = string(in.String())
		case "statusFreelancer":
			out.StatusFreelancer = string(in.String())
		case "paymentAmount":
			out.PaymentAmount = float32(in.Float32Str())
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
func easyjson61e0ab13EncodeGithubComGoParkMailRu20192ComandusInternalModel3(out *jwriter.Writer, in Response) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.Int64(int64(in.ID))
	}
	{
		const prefix string = ",\"freelancerId\":"
		out.RawString(prefix)
		out.Int64(int64(in.FreelancerId))
	}
	{
		const prefix string = ",\"jobId\":"
		out.RawString(prefix)
		out.Int64(int64(in.JobId))
	}
	{
		const prefix string = ",\"files\":"
		out.RawString(prefix)
		out.String(string(in.Files))
	}
	{
		const prefix string = ",\"date\":"
		out.RawString(prefix)
		out.Raw((in.Date).MarshalJSON())
	}
	{
		const prefix string = ",\"statusManager\":"
		out.RawString(prefix)
		out.String(string(in.StatusManager))
	}
	{
		const prefix string = ",\"statusFreelancer\":"
		out.RawString(prefix)
		out.String(string(in.StatusFreelancer))
	}
	{
		const prefix string = ",\"paymentAmount\":"
		out.RawString(prefix)
		out.Float32Str(float32(in.PaymentAmount))
	}
	out.RawByte('}')
}
func easyjson61e0ab13DecodeGithubComGoParkMailRu20192ComandusInternalModel4(in *jlexer.Lexer, out *ExtendFreelancer) {
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
		case "freelancer":
			if in.IsNull() {
				in.Skip()
				out.F = nil
			} else {
				if out.F == nil {
					out.F = new(Freelancer)
				}
				(*out.F).UnmarshalEasyJSON(in)
			}
		case "firstName":
			out.FirstName = string(in.String())
		case "secondName":
			out.SecondName = string(in.String())
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
func easyjson61e0ab13EncodeGithubComGoParkMailRu20192ComandusInternalModel4(out *jwriter.Writer, in ExtendFreelancer) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"freelancer\":"
		out.RawString(prefix[1:])
		if in.F == nil {
			out.RawString("null")
		} else {
			(*in.F).MarshalEasyJSON(out)
		}
	}
	{
		const prefix string = ",\"firstName\":"
		out.RawString(prefix)
		out.String(string(in.FirstName))
	}
	{
		const prefix string = ",\"secondName\":"
		out.RawString(prefix)
		out.String(string(in.SecondName))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ExtendFreelancer) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson61e0ab13EncodeGithubComGoParkMailRu20192ComandusInternalModel4(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ExtendFreelancer) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson61e0ab13EncodeGithubComGoParkMailRu20192ComandusInternalModel4(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ExtendFreelancer) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson61e0ab13DecodeGithubComGoParkMailRu20192ComandusInternalModel4(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ExtendFreelancer) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson61e0ab13DecodeGithubComGoParkMailRu20192ComandusInternalModel4(l, v)
}
