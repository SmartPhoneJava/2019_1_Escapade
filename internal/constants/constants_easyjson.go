// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package constants

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

func easyjson8a25db09DecodeGithubComGoParkMailRu20191EscapadeInternalConstants(in *jlexer.Lexer, out *roomConfiguration) {
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
		case "Set":
			out.Set = bool(in.Bool())
		case "nameMinLength":
			out.NameMin = int(in.Int())
		case "nameMaxLength":
			out.NameMax = int(in.Int())
		case "timeToPrepareMin":
			out.TimeToPrepareMin = int(in.Int())
		case "timeToPrepareMax":
			out.TimeToPrepareMax = int(in.Int())
		case "timeToPlayMin":
			out.TimeToPlayMin = int(in.Int())
		case "timeToPlayMax":
			out.TimeToPlayMax = int(in.Int())
		case "playersMin":
			out.PlayersMin = int(in.Int())
		case "playersMax":
			out.PlayersMax = int(in.Int())
		case "observersMax":
			out.ObserversMax = int(in.Int())
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
func easyjson8a25db09EncodeGithubComGoParkMailRu20191EscapadeInternalConstants(out *jwriter.Writer, in roomConfiguration) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"Set\":"
		out.RawString(prefix[1:])
		out.Bool(bool(in.Set))
	}
	{
		const prefix string = ",\"nameMinLength\":"
		out.RawString(prefix)
		out.Int(int(in.NameMin))
	}
	{
		const prefix string = ",\"nameMaxLength\":"
		out.RawString(prefix)
		out.Int(int(in.NameMax))
	}
	{
		const prefix string = ",\"timeToPrepareMin\":"
		out.RawString(prefix)
		out.Int(int(in.TimeToPrepareMin))
	}
	{
		const prefix string = ",\"timeToPrepareMax\":"
		out.RawString(prefix)
		out.Int(int(in.TimeToPrepareMax))
	}
	{
		const prefix string = ",\"timeToPlayMin\":"
		out.RawString(prefix)
		out.Int(int(in.TimeToPlayMin))
	}
	{
		const prefix string = ",\"timeToPlayMax\":"
		out.RawString(prefix)
		out.Int(int(in.TimeToPlayMax))
	}
	{
		const prefix string = ",\"playersMin\":"
		out.RawString(prefix)
		out.Int(int(in.PlayersMin))
	}
	{
		const prefix string = ",\"playersMax\":"
		out.RawString(prefix)
		out.Int(int(in.PlayersMax))
	}
	{
		const prefix string = ",\"observersMax\":"
		out.RawString(prefix)
		out.Int(int(in.ObserversMax))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v roomConfiguration) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson8a25db09EncodeGithubComGoParkMailRu20191EscapadeInternalConstants(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v roomConfiguration) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson8a25db09EncodeGithubComGoParkMailRu20191EscapadeInternalConstants(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *roomConfiguration) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson8a25db09DecodeGithubComGoParkMailRu20191EscapadeInternalConstants(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *roomConfiguration) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson8a25db09DecodeGithubComGoParkMailRu20191EscapadeInternalConstants(l, v)
}
func easyjson8a25db09DecodeGithubComGoParkMailRu20191EscapadeInternalConstants1(in *jlexer.Lexer, out *fieldConfiguration) {
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
		case "Set":
			out.Set = bool(in.Bool())
		case "widthMin":
			out.WidthMin = int(in.Int())
		case "widthMax":
			out.WidthMax = int(in.Int())
		case "heightMin":
			out.HeightMin = int(in.Int())
		case "heightMax":
			out.HeightMax = int(in.Int())
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
func easyjson8a25db09EncodeGithubComGoParkMailRu20191EscapadeInternalConstants1(out *jwriter.Writer, in fieldConfiguration) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"Set\":"
		out.RawString(prefix[1:])
		out.Bool(bool(in.Set))
	}
	{
		const prefix string = ",\"widthMin\":"
		out.RawString(prefix)
		out.Int(int(in.WidthMin))
	}
	{
		const prefix string = ",\"widthMax\":"
		out.RawString(prefix)
		out.Int(int(in.WidthMax))
	}
	{
		const prefix string = ",\"heightMin\":"
		out.RawString(prefix)
		out.Int(int(in.HeightMin))
	}
	{
		const prefix string = ",\"heightMax\":"
		out.RawString(prefix)
		out.Int(int(in.HeightMax))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v fieldConfiguration) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson8a25db09EncodeGithubComGoParkMailRu20191EscapadeInternalConstants1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v fieldConfiguration) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson8a25db09EncodeGithubComGoParkMailRu20191EscapadeInternalConstants1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *fieldConfiguration) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson8a25db09DecodeGithubComGoParkMailRu20191EscapadeInternalConstants1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *fieldConfiguration) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson8a25db09DecodeGithubComGoParkMailRu20191EscapadeInternalConstants1(l, v)
}
