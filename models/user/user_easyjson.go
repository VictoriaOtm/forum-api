// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package user

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

func easyjson9e1087fdDecodeGithubComVictoriaOtmForumApiModelsUser(in *jlexer.Lexer, out *Arr) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(Arr, 0, 1)
			} else {
				*out = Arr{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v1 User
			if data := in.Raw(); in.Ok() {
				in.AddError((v1).UnmarshalJSON(data))
			}
			*out = append(*out, v1)
			in.WantComma()
		}
		in.Delim(']')
	}
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson9e1087fdEncodeGithubComVictoriaOtmForumApiModelsUser(out *jwriter.Writer, in Arr) {
	if in == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
		out.RawString("null")
	} else {
		out.RawByte('[')
		for v2, v3 := range in {
			if v2 > 0 {
				out.RawByte(',')
			}
			out.Raw((v3).MarshalJSON())
		}
		out.RawByte(']')
	}
}

// MarshalJSON supports json.Marshaler interface
func (v Arr) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson9e1087fdEncodeGithubComVictoriaOtmForumApiModelsUser(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Arr) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson9e1087fdEncodeGithubComVictoriaOtmForumApiModelsUser(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Arr) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson9e1087fdDecodeGithubComVictoriaOtmForumApiModelsUser(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Arr) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson9e1087fdDecodeGithubComVictoriaOtmForumApiModelsUser(l, v)
}
func easyjson9e1087fdDecodeGithubComVictoriaOtmForumApiModelsUser1(in *jlexer.Lexer, out *Update) {
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
		case "about":
			if in.IsNull() {
				in.Skip()
				out.About = nil
			} else {
				if out.About == nil {
					out.About = new(string)
				}
				*out.About = string(in.String())
			}
		case "email":
			if in.IsNull() {
				in.Skip()
				out.Email = nil
			} else {
				if out.Email == nil {
					out.Email = new(string)
				}
				*out.Email = string(in.String())
			}
		case "fullname":
			if in.IsNull() {
				in.Skip()
				out.Fullname = nil
			} else {
				if out.Fullname == nil {
					out.Fullname = new(string)
				}
				*out.Fullname = string(in.String())
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
func easyjson9e1087fdEncodeGithubComVictoriaOtmForumApiModelsUser1(out *jwriter.Writer, in Update) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"about\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if in.About == nil {
			out.RawString("null")
		} else {
			out.String(string(*in.About))
		}
	}
	{
		const prefix string = ",\"email\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if in.Email == nil {
			out.RawString("null")
		} else {
			out.String(string(*in.Email))
		}
	}
	{
		const prefix string = ",\"fullname\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if in.Fullname == nil {
			out.RawString("null")
		} else {
			out.String(string(*in.Fullname))
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Update) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson9e1087fdEncodeGithubComVictoriaOtmForumApiModelsUser1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Update) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson9e1087fdEncodeGithubComVictoriaOtmForumApiModelsUser1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Update) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson9e1087fdDecodeGithubComVictoriaOtmForumApiModelsUser1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Update) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson9e1087fdDecodeGithubComVictoriaOtmForumApiModelsUser1(l, v)
}
func easyjson9e1087fdDecodeGithubComVictoriaOtmForumApiModelsUser2(in *jlexer.Lexer, out *User) {
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
		case "about":
			out.About = string(in.String())
		case "email":
			out.Email = string(in.String())
		case "fullname":
			out.Fullname = string(in.String())
		case "nickname":
			out.Nickname = string(in.String())
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
func easyjson9e1087fdEncodeGithubComVictoriaOtmForumApiModelsUser2(out *jwriter.Writer, in User) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"about\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.About))
	}
	{
		const prefix string = ",\"email\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Email))
	}
	{
		const prefix string = ",\"fullname\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Fullname))
	}
	{
		const prefix string = ",\"nickname\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Nickname))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v User) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson9e1087fdEncodeGithubComVictoriaOtmForumApiModelsUser2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v User) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson9e1087fdEncodeGithubComVictoriaOtmForumApiModelsUser2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *User) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson9e1087fdDecodeGithubComVictoriaOtmForumApiModelsUser2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *User) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson9e1087fdDecodeGithubComVictoriaOtmForumApiModelsUser2(l, v)
}
