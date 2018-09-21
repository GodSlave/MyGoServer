package base

import (
	"crypto/md5"
	"github.com/GodSlave/MyGoServer/log"
)

func GetMd5(allBytes []byte) []byte {

	allBytes[32], allBytes[6] = allBytes[6], allBytes[32]
	allBytes[48], allBytes[31] = allBytes[31], allBytes[48]
	log.Info("AesKey  before is %v", allBytes)
	md5result := md5.Sum(allBytes)
	return md5result[:]
}

func GetMd5T(bytes1 []byte, bytes2 []byte) []byte {
	allBytes := make([]byte, len(bytes1)+len(bytes2))
	copy(allBytes[0:len(bytes1)], bytes1)
	copy(allBytes[len(bytes1):], bytes2)
	return GetMd5(allBytes)
}
