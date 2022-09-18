package apigen

import (
	"fmt"
	"strings"
)

type FieldProperties struct {
	Name     string
	DbName   string
	Type     string
	Repeated bool
	Pointer  bool
	Complex  bool
	Object   *ApiObject
}

type ApiObject struct {
	Name   string
	DbName string
	Fields []*FieldProperties
}

func CreateService(obj *ApiObject) string {

	code := fmt.Sprintf("package %s\n\n", strings.ToLower(obj.Name))
	code += "import (\n\t\"encoding/json\"\n\t\"fmt\"\n\t\"io/ioutil\"\n\t\"net/http\"\n)\n\n"
	code += CreateStruct(obj, structMp, true)
	code += CreateJsonMarshaller(obj)
	code += DeleteFromArrayGen(obj)
	serviceName := obj.GetLocalServiceName()
	serviceStruct := &ApiObject{
		Name: serviceName,
		Fields: []*FieldProperties{
			{
				Name:     obj.GetLocalStoreName(),
				Complex:  true,
				Object:   obj,
				Pointer:  true,
				Repeated: true,
			},
			{
				Name: "Store",
				Type: "string",
			},
		},
	}
	code += CreateRestServiceInterface(obj)

	code += CreateStruct(serviceStruct, structMp, false)
	code += GetGen(serviceStruct, obj)
	code += CreateGen(serviceStruct, obj)
	code += UpdateGen(serviceStruct, obj)
	code += DeleteGen(serviceStruct, obj)
	code += ListGen(serviceStruct, obj)

	code += CreateHttpService(serviceStruct, obj)

	return code
}
