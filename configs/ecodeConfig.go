package configs

import (
	"encoding/json"
	"io/ioutil"
	"strconv"
)

type eCodeConst struct {
	Colors           []string `json:"colors"`
	CodeTypes        []string `json:"codeTypes"`
	CodeTypeIDByText map[string]string
}

var ecodeConfigFilePath = "./configs/ecodeConst.json"

func GetECodeConst() (*eCodeConst, error) {
	fileData, err := ioutil.ReadFile(ecodeConfigFilePath)
	if err != nil {
		return nil, err
	}

	ec := new(eCodeConst)
	err = json.Unmarshal(fileData, ec)
	if err != nil {
		return nil, err
	}

	//生成 codeType.Text -> id的映射
	ec.CodeTypeIDByText = make(map[string]string, 0)
	for i, codeType := range ec.CodeTypes {
		ec.CodeTypeIDByText[codeType] = strconv.FormatInt(int64(i), 10)
	}

	return ec, nil
}
