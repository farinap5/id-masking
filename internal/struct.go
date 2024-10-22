package internal

import "crypto/cipher"

type Encoder struct {
	AES cipher.Block
}