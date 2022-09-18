package apigen

import (
	"fmt"
	"github.com/gobeam/stringy"
	"strings"
)

var structMp = map[string]bool{}

func (obj *ApiObject) GetLocalStoreName() string {
	return fmt.Sprintf("%s_LocalStore", obj.Name)
}

func (obj *ApiObject) GetServiceInterfaceName() string {
	return fmt.Sprintf("%sService", obj.Name)
}

func (obj *ApiObject) GetLocalServiceName() string {
	return fmt.Sprintf("%sLocalService", obj.Name)
}

func (obj *ApiObject) GetHttpServiceName() string {
	return fmt.Sprintf("%sHttpService", obj.Name)
}

func CreateJsonMarshaller(obj *ApiObject) string {
	content := fmt.Sprintf("res,err := json.Marshal(s)\n")
	content += HandleErrorGen("\"\"", "err")
	content += "return string(res),nil"
	code := FunctionGen(
		obj.Name,
		"s",
		"Marshal",
		[]string{},
		[]string{"string,error"},
		content,
	)

	return code
}

func DeleteFromArrayGen(obj *ApiObject) string {
	code := fmt.Sprintf("func remove(s []*%s, i int) []*%s {\n\ts[i] = s[len(s)-1]\n\treturn s[:len(s)-1]\n}\n", obj.Name, obj.Name)
	return code
}

func FunctionGen(structName string, structObjectName string, Name string, arguments []string, returnTypes []string, content string) string {
	code := "func "
	if structName != "" {
		code += fmt.Sprintf("(%s *%s) ", structObjectName, structName)
	}
	code += fmt.Sprintf("%s(%s) ", Name, strings.Join(arguments, ","))
	if len(returnTypes) > 0 {
		code += fmt.Sprintf("(%s){\n", strings.Join(returnTypes, ","))
	} else {
		code += "{\n"
	}

	lines := strings.Split(content, "\n")
	for i, line := range lines {
		if i == len(lines)-1 && lines[i] == "" {
			break
		}
		lines[i] = "\t" + line
	}
	code += strings.Join(lines, "\n")
	code += "\n}\n\n"
	return code
}

func IfGen(condition string, content string) string {

	code := fmt.Sprintf("if %s {\n", condition)
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		if i == len(lines)-1 && lines[i] == "" {
			break
		}
		lines[i] = "\t" + line
	}
	code += strings.Join(lines, "\n")
	code += "\n}\n"
	return code
}

func ForGen(values string, iterator string, content string) string {

	code := fmt.Sprintf("for %s := range %s {\n", values, iterator)
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		if i == len(lines)-1 && lines[i] == "" {
			break
		}
		lines[i] = "\t" + line
	}
	code += strings.Join(lines, "\n")
	code += "\n}\n"
	return code
}

func HandleErrorGen(returnValues ...string) string {
	content := fmt.Sprintf("return %s", strings.Join(returnValues, ","))
	code := IfGen("err != nil", content)
	return code
}

func HandleHttpErrorGen(code string) string {
	return fmt.Sprintf("if err!= nil{\n\thttp.Error(w,err.Error(),%s)\n\treturn\n}\n", code)
}

func CreateRestServiceInterface(obj *ApiObject) string {
	code := fmt.Sprintf("type %sService interface{\n", obj.Name)
	code += fmt.Sprintf("\tGet(string) (*%s, error)\n", obj.Name)
	code += fmt.Sprintf("\tCreate(...*%s) ([]*%s, error)\n", obj.Name, obj.Name)
	code += fmt.Sprintf("\tUpdate(*%s) (*%s, error)\n", obj.Name, obj.Name)
	code += fmt.Sprintf("\tDelete(string) (*%s, error)\n", obj.Name)
	code += fmt.Sprintf("\tList() ([]*%s, error)\n", obj.Name)
	code += fmt.Sprintf("}\n")

	return code
}

func CreateStruct(obj *ApiObject, structMap map[string]bool, withTags bool) string {
	structMap[obj.Name] = true
	totalCode := ""
	code := fmt.Sprintf("type %s struct{\n", obj.Name)

	for _, field := range obj.Fields {
		switch field.Type {
		case "int":
			code += fmt.Sprintf("\t%s ", field.Name)
			if field.Repeated {
				code += "[]"
			}
			if field.Repeated {
				code += "*"
			}
			code += "int"
			if withTags {
				tagName := stringy.New(field.Name).SnakeCase("?", "").ToLower()
				code += fmt.Sprintf("  `json:\"%s,omitempty\"`\n", tagName)
			} else {
				code += "\n"
			}

		case "string":
			code += fmt.Sprintf("\t%s ", field.Name)
			if field.Repeated {
				code += "[]"
			}
			if field.Repeated {
				code += "*"
			}
			code += "string"
			if withTags {
				tagName := stringy.New(field.Name).SnakeCase("?", "").ToLower()
				code += fmt.Sprintf("  `json:\"%s,omitempty\"`\n", tagName)
			} else {
				code += "\n"
			}

		case "bool":
			code += fmt.Sprintf("\t%s ", field.Name)
			if field.Repeated {
				code += "[]"
			}
			if field.Repeated {
				code += "*"
			}
			code += "bool"
			if withTags {
				tagName := stringy.New(field.Name).SnakeCase("?", "").ToLower()
				code += fmt.Sprintf("  `json:\"%s,omitempty\"`\n", tagName)
			} else {
				code += "\n"
			}
		case "interface":
			code += fmt.Sprintf("\t%s %s\n", field.Object.Name, field.Name)

		default:
			if !structMap[field.Object.Name] {
				totalCode += CreateStruct(field.Object, structMap, withTags)
			}

			if field.Complex {
				code += fmt.Sprintf("\t%s ", field.Name)
				if field.Repeated {
					code += "[]"
				}
				if field.Pointer {
					code += "*"
				}
				code += fmt.Sprintf("%s", field.Object.Name)
				if withTags {
					tagName := stringy.New(field.Name).SnakeCase("?", "").ToLower()
					code += fmt.Sprintf("  `json:\"%s,omitempty\"`\n", tagName)
				} else {
					code += "\n"
				}
			}
		}
	}
	code += "}\n"

	totalCode += code
	return totalCode
}

//"type RestInterface interface {\n\tGet(string2 string) (*Npc,error) \n\tCreate(obj ...*Npc) ([]*Npc,error)\n\tUpdate(obj *Npc) (*Npc,error)\n\tDelete(string2 string) (*Npc,error)\n\tList() ([]*Npc,error)\n}"
