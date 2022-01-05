package main

import (
	"bytes"
	"text/template"
)

func renderTemplate(textTemplate string, data any) string {
	tmpl, err := template.New("Template").Parse(textTemplate)
	if err != nil {
		panic(err)
	}
	var out bytes.Buffer
	err = tmpl.Execute(&out, data)
	if err != nil {
		panic(err)
	}

	return out.String()
}

var functionTemplate = `
func {{.Name}}() itertools.Iterator[int] {
	__next := 0
	var __zero int
	__nop := func() {}
	advance := func() (bool, int) {
		switch __next {
		{{range .NextIndices}}
		case {{.}}:
			goto __next_{{.}}
		{{end}}
		}
		
	__next_0:

		{{.Body}}

		return false, __zero
	}
	return itertools.FromAdvance(advance)
}
`

func renderFunction(name, body string, stateCount int) string {
	nextIndices := make([]int, stateCount+1)
	for i := range nextIndices {
		nextIndices[i] = i
	}

	return renderTemplate(functionTemplate, struct {
		Name, Body  string
		NextIndices []int
	}{Name: name, Body: body, NextIndices: nextIndices})
}

var returnTemplate = `
	__next = {{.StateId}}
	return true, {{.ReturnValue}}
__next_{{.StateId}}:
	__nop()
`

func renderReturn(stateId int, returnValue string) string {
	return renderTemplate(returnTemplate, struct {
		StateId     int
		ReturnValue string
	}{StateId: stateId, ReturnValue: returnValue})
}
