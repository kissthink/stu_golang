package dev

type ListString []string
func (l ListString) Len() int        { return len(l) }
func (l *ListString) Append(val string) { *l = append(*l, val) }


