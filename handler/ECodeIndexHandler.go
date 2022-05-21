package handler

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/rroy233/EnterNEU/configs"
	"github.com/rroy233/EnterNEU/utils"
	"io/ioutil"
	"strings"
)

func ECodeIndexHandler(c *gin.Context) {
	data, err := ioutil.ReadFile("./assets/ecode/index.html")
	if err != nil {
		utils.ReturnPlainHtml(c, "找不到./assets/ecode/index.html")
		return
	}

	//判断是否需要插入统计代码
	//请在config.yaml中设置，本项目默认采用https_v6_51_la提供的统计服务
	//统计信息中不会包含cookies等隐私信息
	if configs.Get().StatisticsReport.V651La == true {
		jss := strings.Join([]string{configs.Get().StatisticsReport.V651LaJs1,
			configs.Get().StatisticsReport.V651LaJs2,
			"</head>",
		}, "")
		data = bytes.Replace(data, []byte("</head>"), []byte(jss), 1)
	}

	c.Data(200, gin.MIMEHTML, data)
	return
}
