package template

import (
	"text/template"
	"github.com/GodSlave/MyGoServer/log"
	"os"
	"strings"
	"fmt"
	"path/filepath"
	"reflect"
)

type ModuleInfo struct {
	ModuleName string
	StructName string
	Items      []*StructItem
}

const TYPE_INT = 1
const TYPE_ARRAY = 2
const TYPE_STRING = 3
const TYPE_BOOL = 4
const TYPE_FLOAT = 5

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

func BuildModel(mStruct interface{}, moduleName string) {
	moduleInfo := &ModuleInfo{}
	t := reflect.TypeOf(mStruct)
	moduleInfo.StructName = t.Name()
	moduleInfo.ModuleName = moduleName
	for i := 0; i < t.NumField(); i++ {
		structItem := &StructItem{}
		f := t.Field(i)
		structItem.DisPlayName, _ = f.Tag.Lookup("name")
		structItem.ItemName = f.Name
		switch f.Type.Kind() {
		case reflect.Bool:
			structItem.ItemType = TYPE_BOOL
		case reflect.String:
			structItem.ItemType = TYPE_STRING
		case reflect.Float32, reflect.Float64:
			structItem.ItemType = TYPE_FLOAT
		case reflect.Array:
			structItem.ItemType = TYPE_ARRAY
		case reflect.Int, reflect.Int8, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8:
			structItem.ItemType = TYPE_STRING
		default:
			log.Error("unknowType %s", f.Type.Kind().String())
		}
		moduleInfo.Items = append(moduleInfo.Items, structItem)
	}
	Export(moduleInfo)
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
