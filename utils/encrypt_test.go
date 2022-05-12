package utils

import (
	"github.com/google/uuid"
	"strings"
	"testing"
)

func TestAesCbcEncrypt(t *testing.T) {
	configs.InitConfig("../config.yaml")
	id, _ := uuid.NewUUID()
	key := strings.Replace(id.String(), "-", "", -1)
	t.Log("key:", key)
	out, err := AesCbcEncryptString("test", key)
	if err != nil {
		t.Error(err)
	}
	t.Log("加密结果：" + out)
}

func TestAesCbcDecrypt(t *testing.T) {
	text := "kGQXg15MubiJG-Ty2cnntQ=="
	key := "0a8a5b82c94411ec9d02acbc327cbb49"
	configs.InitConfig("../config.yaml")
	out, err := AesCbcDecryptString(text, key)
	if err != nil {
		t.Error(err)
	}
	t.Log("解密结果：" + out)
}
