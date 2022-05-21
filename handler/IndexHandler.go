package handler

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/rroy233/EnterNEU/configs"
	"github.com/rroy233/EnterNEU/utils"
	"io/ioutil"
	"strings"
)

func IndexHandler(c *gin.Context) {
	index, err := ioutil.ReadFile("./assets/enterneu/index.html")
	if err != nil {
		utils.ReturnPlainHtml(c, "模板读取失败")
		return
	}

	//判断是否需要插入统计代码
	//请在config.yaml中设置，本项目默认采用https_v6_51_la提供的统计服务
	if configs.Get().StatisticsReport.V651La == true {
		jss := strings.Join([]string{configs.Get().StatisticsReport.V651LaJs1,
			configs.Get().StatisticsReport.V651LaJs2,
			"</head>",
		}, "")
		index = bytes.Replace(index, []byte("</head>"), []byte(jss), 1)
	}

	c.Data(200, "text/html;charset=utf-8", index)
	return
}
