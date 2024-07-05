package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tianaiyouqing/captcha-go-demo/config"
	"github.com/tianaiyouqing/tianai-captcha-go/common/model"
	"net/http"
	"strings"
)

func main() {
	fmt.Println(config.Captcha)
	engine := gin.Default()
	engine.Static("/", "./static")

	engine.POST("/gen", func(c *gin.Context) {
		param := c.Query("type")
		if param == "" {
			param = "slider"
		}
		captcha, err := config.Captcha.GenerateCaptcha(&model.GenerateParam{
			CaptchaName: strings.ToUpper(param),
		})
		if err != nil {
			c.JSON(500, gin.H{
				"code": 500,
				"msg":  err.Error(),
			})
			return
		}
		c.JSON(200, gin.H{
			"code":    200,
			"msg":     "success",
			"id":      captcha.Id,
			"captcha": captcha,
		})
	})

	engine.POST("/check", func(c *gin.Context) {
		validParam := new(ValidParam)
		if err := c.ShouldBindJSON(validParam); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		valid, err := config.Captcha.Valid(validParam.Id, &validParam.Data)
		if err != nil {
			c.JSON(500, gin.H{
				"code": 500,
				"msg":  err.Error(),
			})
			return
		}
		if valid.Code == 200 {
			// 验证码验证成功,存自定义token, 后续业务处理,这里方便演示，用captcha的cacheStore演示
			token := "token_" + validParam.Id
			_ = config.Captcha.CacheStore.SetCache(token, map[string]any{}, nil)
			// 返回token给前端，用作二次验证
			c.JSON(200, gin.H{
				"code":  200,
				"msg":   "success",
				"token": token,
			})
		} else {
			c.JSON(200, valid)
		}
	})

	engine.Run(":8080")
}

type ValidParam struct {
	Id   string                  `json:"id" binding:"required"`
	Data model.ImageCaptchaTrack `json:"data" binding:"required"`
}
