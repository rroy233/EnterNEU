package configs

import (
	"encoding/json"
	"io/ioutil"
)

type eCodeConst struct {
	Colors    []string `json:"colors"`
	CodeTypes []string `json:"codeTypes"`
}

func GetECodeConst() (*eCodeConst, error) {
	fileData, err := ioutil.ReadFile("./configs/ecodeConst.json")
	if err != nil {
		return nil, err
	}

	ec := new(eCodeConst)
	err = json.Unmarshal(fileData, ec)
	if err != nil {
		return nil, err
	}
	return ec, nil
}
