package site

import (
    "appengine"
    "html/template"
    "net/http"
)

type Tpl struct {
	Name		string
	File		string
	Tmpl		*template.Template
	Data		[]Data
}

type Data struct{
	Name		string
	Msg			[]Msg
	Input		[]Input
	TxtArea		[]TxtArea
	Sel			[]Sel
}

type Msg struct {
	Name		string
	Text		string
	Class		string
	ClassOk		string
	ClassErr	string
	ClassWarn	string
	ClassAlert	string
}

type Input struct {
	Name		string
	Value		string
	Checked		string
	Extras		Extras
}

type TxtArea struct {
	Name		string
	Value		string
	Cols		int
	Rows		int
	Extras		Extras
}
type Sel struct {
	Name		string
	Size		string
	Opt			[]Opt
	Multiple	string
	Extras		Extras
}

type Opt struct {
	Label		string
	Value		string
	Selected	string
}

type Extras struct {
	Msg			string
	Class		string
	ErrMsg		string
	ErrClass	string
	Disabled	string
	ReadOnly	string
}

var TplPath string

func init() {
	TplPath = "templates/"
}

func (t *Tpl) Add(c appengine.Context, name string, file string) {
	t.Tmpl = template.Must(template.New(name).ParseFiles(TplPath+file))
}

func (t *Tpl) Render(w http.ResponseWriter) {
	if err := t.Tmpl.Execute(w, t.Data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (t *Tpl) AddData(n string) *Data {
	d := &Data { Name: n }
	t.Data = append(t.Data, *d)
	return d
}

func (d *Data) AddMsg(n string, tx string, c string, cok string, cerr string, cwrn string, calt string) Msg {
	msg := Msg { Name: n, Text: tx, Class: c, ClassOk: cok, ClassErr: cerr, ClassWarn: cwrn, ClassAlert: calt }
	d.Msg = append(d.Msg, msg)
	return msg
}

func (d *Data) AddInput(n string, v string, chk string) Input {
	i := Input {Name: n, Value: v, Checked: chk }
	d.Input = append(d.Input, i)
	return i
}

func (d *Data) AddTxtArea(n string, v string, cols int, rows int) TxtArea {
	txt := TxtArea {Name: n, Value: v, Cols: cols, Rows: rows }
	d.TxtArea = append(d.TxtArea, txt)
	return txt
}

func (d *Data) AddSel(n string, s string, m string) Sel {
	sel := Sel {Name: n, Size: s, Multiple: m }
	d.Sel = append(d.Sel, sel)
	return sel
}

func (s *Sel) AddOpt(l string, v string, sel string) Opt {
	opt := Opt {Label: l, Value: v, Selected: sel }
	s.Opt = append(s.Opt, opt)
	return opt
}

// Add Extras falta, para cada tipo
// func (* tipo) addExtras(.....) *tipo { }
