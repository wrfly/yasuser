package utils

import (
	"crypto/md5"
	"fmt"
)

const CODE62 = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const CODE_LENTH = 62

var EDOC = map[string]int{"0": 0, "1": 1, "2": 2, "3": 3, "4": 4, "5": 5, "6": 6, "7": 7, "8": 8, "9": 9, "a": 10, "b": 11, "c": 12, "d": 13, "e": 14, "f": 15, "g": 16, "h": 17, "i": 18, "j": 19, "k": 20, "l": 21, "m": 22, "n": 23, "o": 24, "p": 25, "q": 26, "r": 27, "s": 28, "t": 29, "u": 30, "v": 31, "w": 32, "x": 33, "y": 34, "z": 35, "A": 36, "B": 37, "C": 38, "D": 39, "E": 40, "F": 41, "G": 42, "H": 43, "I": 44, "J": 45, "K": 46, "L": 47, "M": 48, "N": 49, "O": 50, "P": 51, "Q": 52, "R": 53, "S": 54, "T": 55, "U": 56, "V": 57, "W": 58, "X": 59, "Y": 60, "Z": 61}

func MD5(in string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(in)))
}

// CalHash base62
func CalHash(in uint64) string {
	s := encode(in)
	r := ""
	for i := 0; i < 6-len(s); i++ {
		r += "0"
	}
	return fmt.Sprintf("%s%s", r, s)
}

/**
 * 编码 整数 为 base62 字符串
 */
func encode(number uint64) []byte {
	if number == 0 {
		return []byte("0")
	}
	result := make([]byte, 0)
	for number > 0 {
		round := number / CODE_LENTH
		remain := number % CODE_LENTH
		result = append(result, CODE62[remain])
		number = round
	}
	return result
}
