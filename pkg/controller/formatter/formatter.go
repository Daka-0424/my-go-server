package formatter

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/gin-gonic/gin"
)

var DEFAULT_KEY = []byte{
	0x8E, 0x09, 0x6B, 0xD0, 0x40, 0x43, 0xDF, 0xAB, 0x31, 0xE6, 0x97, 0x40, 0x5E, 0x4B, 0x86, 0xA8,
	0xC3, 0xFE, 0x43, 0x85, 0xD3, 0x21, 0x19, 0x10, 0x1A, 0x38, 0xE1, 0x38, 0xE0, 0x09, 0x03, 0x9D,
}

var DEFAULT_IV = []byte{
	0x51, 0xFA, 0xCE, 0x3C, 0x7D, 0xB7, 0xDC, 0xDF, 0x33, 0x00, 0xD9, 0x71, 0x7B, 0xB4, 0x3D, 0x02,
}

const (
	KEY_SIZE = 32
	IV_SIZE  = 16
)

const CRYPTO_CACHE_KEY = "crypto_"

func Respond(ctx *gin.Context, status int, v any) {
	if strings.Contains(ctx.Request.Header.Get("Accept"), gin.MIMEJSON) {
		ctx.JSON(status, v)
		return
	}

	key, iv := getKeyAndIV(ctx)

	json, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}

	encrypted, err := Encrypt(json, key, iv)
	if err != nil {
		panic(err)
	}

	ctx.Writer.Header().Add("Context-Type", "application/octet_stream")
	ctx.Writer.WriteHeader(status)
	ctx.Writer.Write(encrypted)
}

func ShouldBind(ctx *gin.Context, v any) error {
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		return err
	}

	if ctx.Request.Header.Get("Content-Type") == gin.MIMEJSON {
		return json.Unmarshal(body, v)
	}

	key, iv := getKeyAndIV(ctx)

	decrypted, err := Decrypt(body, key, iv)
	if err != nil {
		return err
	}

	return json.Unmarshal(decrypted, v)
}

func Encrypt(data, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	paddedData := pad(data)
	stream := cipher.NewCBCEncrypter(block, iv)
	encrypted := make([]byte, len(paddedData))
	stream.CryptBlocks(encrypted, paddedData)

	return encrypted, nil
}

func Decrypt(data, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	stream := cipher.NewCBCDecrypter(block, iv)
	decrypted := make([]byte, len(data))
	stream.CryptBlocks(decrypted, data)

	return unpad(decrypted), nil
}

func pad(data []byte) []byte {
	padSize := aes.BlockSize - (len(data) % aes.BlockSize)
	return append(data, bytes.Repeat([]byte{byte(padSize)}, padSize)...)
}

func unpad(data []byte) []byte {
	padSize := int(data[len(data)-1])
	return data[:len(data)-padSize]
}

func getKeyAndIV(ctx *gin.Context) ([]byte, []byte) {
	v, ok := ctx.Get("cryptoKey")
	if !ok {
		return DEFAULT_KEY, DEFAULT_IV
	}

	key := v.([]byte)

	v, ok = ctx.Get("cryptoIV")
	if !ok {
		return DEFAULT_KEY, DEFAULT_IV
	}

	iv := v.([]byte)

	return key, iv
}
