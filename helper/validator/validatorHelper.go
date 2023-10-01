package validatorHelper

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
)

var _trans ut.Translator

// InitTrans  初始化验证器
func InitTrans(locale string) (err error) {
	// 修改gin中的validator引擎属性,实现自定制
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		zhT := zh.New() // 中文翻译器
		enT := en.New() // 英文翻译器

		// 第一个参数是备用(fallback)的语言环境
		// 后面的参数应该是支持的语言环境(支持多个)
		// uni :=ut.New(zhT,zhT)也是可以的
		uni := ut.New(enT, zhT, enT)

		// locale 取决于http请求头的Accept-Language
		var ok bool
		// 也可以使用uni.FindTranslator(...)传入多个locale进行查找
		_trans, ok = uni.GetTranslator(locale)
		if !ok {
			panic("语言环境获取失败" + locale)
		}
		// 注册翻译器
		switch locale {
		case "en":
			err = enTranslations.RegisterDefaultTranslations(v, _trans)
		case "zh":
			err = zhTranslations.RegisterDefaultTranslations(v, _trans)
		default:
			err = enTranslations.RegisterDefaultTranslations(v, _trans)
		}
		fmt.Printf("%s翻译定制成功", locale)
		return nil
	}
	return nil
}

// GetTrans 获取当前翻译接口
func GetTrans() ut.Translator {
	return _trans
}
