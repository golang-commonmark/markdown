// Copyright 2017 James.gx_Chen@htc.com All rights reserved.

package markdown

import (
	"unicode"
	"unicode/utf8"
)

func scanSupScript(s *StateInline, start int) (canOpen bool, canClose bool, count int) {
	pos := start
	max := s.PosMax
	src := s.Src
	marker := src[start]

	lastChar, lastLen := utf8.DecodeLastRuneInString(src[:start])

	for pos < max && src[pos] == marker {
		pos++
	}
	count = pos - start

	nextChar, nextLen := utf8.DecodeRuneInString(src[pos:])

	isLastSpaceOrStart := lastLen == 0 || unicode.IsSpace(lastChar)
	isNextSpaceOrEnd := nextLen == 0 || unicode.IsSpace(nextChar)
	isLastPunct := !isLastSpaceOrStart && (isMarkdownPunct(lastChar) || unicode.IsPunct(lastChar))
	isNextPunct := !isNextSpaceOrEnd && (isMarkdownPunct(nextChar) || unicode.IsPunct(nextChar))

	leftFlanking := !isNextSpaceOrEnd && (!isNextPunct || isLastSpaceOrStart || isLastPunct)
	rightFlanking := !isLastSpaceOrStart && (!isLastPunct || isNextSpaceOrEnd || isNextPunct)

	canOpen = leftFlanking
	canClose = rightFlanking

	return
}

var sup [256]bool

func init() {
	sup['^'] = true
}

func ruleSupScript(s *StateInline, silent bool) (_ bool) {
	src := s.Src
	max := s.PosMax
	start := s.Pos
	marker := src[start]

	if !sup[marker] {
		return
	}

	if silent {
		return
	}

	canOpen, _, startCount := scanSupScript(s, start)
	s.Pos += startCount
	if !canOpen {
		s.Pending.WriteString(src[start:s.Pos])
		return true
	}

	stack := []int{startCount}
	found := false

	for s.Pos < max {
		if src[s.Pos] == marker {
			canOpen, canClose, count := scanSupScript(s, s.Pos)

			if canClose {
				oldCount := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				newCount := count

				for oldCount != newCount {
					if newCount < oldCount {
						stack = append(stack, oldCount-newCount)
						break
					}

					newCount -= oldCount

					if len(stack) == 0 {
						break
					}

					s.Pos += oldCount
					oldCount = stack[len(stack)-1]
					stack = stack[:len(stack)-1]
				}

				if len(stack) == 0 {
					startCount = oldCount
					found = true
					break
				}

				s.Pos += count
				continue
			}

			if canOpen {
				stack = append(stack, count)
			}

			s.Pos += count
			continue
		}

		s.Md.Inline.SkipToken(s)
	}

	if !found {
		s.Pos = start
		return
	}

	s.PosMax = s.Pos
	s.Pos = start + startCount

	count := startCount
	if count > 0 {
		s.PushOpeningToken(&SupScriptOpen{})
	}

	s.Md.Inline.Tokenize(s)

	if count%2 != 0 {
		s.PushClosingToken(&SupScriptClose{})
	}

	s.Pos = s.PosMax + startCount
	s.PosMax = max

	return true
}
