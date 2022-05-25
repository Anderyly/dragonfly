/*
 * @Author anderyly
 * @email admin@aaayun.cc
 * @link http://blog.aaayun.cc/
 * @copyright Copyright (c) 2022
 */

package ay

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"log"
	"reflect"

	//en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

type Validator struct {
}

var (
	uni *ut.UniversalTranslator
	//validate *validator.Validate
	trans ut.Translator
)

func init() {
	//注册翻译器
	zh := zh.New()
	uni = ut.New(zh, zh)
	trans, _ = uni.GetTranslator("zh")
	validate := binding.Validator.Engine().(*validator.Validate)
	err := zh_translations.RegisterDefaultTranslations(validate, trans)
	if err != nil {
		log.Println(err)
	}
	//en_translations.RegisterDefaultTranslations(validate, trans)
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := fld.Tag.Get("label")
		return name
	})
}

//Translate 翻译错误信息
func (class Validator) Translate(err error) string {
	//var result = make(map[string][]string)
	errors := err.(validator.ValidationErrors)
	res := ""
	for _, err := range errors {
		//result[err.Field()] = append(result[err.Field()], err.Translate(trans))
		res += err.Translate(trans)
	}
	return res
}
