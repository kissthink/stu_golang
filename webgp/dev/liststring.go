package dev


type ListString []string

func (l ListString) Len() int           { return len(l) }
func (l *ListString) Append(val string) { *l = append(*l, val) }

type LStr []string

func (l LStr) Len() int { return len(l) }
func (l *LStr) Append(val string) {
	*l = append(*l, val)
}
