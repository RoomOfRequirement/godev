package strutils

import "strings"

// Map ...
func Map(inStrs []string, mapFunc func(string) string) []string {
	for i := 0; i < len(inStrs); i++ {
		inStrs[i] = mapFunc(inStrs[i])
	}
	return inStrs
}

// Reduce ...
//	notice: useless
func Reduce(inStrs []string, reduceFunc func([]string) string) string {
	return reduceFunc(inStrs)
}

// MapReduce ...
func MapReduce(inStrs []string, mapFunc func(string) string, reduceFunc func([]string) string) string {
	for i := 0; i < len(inStrs); i++ {
		inStrs[i] = mapFunc(inStrs[i])
	}
	return reduceFunc(inStrs)
}

// SplitMapReduce ...
func SplitMapReduce(inString, delimiter string, mapFunc func(string) string, reduceFunc func([]string) string) string {
	strs := strings.Split(inString, delimiter)
	for i := 0; i < len(strs); i++ {
		strs[i] = mapFunc(strs[i])
	}
	return reduceFunc(strs)
}
