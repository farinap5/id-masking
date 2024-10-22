package internal

import (
	"bytes"
	"crypto/aes"
	b64 "encoding/base64"
	"strconv"
	"strings"
)

func Init(key string) *Encoder {
	block,err := aes.NewCipher([]byte(key))
	if err != nil {
		panic(err.Error())
	}
	enc := new(Encoder)
	enc.AES = block
	return enc
}

// PKCS7 Padding
func Padding(data []byte, blockSize int) []byte {
	padSize := blockSize - len(data)%blockSize
	padding := bytes.Repeat([]byte{byte(padSize)}, padSize)
	return append(data, padding...)
}
func Unpadding(data []byte) []byte {
	length := len(data)
	unpadding := int(data[length-1])
	return data[:(length - unpadding)]
}

func (e *Encoder) Encode(nounce string, id int) string {
	idStr := strconv.Itoa(id)
	data := []byte(nounce + ":" + idStr)

	blockSize := aes.BlockSize
	data = Padding(data, blockSize)
	ciphertext := make([]byte, len(data))
	e.AES.Encrypt(ciphertext, data)
	b64enc := b64.StdEncoding.EncodeToString(ciphertext)
	return b64enc
}

func (e *Encoder) Decode(data string) (string, int) {
	bts, err := b64.StdEncoding.DecodeString(data)
	if err != nil {
		panic(err.Error())
	}

	deciphertext := make([]byte, len(bts))
	e.AES.Decrypt(deciphertext, bts)
	deciphertext = Unpadding(deciphertext)

	lst := strings.Split(string(deciphertext), ":")
	nounce := lst[0]
	idStr := lst[1]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		panic(err.Error())
	}

	return nounce, id
}