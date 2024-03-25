package dao

import (
	"bytes"
	"encoding/base64"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"image"
	"image/png"
	"strconv"
)

var Key *otp.Key

func ImageToBase64(img image.Image) (string, error) {
	// 创建一个缓冲区
	var buf bytes.Buffer

	// 将图片数据写入缓冲区
	if err := png.Encode(&buf, img); err != nil {
		return "", err
	}

	// 将图片数据转换为Base64编码
	encodedStr := base64.StdEncoding.EncodeToString(buf.Bytes())

	return encodedStr, nil
}

func GenerateQRCode(uid int64) (img image.Image, err error) {
	Key, err = totp.Generate(totp.GenerateOpts{
		Issuer: strconv.FormatInt(uid, 10),
	})
	if err != nil {
		return
	}

	//生成QR code
	img, err = Key.Image(200, 200)
	if err != nil {
		return
	}
	return
}
