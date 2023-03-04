package utils

import "strconv"

func StringToUint32(s string) (uint32, error) {
	u64, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0, err
	}

	return uint32(u64), nil
}
