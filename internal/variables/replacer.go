package variables

import (
	"bytes"
	"text/template"
)

type Replacer struct {
	vars *TemplateVars
}

func NewReplacer(vars *TemplateVars) *Replacer {
	return &Replacer{vars: vars}
}

func (r *Replacer) Replace(content []byte) ([]byte, error) {
	tmpl, err := template.New("content").Parse(string(content))
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, r.vars); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}