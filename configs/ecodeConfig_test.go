package configs

import "testing"

func TestGetECodeConst(t *testing.T) {
	ecodeConfigFilePath = "../configs/ecodeConst.json"
	cst, err := GetECodeConst()
	if err != nil {
		t.Error(err)
	}
	if cst.CodeTypeIDByText["禁止出校"] == "8" {
		t.Log("ok")
	}
}
