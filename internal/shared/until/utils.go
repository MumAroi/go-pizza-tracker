package until

import (
	"html/template"

	"github.com/gin-gonic/gin"
)

func loadTemplate(router *gin.Engine) error {
	functions := template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
	}
	temp, err := template.New("").Funcs(functions).ParseGlob("templates/*.tmpl")

	if err != nil {
		return err
	}

	router.SetHTMLTemplate(temp)
	return nil

}
