package commands

import (
	"custom-redis/store"
	"errors"
	"math"
	"strconv"
)

// //exists
// 1. invalid: return "ERR value is not an integer"
// 2. valid and within range: ok
// 3. valid but on border :"ERR integer overflow"
// 4. valid but out of range: "ERR value is not an integer or out of range"

// //dont exists
// 1. take number=0

func ProcessIncrement(request CommandRequest) (string, error) {
	entry := store.GetEntry(request.Key)
	if entry.Value == "" {
		entry.Value = "0"
	}

	num, err := strconv.ParseInt(entry.Value, 10, 64)
	if err != nil {
		//outside the signed 64-bit range
		if errors.Is(err, strconv.ErrRange) {
			return "", errors.New("value is not an integer or out of range")
		}
		// Invalid integer syntax.
		return "", errors.New("value is not an integer")
	}

	if num == math.MaxInt64 {
		return "", errors.New("integer overflow")
	}

	entry.Value = strconv.FormatInt(num+1, 10)
	store.SetEntry(request.Key, entry)
	return entry.Value, nil
}
