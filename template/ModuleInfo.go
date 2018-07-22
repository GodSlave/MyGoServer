package template

type ModuleInfo struct {
	ModuleName string
	StructName string
	Items      []*StructItem
}

type StructItem struct {
	ItemName string
	ItemType int
}
