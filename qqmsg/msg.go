package qqmsg

import (
	"encoding/hex"
	"fmt"
)

const (
	MsgText         = 1
	MsgFace         = 2
	MsgGroupImage   = 3
	MsgPrivateImage = 6
	MsgVoice        = 7
	MsgNickName     = 18
	MsgVideo        = 26
)

type (
	Msg struct {
		Header         Header
		Elements       []MsgElem
		SenderNickname string
	}

	MsgElem interface {
		Type() ElementType
	}

	ElementType int
)

type Header struct {
	Time       uint32
	Rand       uint32
	Color      uint32
	FontSize   uint8
	FontSylte  uint8
	Charset    uint8
	FontFamily uint8
	FontName   string
}

func Unpack(b []byte) Msg {
	// MSG
	var msg Msg
	var header Header

	buf := NewBuffer(b)
	buf.Skip(8)
	header.Time = buf.Uint32()
	header.Rand = buf.Uint32()
	header.Color = buf.Uint32()
	header.FontSize = buf.Byte()
	header.FontSylte = buf.Byte()
	header.Charset = buf.Byte()
	header.FontFamily = buf.Byte()
	// 24:
	fontName := buf.Read(int(buf.Uint16()))
	header.FontName = DecodeUtf16(fontName)
	msg.Header = header

	buf.Skip(2)

	for !buf.empty() {
		t, _, v := buf.TLV()
		switch t {
		case MsgNickName:
			msg.SenderNickname = DecodeNickname(v)
		default:
			if d, ok := MsgDecoders[t]; ok {
				msg.Elements = append(msg.Elements, d(v))
			}
		}
	}
	return msg
}

var MsgDecoders = map[uint8]func([]byte) MsgElem{
	MsgText:         DecodeTextMsg,
	MsgFace:         DecodeFace,
	MsgGroupImage:   DecodeImage,
	MsgPrivateImage: DecodeImage,
	MsgVoice:        DecodeVoice,
	MsgVideo:        DecodeVideo,
}

func DecodeNickname(b []byte) string {
	buf := NewBuffer(b)
	for !buf.empty() {
		t, _, v := buf.TLV()
		switch t {
		case 1, 2:
			return DecodeUtf16(v)
		}
	}
	return ""
}

func DecodeTextMsg(b []byte) MsgElem {
	buf := NewBuffer(b)
	for !buf.empty() {
		t, _, v := buf.TLV()
		switch t {
		case 1:
			return &TextElement{Content: DecodeUtf16(v)}
		}
	}
	return nil
}

func DecodeFace(b []byte) MsgElem {
	buf := NewBuffer(b)
	for !buf.empty() {
		t, _, v := buf.TLV()
		switch t {
		case 1:
			var id int
			for len(v) > 0 {
				id = int(v[0]) | id<<8
				v = v[1:]
			}
			return &FaceElement{
				Id:   id,
				Name: faceMap[id],
			}
		}
	}
	return nil
}

func DecodeImage(b []byte) MsgElem {
	var elem ImageElement
	buf := NewBuffer(b)
	for !buf.empty() {
		t, _, v := buf.TLV()
		switch t {
		case 1:
			elem.Hash = v
		case 2:
			elem.Path = DecodeUtf16(v)
		}
	}
	return &elem
}

func DecodeVoice(b []byte) MsgElem {
	buf := NewBuffer(b)
	for !buf.empty() {
		t, _, v := buf.TLV()
		switch t {
		case 1:
			return &VoiceElement{Hash: v}
		}
	}
	return nil
}

func DecodeVideo(b []byte) MsgElem {
	buf := NewBuffer(b)
	for !buf.empty() {
		t, _, v := buf.TLV()
		switch t {
		case 1:
			h := v[244 : 244+16]
			for i := range h {
				h[i] ^= 0xEF
			}
			return &VideoElement{Hash: h}
		}
	}
	return nil
}

func EncodeMsg(msg Msg) string {
	encodeElem := func(elems []MsgElem) string {
		ok := ""
		for _, elem := range elems {
			switch e := elem.(type) {
			case *TextElement:
				ok += e.Content
			case *ImageElement:
				ok += fmt.Sprintf("[t:img,path=%s,hash=%s]", e.Path, hex.EncodeToString(e.Hash))
			case *VoiceElement:
				ok += fmt.Sprintf("[t:voice,file=%s,hash=%s.amr]", encodeB48(e.Hash), hex.EncodeToString(e.Hash))
			case *VideoElement:
				ok += fmt.Sprintf("[t:video,hash=%s]", hex.EncodeToString(e.Hash))
			case *FaceElement:
				ok += fmt.Sprintf("[t:face,id=%d]", e.Id)
			}
		}
		return ok
	}
	return encodeElem(msg.Elements)
}
