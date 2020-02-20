package strutils

import "errors"

// DistanceType ...
type DistanceType int

const (
	// Levenshtein ...
	Levenshtein DistanceType = iota
	// LongestCommonSubsequence ...
	LongestCommonSubsequence
	// Hamming ...
	Hamming
	// TODO: Damerau-Levenshtein distance and Jaro distance
	// DamerauLevenshtein
	// Jaro
)

// EditDistance ...
/*
 * https://en.wikipedia.org/wiki/Edit_distance
 * The Levenshtein distance allows deletion, insertion and substitution.
 * The Longest common subsequence (LCS) distance allows only insertion and deletion, not substitution.
 * The Hamming distance allows only substitution, hence, it only applies to strings of the same length.
 * The Damerau–Levenshtein distance allows insertion, deletion, substitution, and the transposition of two adjacent characters.
 * The Jaro distance allows only transposition.
 */
func EditDistance(t DistanceType, s1, s2 string) (int, error) {
	switch t {
	case Levenshtein:
		return levenshtein(s1, s2), nil
	case LongestCommonSubsequence:
		return lcs(s1, s2), nil
	case Hamming:
		return hamming(s1, s2)
	default:
		return -1, errors.New("unsupported distance type")
	}
}

// https://en.wikibooks.org/wiki/Algorithm_Implementation/Strings/Levenshtein_distance
func levenshtein(s1, s2 string) int {
	s1len, s2len := len(s1), len(s2)
	if s1len == 0 || s2len == 0 {
		return s1len + s2len
	}
	column := make([]int, s1len+1)
	var x, y, lastDiag, oldDiag, inc int
	for y = 1; y < s1len+1; y++ {
		column[y] = y
	}
	for x = 1; x < s2len+1; x++ {
		column[0] = x
		lastDiag = x - 1
		for y = 1; y < s1len+1; y++ {
			if s1[y-1] == s2[x-1] {
				inc = 0
			} else {
				inc = 1
			}
			oldDiag = column[y]
			column[y] = min3(column[y]+1, column[y-1]+1, lastDiag+inc)
			lastDiag = oldDiag
		}
	}
	return column[s1len]
}

func min3(x, y, z int) int {
	if x < y {
		if x < z {
			return x
		}
	} else {
		if y < z {
			return y
		}
	}
	return z
}

// https://en.wikibooks.org/wiki/Algorithm_Implementation/Strings/Longest_common_substring
func lcs(s1 string, s2 string) int {
	s1len, s2len := len(s1), len(s2)
	if s1len == 0 || s2len == 0 {
		return s1len + s2len
	}
	var m = make([][]int, s1len+1)
	for i := 0; i < len(m); i++ {
		m[i] = make([]int, s2len+1)
	}
	longest := 0
	for x := 1; x < s1len+1; x++ {
		for y := 1; y < s2len+1; y++ {
			if s1[x-1] == s2[y-1] {
				m[x][y] = m[x-1][y-1] + 1
				if m[x][y] > longest {
					longest = m[x][y]
				}
			}
		}
	}
	return s2len - 1 - longest
}

// only for equal length strings
func hamming(s1, s2 string) (int, error) {
	s1len, s2len := len(s1), len(s2)
	if s1len != s2len {
		return -1, errors.New("undefined for sequences of unequal length")
	}
	cnt := 0
	for i := 0; i < s1len; i++ {
		if s1[i] != s2[i] {
			cnt++
		}
	}
	return cnt, nil
}
