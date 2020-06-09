package main

import (
	"bytes"
	"strconv"
)

type list struct {
	value int
	next  *list
}

func (l *list) add(v int) *list {
	if l == nil {
		return &list{v, nil}
	}

	head := l
	for l != nil {
		if l.next == nil {
			l.next = &list{v, nil}
			break
		}

		l = l.next
	}

	return head
}

func (l *list) insert(v int) *list {
	if l == nil {
		return &list{v, nil}
	}
	curr := &list{v, l.next}
	l.next = curr

	return curr
}

func (l *list) slice() []int {
	var s []int
	for l != nil {
		s = append(s, l.value)
		l = l.next
	}

	return s
}

func (l *list) String() string {
	s := l.slice()
	var b bytes.Buffer
	for _, v := range s {
		b.WriteString(strconv.Itoa(v))
		b.WriteString(" ")
	}

	return b.String()
}

func mergeSorted(l1 *list, l2 *list) *list {
	merged := &list{}
	start := merged
	for l1 != nil && l2 != nil {
		if l1.value < l2.value {
			merged.next = l1
			l1 = l1.next
		} else {
			merged.next = l2
			l2 = l2.next
		}
		merged = merged.next
	}

	if l1 != nil {
		merged.next = l1
	}

	if l2 != nil {
		merged.next = l2
	}

	return start.next
}

func reverseSublist(l *list, begin int, end int) {
}

func cycle(list list) bool {
	return false
}
