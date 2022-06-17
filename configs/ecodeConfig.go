package configs

import (
	"encoding/json"
	"io/ioutil"
	"strconv"
)

type eCodeConst struct {
	Colors              []string `json:"colors"`
	CodeTypes           []string `json:"codeTypes"`
	Announcement        string   `json:"announcement"`
	CodeTypeIndexByText map[string]string
	CodeTypeTextByIndex map[string]string
	ColorByCodeTypeID   map[string]string
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

	//初始化map
	ec.CodeTypeIndexByText = make(map[string]string)
	ec.CodeTypeTextByIndex = make(map[string]string)
	ec.ColorByCodeTypeID = make(map[string]string)

	//生成 codeType.Text -> id的映射
	ec.CodeTypeIndexByText = make(map[string]string, 0)
	for i, codeType := range ec.CodeTypes {
		ec.CodeTypeIndexByText[codeType] = strconv.FormatInt(int64(i), 10)
		ec.CodeTypeTextByIndex[strconv.FormatInt(int64(i), 10)] = codeType
		ec.ColorByCodeTypeID[strconv.FormatInt(int64(i), 10)] = ec.Colors[i]
	}

	return ec, nil
}
