package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// GetErrorMsg 获取验证错误信息
func GetErrorMsg(err error) string {
	if err == nil {
		return ""
	}

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return err.Error()
	}

	var messages []string
	for _, e := range validationErrors {
		messages = append(messages, getSingleMessage(e))
	}

	return strings.Join(messages, "; ")
}

func getSingleMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s 是必填的", translateField(e.Field()))
	case "min":
		return fmt.Sprintf("%s 长度不能少于 %s", translateField(e.Field()), e.Param())
	case "max":
		return fmt.Sprintf("%s 长度不能超过 %s", translateField(e.Field()), e.Param())
	case "email":
		return fmt.Sprintf("%s 格式不正确", translateField(e.Field()))
	case "url":
		return fmt.Sprintf("%s 格式不正确", translateField(e.Field()))
	case "numeric":
		return fmt.Sprintf("%s 必须是数字", translateField(e.Field()))
	case "between":
		return fmt.Sprintf("%s 必须在 %s 范围内", translateField(e.Field()), e.Param())
	default:
		return fmt.Sprintf("%s 验证失败", translateField(e.Field()))
	}
}

func translateField(field string) string {
	translations := map[string]string{
		"Username":      "用户名",
		"Password":      "密码",
		"Email":         "邮箱",
		"Nickname":      "昵称",
		"Title":         "标题",
		"Content":       "内容",
		"Description":   "描述",
		"Name":          "名称",
		"URL":           "链接",
		"OldPassword":   "旧密码",
		"NewPassword":   "新密码",
		"ConfirmPassword": "确认密码",
	}

	if translated, ok := translations[field]; ok {
		return translated
	}
	return field
}
