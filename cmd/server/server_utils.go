package main

import "strconv"

func intToBytes(i int) []byte {
	return []byte(strconv.Itoa(i))
}

func bytesToInt(bs []byte) (int, error) {
	return strconv.Atoi(string(bs))
}
