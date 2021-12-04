package api

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// getRange tries to parse the ranges in http header
// Both from and to are inclusive
func getRange(header string, totalSize int64) (from, to int64, err error) {
	// If the header is not present just return the whole file
	if header == "" {
		return 0, totalSize - 1, nil
	}
	// Remove prefix
	if !strings.HasPrefix(header, "bytes=") {
		return 0, 0, errors.New("invalid range")
	}
	// Try to parse the header
	headerSplit := strings.Split(header[len("bytes="):], "-")
	// Check normal string
	if len(headerSplit) != 2 {
		return 0, 0, errors.New("invalid header split size")
	}
	// Parse from
	suffixMode := false
	if headerSplit[0] == "" {
		suffixMode = true
	} else {
		from, err = strconv.ParseInt(headerSplit[0], 10, 64)
		if err != nil {
			err = fmt.Errorf("cannot parse from: %w", err)
			return
		}
	}
	// Parse to
	if headerSplit[1] == "" {
		to = totalSize - 1
	} else {
		to, err = strconv.ParseInt(headerSplit[1], 10, 64)
		if err != nil {
			err = fmt.Errorf("cannot parse to: %w", err)
			return
		}
		if suffixMode {
			from = totalSize - to
			to = totalSize - 1
		}
	}
	// Check ranges
	if from < 0 || to > totalSize || to < from || to < 0 {
		return 0, 0, errors.New("invalid range")
	}
	return
}
