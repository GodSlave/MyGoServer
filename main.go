package main

import "github.com/GodSlave/MyGoServer/template"

var (
	Version string
	Build   string
)

func main() {
	//app := app2.NewApp()
	//web_Module := web.Module()
	//webModule := web_Module.(*web.ModuleWeb)
	//webModule.Router = gin.Default()
	//app.Run(gate.Module(), userModule.Module(), httpGate.Module(), web_Module)

	template.Export(&template.ModuleInfo{
		ModuleName: "test",
		StructName: "info",
		Items: []*template.StructItem{
			&template.StructItem{
				ItemName: "Id",
				ItemType: template.TYPE_INT,
			},
			&template.StructItem{
				ItemName: "Name",
				ItemType: template.TYPE_STRING,
			},
		},
	})

}
