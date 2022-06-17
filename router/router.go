package router

import (
	"github.com/gin-gonic/gin"
	"github.com/rroy233/EnterNEU/handler"
	"github.com/rroy233/EnterNEU/middlewares"
)

func Register(engine *gin.Engine) {
	//e马桶
	eCodeGroup := engine.Group("/ecode")
	eCodeGroup.Use(middlewares.ECodeMiddleWare)
	{
		eCodeGroup.Static("/js", "./assets/ecode/js")
		eCodeGroup.Static("/css", "./assets/ecode/css")
		eCodeGroup.Static("/fonts", "./assets/ecode/fonts")
		eCodeGroup.Static("/img", "./assets/ecode/img")

		eCodeGroup.GET("/", handler.ECodeIndexHandler)

		eCodeGroup.GET("/api/user-info", handler.ECodeAPIUserInfo)
		eCodeGroup.GET("/api/register", handler.ECodeAPIRegister)
		eCodeGroup.GET("/api/user-face", handler.ECodeAPIUserFace)

		eCodeGroup.GET("/api/qr-code", handler.ECodeAPIQrCodeHandler)
		eCodeGroup.GET("/api/person-health", handler.ECodeAPIPersonHealthHandler)
		eCodeGroup.GET("/api/user-priority", handler.ECodeAPIUserPriorityHandler)
		eCodeGroup.GET("/api/user-vaccination", handler.ECodeAPIUserVaccinationHandler)
		eCodeGroup.GET("/api/user-examination", handler.ECodeAPIUserExaminationHandler)

	}

	//fuckNeu
	engine.GET("/", handler.IndexHandler)
	engine.Static("/js", "./assets/enterneu/js")
	engine.Static("/css", "./assets/enterneu/css")
	apiGroup := engine.Group("/api")
	{
		apiGroup.GET("/:token/:key", handler.TokenCheckHandler)
		apiGroup.GET("/:token/:key/shadowrocket", handler.TokenShadowrocketHandler)

		//头像预览
		apiGroup.GET("/viewImage", handler.APIViewImageHandler)

		//获取说明MarkDown
		apiGroup.GET("/tips", handler.APITipsHandler)

		//生成随机秘钥
		apiGroup.GET("/randKey", handler.APIRandKeyHandler)

		//存储基本信息
		apiGroup.POST("/store", handler.APIStoreHandler)

		//获取配置
		apiGroup.GET("/config", handler.APIConfigHandler)

		//上传头像
		apiGroup.POST("/upload", handler.APIUploadHandler)

		//获取存储的信息
		apiGroup.GET("/:token/:key/status", handler.APIStatusHandler)
		//获取存储的信息
		apiGroup.POST("/:token/:key/delete", handler.APIDeleteHandler)
	}

}
