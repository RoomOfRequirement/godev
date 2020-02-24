package utils

import (
	"log"
	"time"
)

// RoughTiming ...
//	use it to roughly profile func execution time
//	remember: defer func costs time itself
//	use benchmark for accurate time profiling
func RoughTiming(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s costs %s", name, elapsed)
}
