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

const TYPE_INT = "int"
const TYPE_ARRAY = "array"
const TYPE_STRING = "string"
const TYPE_BOOL = "bool"
const TYPE_FLOAT = "float"

type StructItem struct {
	ItemName    string
	ItemType    string
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
		case reflect.Array,reflect.Slice:
			structItem.ItemType = TYPE_ARRAY
		case reflect.Int, reflect.Int8, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8:
			structItem.ItemType = TYPE_INT
		default:
			log.Error("unknowType %s  %s", f.Type.Kind().String(),f.Name)
		}
		moduleInfo.Items = append(moduleInfo.Items, structItem)
	}
	Export(moduleInfo)
}

func Export(m *ModuleInfo) {
	err := os.MkdirAll("./template/public/", os.ModePerm)
	err = os.MkdirAll("./template/backend/", os.ModePerm)
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
	printerBg.F, err = os.OpenFile("./template/backend"+string(filepath.Separator)+m.StructName+"Backend.go", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
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
