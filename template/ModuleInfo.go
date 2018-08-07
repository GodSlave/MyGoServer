package template

import (
	"text/template"
	"github.com/GodSlave/MyGoServer/log"
	"os"
	"strings"
	"fmt"
	"path/filepath"
)

type ModuleInfo struct {
	ModuleName string
	StructName string
	Items      []*StructItem

}

const TYPE_INT = 1
const TYPE_ARRAY = 2
const TYPE_STRING = 3

type StructItem struct {
	ItemName    string
	ItemType    int
	DisPlayName string
}

type ModulesItem struct {
	ItemName    string
	DisPlayName string
}

type Printer struct {
	F *os.File
}

func (printer *Printer) Write(p []byte) (n int, err error) {
	content := string(p)
	content = strings.Replace(content, "&#34;", "\"", -1)
	content = strings.Replace(content, "&lt;", "<", -1)
	_, err = printer.F.Write([]byte(content))
	if err != nil {
		fmt.Println(err.Error())
	}
	return 0, err
}

func Export(m *ModuleInfo) {

	tplhtml, err := template.New("serverhtml.template").Parse(HtmlContent)
	if err != nil {
		log.Error(err.Error())
	}
	var printer = Printer{}
	printer.F, err = os.OpenFile("./template/public"+string(filepath.Separator)+m.StructName+".html", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	err = tplhtml.Execute(&printer, m)
	if err != nil {
		log.Error(err.Error())
	}

	tplBg, err := template.New("serverbg.template").Parse(BGWebContent)
	if err != nil {
		log.Error(err.Error())
	}
	var printerBg = Printer{}
	printerBg.F, err = os.OpenFile("./template/public"+string(filepath.Separator)+m.StructName+"Backend.go", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	err = tplBg.Execute(&printerBg, m)
	if err != nil {
		log.Error(err.Error())
	}

	tplJS, err := template.New("serverjs.template").Parse(JSContent)
	if err != nil {
		log.Error(err.Error())
	}
	var printerJs = Printer{}
	printerJs.F, err = os.OpenFile("./template/public"+string(filepath.Separator)+m.StructName+".js", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	err = tplJS.Execute(&printerJs, m)
	if err != nil {
		log.Error(err.Error())
	}

}
