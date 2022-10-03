package qqmsg

type TextElement struct {
	Content string
}

type ImageElement struct {
	Path string
	Hash string
}

type FaceElement struct {
	Id   int
	Name string
}

type VoiceElement struct {
	Hash string
}

type VideoElement struct {
	Hash string
}

func (e *TextElement) Type() ElementType {
	return MsgText
}

func (e *ImageElement) Type() ElementType {
	return MsgGroupImage
}

func (e *FaceElement) Type() ElementType {
	return MsgFace
}

func (e *VoiceElement) Type() ElementType {
	return MsgVoice
}

func (e *VideoElement) Type() ElementType {
	return MsgVideo
}
