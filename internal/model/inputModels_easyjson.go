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

func easyjsonA384fcacDecodeGithubComGoParkMailRu20192ComandusInternalModel(in *jlexer.Lexer, out *SearchParams) {
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
		case "minGrade":
			out.MinGrade = int64(in.Int64())
		case "maxGrade":
			out.MaxGrade = int64(in.Int64())
		case "minPaymentAmount":
			out.MinPaymentAmount = float64(in.Float64())
		case "maxPaymentAmount":
			out.MaxPaymentAmount = float64(in.Float64())
		case "country":
			out.Country = int64(in.Int64())
		case "city":
			out.City = int64(in.Int64())
		case "minProposalCount":
			out.MinProposals = int64(in.Int64())
		case "maxProposalCount":
			out.MaxProposals = int64(in.Int64())
		case "experienceLevel":
			if in.IsNull() {
				in.Skip()
			} else {
				in.Delim('[')
				v1 := 0
				for !in.IsDelim(']') {
					if v1 < 3 {
						(out.ExperienceLevel)[v1] = bool(in.Bool())
						v1++
					} else {
						in.SkipRecursive()
					}
					in.WantComma()
				}
				in.Delim(']')
			}
		case "desc":
			out.Desc = bool(in.Bool())
		case "jobTypeId":
			out.JobType = int64(in.Int64())
		case "limit":
			out.Limit = int64(in.Int64())
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
func easyjsonA384fcacEncodeGithubComGoParkMailRu20192ComandusInternalModel(out *jwriter.Writer, in SearchParams) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"minGrade\":"
		out.RawString(prefix[1:])
		out.Int64(int64(in.MinGrade))
	}
	{
		const prefix string = ",\"maxGrade\":"
		out.RawString(prefix)
		out.Int64(int64(in.MaxGrade))
	}
	{
		const prefix string = ",\"minPaymentAmount\":"
		out.RawString(prefix)
		out.Float64(float64(in.MinPaymentAmount))
	}
	{
		const prefix string = ",\"maxPaymentAmount\":"
		out.RawString(prefix)
		out.Float64(float64(in.MaxPaymentAmount))
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
		const prefix string = ",\"minProposalCount\":"
		out.RawString(prefix)
		out.Int64(int64(in.MinProposals))
	}
	{
		const prefix string = ",\"maxProposalCount\":"
		out.RawString(prefix)
		out.Int64(int64(in.MaxProposals))
	}
	{
		const prefix string = ",\"experienceLevel\":"
		out.RawString(prefix)
		out.RawByte('[')
		for v2 := range in.ExperienceLevel {
			if v2 > 0 {
				out.RawByte(',')
			}
			out.Bool(bool((in.ExperienceLevel)[v2]))
		}
		out.RawByte(']')
	}
	{
		const prefix string = ",\"desc\":"
		out.RawString(prefix)
		out.Bool(bool(in.Desc))
	}
	{
		const prefix string = ",\"jobTypeId\":"
		out.RawString(prefix)
		out.Int64(int64(in.JobType))
	}
	{
		const prefix string = ",\"limit\":"
		out.RawString(prefix)
		out.Int64(int64(in.Limit))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v SearchParams) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonA384fcacEncodeGithubComGoParkMailRu20192ComandusInternalModel(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v SearchParams) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonA384fcacEncodeGithubComGoParkMailRu20192ComandusInternalModel(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *SearchParams) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonA384fcacDecodeGithubComGoParkMailRu20192ComandusInternalModel(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *SearchParams) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonA384fcacDecodeGithubComGoParkMailRu20192ComandusInternalModel(l, v)
}
func easyjsonA384fcacDecodeGithubComGoParkMailRu20192ComandusInternalModel1(in *jlexer.Lexer, out *ReviewInput) {
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
		case "grade":
			out.Grade = int(in.Int())
		case "comment":
			out.Comment = string(in.String())
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
func easyjsonA384fcacEncodeGithubComGoParkMailRu20192ComandusInternalModel1(out *jwriter.Writer, in ReviewInput) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"grade\":"
		out.RawString(prefix[1:])
		out.Int(int(in.Grade))
	}
	{
		const prefix string = ",\"comment\":"
		out.RawString(prefix)
		out.String(string(in.Comment))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ReviewInput) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonA384fcacEncodeGithubComGoParkMailRu20192ComandusInternalModel1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ReviewInput) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonA384fcacEncodeGithubComGoParkMailRu20192ComandusInternalModel1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ReviewInput) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonA384fcacDecodeGithubComGoParkMailRu20192ComandusInternalModel1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ReviewInput) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonA384fcacDecodeGithubComGoParkMailRu20192ComandusInternalModel1(l, v)
}
func easyjsonA384fcacDecodeGithubComGoParkMailRu20192ComandusInternalModel2(in *jlexer.Lexer, out *Notification) {
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
		case "new_messages":
			out.NewMessages = bool(in.Bool())
		case "new_projects":
			out.NewProjects = bool(in.Bool())
		case "news_service":
			out.NewsFromService = bool(in.Bool())
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
func easyjsonA384fcacEncodeGithubComGoParkMailRu20192ComandusInternalModel2(out *jwriter.Writer, in Notification) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"new_messages\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Bool(bool(in.NewMessages))
	}
	{
		const prefix string = ",\"new_projects\":"
		out.RawString(prefix)
		out.Bool(bool(in.NewProjects))
	}
	{
		const prefix string = ",\"news_service\":"
		out.RawString(prefix)
		out.Bool(bool(in.NewsFromService))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Notification) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonA384fcacEncodeGithubComGoParkMailRu20192ComandusInternalModel2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Notification) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonA384fcacEncodeGithubComGoParkMailRu20192ComandusInternalModel2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Notification) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonA384fcacDecodeGithubComGoParkMailRu20192ComandusInternalModel2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Notification) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonA384fcacDecodeGithubComGoParkMailRu20192ComandusInternalModel2(l, v)
}
func easyjsonA384fcacDecodeGithubComGoParkMailRu20192ComandusInternalModel3(in *jlexer.Lexer, out *InnerInfo) {
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
		case "user_id":
			out.UserID = int64(in.Int64())
		case "who_see_profile":
			out.WhoSeeProfile = string(in.String())
		case "control_question":
			out.ControlQuestion = string(in.String())
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
func easyjsonA384fcacEncodeGithubComGoParkMailRu20192ComandusInternalModel3(out *jwriter.Writer, in InnerInfo) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"user_id\":"
		out.RawString(prefix[1:])
		out.Int64(int64(in.UserID))
	}
	{
		const prefix string = ",\"who_see_profile\":"
		out.RawString(prefix)
		out.String(string(in.WhoSeeProfile))
	}
	{
		const prefix string = ",\"control_question\":"
		out.RawString(prefix)
		out.String(string(in.ControlQuestion))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v InnerInfo) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonA384fcacEncodeGithubComGoParkMailRu20192ComandusInternalModel3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v InnerInfo) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonA384fcacEncodeGithubComGoParkMailRu20192ComandusInternalModel3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *InnerInfo) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonA384fcacDecodeGithubComGoParkMailRu20192ComandusInternalModel3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *InnerInfo) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonA384fcacDecodeGithubComGoParkMailRu20192ComandusInternalModel3(l, v)
}
func easyjsonA384fcacDecodeGithubComGoParkMailRu20192ComandusInternalModel4(in *jlexer.Lexer, out *ContractInput) {
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
		case "paymentAmount":
			out.PaymentAmount = float32(in.Float32())
		case "timeEstimation":
			out.TimeEstimation = int(in.Int())
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
func easyjsonA384fcacEncodeGithubComGoParkMailRu20192ComandusInternalModel4(out *jwriter.Writer, in ContractInput) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"paymentAmount\":"
		out.RawString(prefix[1:])
		out.Float32(float32(in.PaymentAmount))
	}
	{
		const prefix string = ",\"timeEstimation\":"
		out.RawString(prefix)
		out.Int(int(in.TimeEstimation))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ContractInput) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonA384fcacEncodeGithubComGoParkMailRu20192ComandusInternalModel4(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ContractInput) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonA384fcacEncodeGithubComGoParkMailRu20192ComandusInternalModel4(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ContractInput) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonA384fcacDecodeGithubComGoParkMailRu20192ComandusInternalModel4(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ContractInput) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonA384fcacDecodeGithubComGoParkMailRu20192ComandusInternalModel4(l, v)
}
func easyjsonA384fcacDecodeGithubComGoParkMailRu20192ComandusInternalModel5(in *jlexer.Lexer, out *BodyPassword) {
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
		case "password":
			out.Password = string(in.String())
		case "newPassword":
			out.NewPassword = string(in.String())
		case "newPasswordConfirmation":
			out.NewPasswordConfirmation = string(in.String())
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
func easyjsonA384fcacEncodeGithubComGoParkMailRu20192ComandusInternalModel5(out *jwriter.Writer, in BodyPassword) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"password\":"
		out.RawString(prefix[1:])
		out.String(string(in.Password))
	}
	{
		const prefix string = ",\"newPassword\":"
		out.RawString(prefix)
		out.String(string(in.NewPassword))
	}
	{
		const prefix string = ",\"newPasswordConfirmation\":"
		out.RawString(prefix)
		out.String(string(in.NewPasswordConfirmation))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v BodyPassword) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonA384fcacEncodeGithubComGoParkMailRu20192ComandusInternalModel5(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v BodyPassword) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonA384fcacEncodeGithubComGoParkMailRu20192ComandusInternalModel5(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *BodyPassword) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonA384fcacDecodeGithubComGoParkMailRu20192ComandusInternalModel5(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *BodyPassword) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonA384fcacDecodeGithubComGoParkMailRu20192ComandusInternalModel5(l, v)
}
