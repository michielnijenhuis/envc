package cmd

import "strings"

func matchPattern(s string, pattern string, si int, pi int, caseSensitive bool) bool {
	if (s == "" && pattern == "") || (si >= len(s) && pi >= len(pattern)) {
		return true
	}

	if !caseSensitive {
		s = strings.ToLower(s)
		pattern = strings.ToLower(pattern)
	}

	if pi < len(pattern) && pattern[pi] == '*' {
		if matchPattern(s, pattern, si, pi+1, true) {
			return true
		}

		return si < len(s) && matchPattern(s, pattern, si+1, pi, true)
	}

	if (pi < len(pattern) && pattern[pi] == '?') || (si < len(s) && pi < len(pattern) && s[si] == pattern[pi]) {
		return matchPattern(s, pattern, si+1, pi+1, true)
	}

	return false
}

func applySkipPattern(index map[string]int, pattern string) map[string]int {
	skip := strings.Split(pattern, ",")

	for _, pattern := range skip {
		pattern = strings.TrimSpace(pattern)
		if pattern == "" {
			continue
		}

		for k := range index {
			if matchPattern(k, pattern, 0, 0, false) {
				delete(index, k)
			}
		}
	}

	return index
}

func applyTargetPattern(index map[string]int, p string) map[string]int {
	pattern := strings.Split(p, ",")
	newKeyIndex := make(map[string]int, len(index))

	for _, p := range pattern {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}

		for k := range index {
			if matchPattern(k, p, 0, 0, false) {
				newKeyIndex[k] = 0
			}
		}
	}

	return newKeyIndex
}
