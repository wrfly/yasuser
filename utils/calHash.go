package utils

import (
	"crypto/md5"
	"fmt"
	"math"
	"strings"
)

// thanks to http://www.01happy.com/golang-base62-encode/

const code62 = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const codeLenth = 62

var edoc = map[string]int{"0": 0, "1": 1, "2": 2, "3": 3, "4": 4, "5": 5, "6": 6, "7": 7, "8": 8, "9": 9, "a": 10, "b": 11, "c": 12, "d": 13, "e": 14, "f": 15, "g": 16, "h": 17, "i": 18, "j": 19, "k": 20, "l": 21, "m": 22, "n": 23, "o": 24, "p": 25, "q": 26, "r": 27, "s": 28, "t": 29, "u": 30, "v": 31, "w": 32, "x": 33, "y": 34, "z": 35, "A": 36, "B": 37, "C": 38, "D": 39, "E": 40, "F": 41, "G": 42, "H": 43, "I": 44, "J": 45, "K": 46, "L": 47, "M": 48, "N": 49, "O": 50, "P": 51, "Q": 52, "R": 53, "S": 54, "T": 55, "U": 56, "V": 57, "W": 58, "X": 59, "Y": 60, "Z": 61}

func MD5(in string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(in)))
}

// CalHash base62, max is 56800235583
func CalHash(in int64) string {
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
func encode(number int64) []byte {
	if number == 0 {
		return []byte("0")
	}
	result := make([]byte, 0)
	for number > 0 {
		round := number / codeLenth
		remain := number % codeLenth
		result = append(result, code62[remain])
		number = round
	}
	l := len(result)
	for i := 0; i < l-i; i++ {
		result[i], result[l-i-1] = result[l-i-1], result[i]
	}
	return result
}

/**
 * 解码字符串为整数
 */
func decode(str string) int {
	str = strings.TrimSpace(str)
	result := 0
	for index, char := range []byte(str) {
		result += edoc[string(char)] * int(math.Pow(codeLenth, float64(index)))
	}
	return result
}
