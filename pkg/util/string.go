package util

import "strconv"

// StringToUint64 字符串转uint64
func StringToUint64(str string) (uint64, error) {
	if str == "" {
		return 0, nil
	}
	valInt, err := strconv.Atoi(str)
	if err != nil {
		return 0, err
	}

	return uint64(valInt), nil
}

// StringToInt 字符串转int
func StringToInt(str string) (int, error) {
	if str == "" {
		return 0, nil
	}
	valInt, err := strconv.Atoi(str)
	if err != nil {
		return 0, err
	}

	return valInt, nil
}
