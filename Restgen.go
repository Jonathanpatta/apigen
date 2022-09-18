package apigen

import (
	"fmt"
)

func GetGen(service *ApiObject, obj *ApiObject) string {
	forContent := IfGen("row.Id == id", fmt.Sprintf("return s.%s[i],nil", obj.GetLocalStoreName()))

	content := ForGen("i,row", "s."+obj.GetLocalStoreName(), forContent)

	content += fmt.Sprintf("return nil,nil\n")

	code := FunctionGen(
		service.Name,
		"s",
		"Get",
		[]string{"id string"},
		[]string{"*" + obj.Name, "error"},
		content,
	)
	return code
}

func CreateGen(service *ApiObject, obj *ApiObject) string {
	//code := fmt.Sprintf("func (s *%s) Create(obj ...*%s) ([]*%s,error){\n", service.Name, obj.Name, obj.Name)
	//code += fmt.Sprintf("\ts.%s = append(s.%s, obj...)\n", obj.GetLocalStoreName(), obj.GetLocalStoreName())
	//code += fmt.Sprintf("\treturn obj,nil\n")
	//code += "}\n"

	content := fmt.Sprintf("s.%s = append(s.%s, obj...)\n", obj.GetLocalStoreName(), obj.GetLocalStoreName())
	content += fmt.Sprintf("return obj,nil\n")

	code := FunctionGen(
		service.Name,
		"s",
		"Create",
		[]string{"obj ...*" + obj.Name},
		[]string{"[]*" + obj.Name, "error"},
		content,
	)
	return code
}

func UpdateGen(service *ApiObject, obj *ApiObject) string {

	ifContent := fmt.Sprintf("s.%s[i] = obj\n", obj.GetLocalStoreName())
	ifContent += fmt.Sprintf("return s.%s[i],nil", obj.GetLocalStoreName())
	content := ForGen("i,row", "s."+obj.GetLocalStoreName(), IfGen("row.Id == obj.Id", ifContent))
	content += "return nil,nil"
	code := FunctionGen(
		service.Name,
		"s",
		"Update",
		[]string{"obj *" + obj.Name},
		[]string{"*" + obj.Name, "error"},
		content,
	)
	return code
}

func DeleteGen(service *ApiObject, obj *ApiObject) string {

	ifContent := fmt.Sprintf("s.%s = remove(s.%s,i)\n", obj.GetLocalStoreName(), obj.GetLocalStoreName())
	ifContent += fmt.Sprintf("return row,nil")
	content := ForGen("i,row", "s."+obj.GetLocalStoreName(), IfGen("row.Id == id", ifContent))
	content += "return nil,nil"

	code := FunctionGen(
		service.Name,
		"s",
		"Delete",
		[]string{"id string"},
		[]string{"*" + obj.Name, "error"},
		content,
	)
	return code
}

func ListGen(service *ApiObject, obj *ApiObject) string {
	//code := fmt.Sprintf("func (s *%s) List() ([]*%s,error){\n", service.Name, obj.Name)
	//code += fmt.Sprintf("\treturn s.%s,nil\n", obj.GetLocalStoreName())
	//code += "}\n"
	code := FunctionGen(
		service.Name,
		"s",
		"List",
		[]string{},
		[]string{"[]*" + obj.Name, "error"},
		fmt.Sprintf("return s.%s,nil", obj.GetLocalStoreName()),
	)
	return code
}
