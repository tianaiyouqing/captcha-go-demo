package config

import (
	"github.com/tianaiyouqing/tianai-captcha-go/application"
	"github.com/tianaiyouqing/tianai-captcha-go/common/model"
	"github.com/tianaiyouqing/tianai-captcha-go/resource"
	"os"
	"time"
)

var Captcha *application.TianAiCaptchaApplication

func init() {
	builder := application.NewBuilder()

	store := resource.NewMemoryImageCaptchaResourceStore()

	// 导入系统自带的模板
	for _, template := range resource.GetDefaultSliderTemplates() {
		store.AddTemplate("SLIDER", template)
	}
	for _, template := range resource.GetDefaultRotateTemplate() {
		store.AddTemplate("ROTATE", template)
	}

	// 导入自定义背景图片，这里演示吧 resources下的图片导入
	folder, _ := os.ReadDir("./resources")
	for _, file := range folder {
		R := &model.Resource{
			ResourceType: "file",
			Data:         "./resources/" + file.Name(),
		}
		store.AddResource("SLIDER", R)
		store.AddResource("ROTATE", R)
		store.AddResource("WORD_IMAGE_CLICK", R)
	}

	// 设置资源存储器
	builder.SetResourceStore(store)
	// 设置缓冲存储器， 默认是内存存储器， 如需要扩展redis之类， 可自定义实现 `application.CacheStore` 接口
	builder.SetCacheStore(application.NewMemoryCacheStore(5*time.Minute, 5*time.Minute))
	// 添加验证码生成器 滑块验证码、 旋转验证码、 文字点击验证码
	builder.AddProvider(application.CreateSliderProvider())
	builder.AddProvider(application.CreateRotateProvider())
	builder.AddProvider(application.CreateWordClickProvider(nil))
	Captcha = builder.Build()
}
