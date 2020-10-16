package main

import (
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	"image"
	_ "image/jpeg"
	"net/http"
	"net/url"
	"strconv"

	"decode-qrcode/decode"
	"decode-qrcode/utils"
	"github.com/gin-gonic/gin"
	logs "github.com/sirupsen/logrus"
)

func main() {
	configInfo := decode.New()
	configInfo.Parse(utils.GetConfigFile())

	if len(configInfo.Whitelist) <= 0 {
		logs.WithFields(logs.Fields{
			"module": "Main方法",
		}).Fatalln("必须配置白名单")
	}

	logs.WithFields(logs.Fields{
		"module": "Main方法",
	}).Println("初始化程序...")

	route := gin.Default()

	// ip合法性中间件
	route.Use(func(context *gin.Context) {
		clientIp := context.ClientIP()
		for _, whiteIp := range configInfo.Whitelist {
			if clientIp == whiteIp {
				return
			}
		}
		context.AbortWithStatusJSON(http.StatusOK,
			utils.DecodeResultJson(http.StatusInternalServerError, "IP不合法，请联系管理员"))
	})

	route.GET("/api/internal/qrcode/decode.htm", func(context *gin.Context) {
		imgUrl := context.Query("imgUrl")
		if imgUrl == "" {
			context.JSON(http.StatusInternalServerError,
				utils.DecodeResultJson(http.StatusInternalServerError, "imgUrl参数不能为空"))
			return
		}
		imgUrl, _ = url.QueryUnescape(imgUrl)
		// 判断url是否合法
		url, _ := url.Parse(imgUrl)
		if !url.IsAbs() {
			context.JSON(http.StatusInternalServerError,
				utils.DecodeResultJson(http.StatusInternalServerError, "imgUrl参数中URL不是合法URL"))
			return
		}

		response, err := http.Get(imgUrl)
		if err != nil {
			context.JSON(http.StatusInternalServerError,
				utils.DecodeResultJson(http.StatusInternalServerError, "解码imgUrl参数异常，异常信息:"+err.Error()))
			return
		}
		if response != nil {
			defer response.Body.Close()
		}
		// 解码二维码
		img, _, _ := image.Decode(response.Body)
		bmp, _ := gozxing.NewBinaryBitmapFromImage(img)
		qrReader := qrcode.NewQRCodeReader()
		qrcodeContent, _ := qrReader.Decode(bmp, nil)
		context.JSON(http.StatusOK, utils.DecodeResultJson(http.StatusOK, qrcodeContent.String()))
	})

	route.Run(":" + strconv.Itoa(configInfo.Port))
}
