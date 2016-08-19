package main

type StatusLine struct {
	Lines []StatusSegment
}

func (l *StatusLine) Prepend(behind *StatusLine) {
	l.Lines = append(behind.Lines, l.Lines...)
}

func (l *StatusLine) Add() *StatusSegment {
	var segment StatusSegment = make(StatusSegment)
	l.Lines = append(l.Lines, segment)
	return &segment
}

func (l *StatusLine) AddLabel(label string) {
	l.Add().
		SetFullText(label).
		SetColor(COLOR_LABEL).
		SetSeparator(false).
		SetSeparatorWidth(0).
		SetMinWidthString("9h59m59s")
}

type StatusSegment map[string]interface{}

func (l *StatusSegment) SetFullText(text string) *StatusSegment {
	(*l)["full_text"] = text
	return l
}

func (l *StatusSegment) SetColor(color string) *StatusSegment {
	(*l)["color"] = color
	return l
}

func (l *StatusSegment) SetSeparator(s bool) *StatusSegment {
	(*l)["separator"] = s
	return l
}

func (l *StatusSegment) SetSeparatorWidth(w int) *StatusSegment {
	(*l)["separator_block_width"] = w
	return l
}

func (l *StatusSegment) SetMinWidthString(s string) *StatusSegment {
	(*l)["min_width"] = s
	return l
}
