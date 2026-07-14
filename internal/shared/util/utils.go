package util

import (
	"encoding/json"
	"html/template"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func LoadTemplate(router *gin.Engine) error {
	functions := template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
		"json": func(v interface{}) template.JS {
			b, _ := json.Marshal(v)
			return template.JS(b)
		}}
	temp, err := template.New("").Funcs(functions).ParseGlob("templates/*.tmpl")

	if err != nil {
		return err
	}

	router.SetHTMLTemplate(temp)
	return nil

}

func SetSessionValue(c *gin.Context, key string, value any) error {
	session := sessions.Default(c)
	session.Set(key, value)
	return session.Save()
}

func ClearSession(c *gin.Context) error {
	session := sessions.Default(c)
	session.Clear()
	return session.Save()
}

func GetSessionString(c *gin.Context, key string) string {
	session := sessions.Default(c)
	val := session.Get(key)
	if val == nil {
		return ""
	}

	str, _ := val.(string)
	return str
}
