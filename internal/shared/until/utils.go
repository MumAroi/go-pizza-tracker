package until

import (
	"encoding/json"
	"html/template"

	"github.com/gin-gonic/gin"
)

type Template interface {
	LoadTemplates(router *gin.Engine) error
}

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
