package formatter

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/Daka-0424/my-go-server/config"
	"github.com/gin-gonic/gin"
)

const (
	KEY_SIZE = 32
	IV_SIZE  = 16
)

const CRYPTO_CACHE_KEY = "crypto_"

func Respond(ctx *gin.Context, cfg *config.Config, status int, v any) {
	if cfg.IsDevelopment() {
		if strings.Contains(ctx.Request.Header.Get("Accept"), gin.MIMEJSON) {
			ctx.JSON(status, v)
			return
		}
	}

	key, iv := getKeyAndIV(ctx, cfg)

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

func ShouldBind(ctx *gin.Context, cfg *config.Config, v any) error {
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		return err
	}

	if cfg.IsDevelopment() {
		if ctx.Request.Header.Get("Content-Type") == gin.MIMEJSON {
			return json.Unmarshal(body, v)
		}
	}

	key, iv := getKeyAndIV(ctx, cfg)

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

func getKeyAndIV(ctx *gin.Context, cfg *config.Config) ([]byte, []byte) {
	v, ok := ctx.Get("cryptoKey")
	if !ok {
		return cfg.RequestKeyIv.DefaultKey, cfg.RequestKeyIv.DefaultIv
	}

	key := v.([]byte)

	v, ok = ctx.Get("cryptoIV")
	if !ok {
		return cfg.RequestKeyIv.DefaultKey, cfg.RequestKeyIv.DefaultIv
	}

	iv := v.([]byte)

	return key, iv
}
