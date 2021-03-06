// Copyright 2014 mqantserver Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

func NewAesEncrypt(key []byte) (aes *AesEncrypt, err error) {
	keyLen := len(key)
	if keyLen < 16 {
		err = fmt.Errorf("The length of res key shall not be less than 16")
		return
	}
	aes = &AesEncrypt{
		Key: key,
	}
	return aes, nil
}

type AesEncrypt struct {
	Key []byte
}

func (this *AesEncrypt) getKey() []byte {
	keyLen := len(this.Key)
	if keyLen < 16 {
		panic("The length of res key shall not be less than 16")
	}
	arrKey := []byte(this.Key)
	if keyLen >= 32 {
		//取前32个字节
		return arrKey[:32]
	}
	if keyLen >= 24 {
		//取前24个字节
		return arrKey[:24]
	}
	//取前16个字节
	return arrKey[:16]
}

//加密字符串
func (this *AesEncrypt) Encrypt(strMesg string) ([]byte, error) {
	return this.EncryptBytes([]byte(strMesg))
}

//加密字符串
func (this *AesEncrypt) EncryptBytes(data []byte) ([]byte, error) {
	key := this.getKey()
	var iv = []byte(key)[:aes.BlockSize]
	encrypted := make([]byte, len(data))
	aesBlockEncrypter, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesEncrypter := cipher.NewCTR(aesBlockEncrypter, iv)
	aesEncrypter.XORKeyStream(encrypted, data)
	return encrypted, nil
}

//解密字符串
func (this *AesEncrypt) Decrypt(src []byte) (result []byte, err error) {
	defer func() {
		//错误处理
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()
	key := this.getKey()
	var iv = []byte(key)[:aes.BlockSize]
	decrypted := make([]byte, len(src))
	var aesBlockDecrypter cipher.Block
	aesBlockDecrypter, err = aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}
	aesDecrypter := cipher.NewCTR(aesBlockDecrypter, iv)
	aesDecrypter.XORKeyStream(decrypted, src)
	return decrypted, nil
}
