package formatter

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"

	"github.com/Daka-0424/my-go-server/language"
	"github.com/Daka-0424/my-go-server/pkg/usecase/model"
	"github.com/Daka-0424/my-go-server/pkg/usecase/model/response"
	"github.com/Songmu/flextime"
	"github.com/gin-gonic/gin"
)

var DEFAULT_KEY = []byte{
	0x8E, 0x09, 0x6B, 0xD0, 0x40, 0x43, 0xDF, 0xAB, 0x31, 0xE6, 0x97, 0x40, 0x5E, 0x4B, 0x86, 0xA8,
	0xC3, 0xFE, 0x43, 0x85, 0xD3, 0x21, 0x19, 0x10, 0x1A, 0x38, 0xE1, 0x38, 0xE0, 0x09, 0x03, 0x9D,
}

var DEFAULT_IV = []byte{
	0xFA, 0xCE, 0x3C, 0x7D, 0xB7, 0xDC, 0xDF, 0x33, 0x00, 0xD9, 0x71, 0x7B, 0xB4, 0x3D, 0x02, 0x51,
}

const (
	KEY_SIZE = 32
	IV_SIZE  = 16
)

const CRYPTO_CACHE_KEY = "crypto_"

func Respond(ctx *gin.Context, status int, v any) {
	// Serverの時間をHeaderに入れて返す
	now := flextime.Now()
	ctx.Writer.Header().Set("Access-Control-Expose-Headers", "Server-Time")
	ctx.Writer.Header().Add("Server-Time", now.Format("2006-01-02 15:04:05"))

	json, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	log.Println("Response", status, string(json))

	if strings.Contains(ctx.Request.Header.Get("Accept"), gin.MIMEJSON) {
		ctx.JSON(status, v)
		return
	}

	key, iv := getKeyAndIV(ctx)
	encrypted, err := Encrypt(json, key, iv)
	if err != nil {
		panic(err)
	}

	ctx.Writer.Header().Add("Content-Type", "application/octet-stream")
	ctx.Writer.WriteHeader(status)
	ctx.Writer.Write(encrypted)
}

func ShouldBind(ctx *gin.Context, v any, localizar *language.Localizer) *response.AppError {
	lang := ctx.Request.Header.Get("Accept-Language")
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		return response.NewErrBadRequest(model.E9901, localizar.MustLocalize(model.E9901, lang, nil))
	}

	if ctx.Request.Header.Get("Content-Type") != gin.MIMEJSON {
		key, iv := getKeyAndIV(ctx)

		decrypted, err := Decrypt(body, key, iv)
		if err != nil {
			return response.NewErrBadRequest(model.E9901, localizar.MustLocalize(model.E9901, lang, nil))
		}
		body = decrypted
	}

	if err := json.Unmarshal(body, v); err != nil {
		return response.NewErrBadRequest(model.E9901, localizar.MustLocalize(model.E9901, lang, nil))
	}

	log.Println("Request", ctx.Request.Method, ctx.Request.URL.Path, string(body))

	return nil
}

func pad(data []byte) []byte {
	padSize := aes.BlockSize - (len(data) % aes.BlockSize)
	return append(data, bytes.Repeat([]byte{byte(padSize)}, padSize)...)
}

func unpad(data []byte) []byte {
	padSize := int(data[len(data)-1])
	return data[:len(data)-padSize]
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
