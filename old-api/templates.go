package main

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"

	"github.com/kiasaki/batbelt/mst"
)

func loadLayoutTemplate() *template.Template {
	t, err := template.New("page").Parse(layoutContents)
	mst.MustNotErr(err)
	return t
}

func executeTemplate(t *template.Template, data interface{}) ([]byte, error) {
	var out bytes.Buffer
	err := t.ExecuteTemplate(&out, "layout", data)
	if err != nil {
		return []byte{}, err
	}
	return out.Bytes(), nil
}

type TemplateMap map[string]*template.Template

func (tm *TemplateMap) RenderPage(w http.ResponseWriter, name string, data interface{}) {
	if tmpl, ok := (*tm)[name]; ok {
		if doc, err := executeTemplate(tmpl, data); err == nil {
			w.Header().Set("Content-Type", "text/html")
			w.Write(doc)
		} else {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf("500 - Error occured while rendering page (%s)", name)))
		}
	} else {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("500 - Can't render a template not present in the TemplateMap (%s)", name)))
	}
}
