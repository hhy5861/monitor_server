package utils

import (
	"encoding/json"
	"crypto/md5"
	"encoding/hex"
)

func GetMD5Hash(data interface{}) string {
	bytes, err := json.Marshal(data)
	if err == nil {
		md5Ctx := md5.New()
		md5Ctx.Write(bytes)
		cipherStr := md5Ctx.Sum(nil)

		return hex.EncodeToString(cipherStr)
	}

	return ""
}
