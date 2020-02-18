package strutils

// Contains returns true is ss slice contains s
//	for substring, use `strings.Contains(s string, substring string)`
func Contains(s string, ss []string) bool {
	for _, str := range ss {
		if str == s {
			return true
		}
	}
	return false
}

// LongestCommonSubstring returns the longest common substring of s1 and s2
//	reference: https://en.wikibooks.org/wiki/Algorithm_Implementation/Strings/Longest_common_substring
//	O(nm), O(nm)
func LongestCommonSubstring(s1 string, s2 string) string {
	if s1 == "" || s2 == "" {
		return ""
	}
	var m = make([][]int, 1+len(s1))
	for i := 0; i < len(m); i++ {
		m[i] = make([]int, 1+len(s2))
	}
	longest := 0
	xLongest := 0
	for x := 1; x < 1+len(s1); x++ {
		for y := 1; y < 1+len(s2); y++ {
			if s1[x-1] == s2[y-1] {
				m[x][y] = m[x-1][y-1] + 1
				if m[x][y] > longest {
					longest = m[x][y]
					xLongest = x
				}
			}
		}
	}
	return s1[xLongest-longest : xLongest]
}
