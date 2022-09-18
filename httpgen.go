package apigen

import "fmt"

func CreateHttpService(serviceObj *ApiObject, obj *ApiObject) string {
	code := ""
	httpServiceStruct := &ApiObject{
		Name: obj.GetHttpServiceName(),
		Fields: []*FieldProperties{
			{
				Name: obj.GetServiceInterfaceName(),
				Type: "interface",
				Object: &ApiObject{
					Name: obj.GetServiceInterfaceName(),
				},
			},
			{
				Name: "Store",
				Type: "string",
			},
			{
				Name: "EndpointName",
				Type: "string",
			},
		},
	}
	code += CreateStruct(httpServiceStruct, structMp, false)
	code += HttpGetGen(httpServiceStruct, serviceObj, obj)
	code += HttpCreateGen(httpServiceStruct, serviceObj, obj)
	code += HttpUpdateGen(httpServiceStruct, serviceObj, obj)
	code += HttpDeleteGen(httpServiceStruct, serviceObj, obj)
	code += HttpListGen(httpServiceStruct, serviceObj, obj)
	code += AddSubrouterGen(httpServiceStruct, serviceObj, obj)

	//funcName := fmt.Sprintf("CreateNew%s",serviceObj.GetHttpServiceName())
	//
	//createHttpServiceFunc := FunctionGen("","",funcName,nil,[]string{"[]*"+serviceObj.GetHttpServiceName()},content)

	return code
}

func HttpGetGen(httpService *ApiObject, service *ApiObject, obj *ApiObject) string {
	content := fmt.Sprintf("data, err := ioutil.ReadAll(r.Body)\n")
	content += HandleHttpErrorGen("400")
	content += fmt.Sprintf("requestData := make(map[string]string, 0)\n")
	content += fmt.Sprintf("err = json.Unmarshal(data, &requestData)\n")
	content += HandleHttpErrorGen("400")
	content += fmt.Sprintf("res,err := s.%s.Get(requestData[\"id\"])\n", obj.GetServiceInterfaceName())
	content += HandleHttpErrorGen("400")
	content += fmt.Sprintf("resData,err := res.Marshal()\n")
	content += HandleHttpErrorGen("500")
	content += fmt.Sprintf("fmt.Fprint(w,string(resData))\n")

	code := FunctionGen(
		httpService.Name,
		"s",
		"GetHandler",
		[]string{"w http.ResponseWriter", " r *http.Request"},
		[]string{},
		content,
	)
	return code
}

func HttpCreateGen(httpService *ApiObject, service *ApiObject, obj *ApiObject) string {
	content := fmt.Sprintf("data, err := ioutil.ReadAll(r.Body)\n")
	content += HandleHttpErrorGen("400")
	content += fmt.Sprintf("requestData := make(map[string][]*%s, 0)\n", obj.Name)
	content += fmt.Sprintf("err = json.Unmarshal(data, &requestData)\n")
	content += HandleHttpErrorGen("400")
	content += fmt.Sprintf("res,err := s.%s.Create(requestData[\"data\"]...)\n", obj.GetServiceInterfaceName())
	content += HandleHttpErrorGen("400")
	content += fmt.Sprintf("resData,err := json.Marshal(res)\n")
	content += HandleHttpErrorGen("500")
	content += fmt.Sprintf("fmt.Fprint(w,string(resData))\n")

	code := FunctionGen(
		httpService.Name,
		"s",
		"CreateHandler",
		[]string{"w http.ResponseWriter", " r *http.Request"},
		[]string{},
		content,
	)
	return code
}

func HttpUpdateGen(httpService *ApiObject, service *ApiObject, obj *ApiObject) string {
	content := fmt.Sprintf("data, err := ioutil.ReadAll(r.Body)\n")
	content += HandleHttpErrorGen("400")
	content += fmt.Sprintf("requestData := make(map[string]*%s, 0)\n", obj.Name)
	content += fmt.Sprintf("err = json.Unmarshal(data, &requestData)\n")
	content += HandleHttpErrorGen("400")
	content += fmt.Sprintf("res,err := s.%s.Update(requestData[\"data\"])\n", obj.GetServiceInterfaceName())
	content += HandleHttpErrorGen("400")
	content += fmt.Sprintf("resData,err := res.Marshal()\n")
	content += HandleHttpErrorGen("500")
	content += fmt.Sprintf("fmt.Fprint(w,string(resData))\n")

	code := FunctionGen(
		httpService.Name,
		"s",
		"UpdateHandler",
		[]string{"w http.ResponseWriter", " r *http.Request"},
		[]string{},
		content,
	)
	return code
}

func HttpDeleteGen(httpService *ApiObject, service *ApiObject, obj *ApiObject) string {
	content := fmt.Sprintf("data, err := ioutil.ReadAll(r.Body)\n")
	content += HandleHttpErrorGen("400")
	content += fmt.Sprintf("requestData := make(map[string]string, 0)\n")
	content += fmt.Sprintf("err = json.Unmarshal(data, &requestData)\n")
	content += HandleHttpErrorGen("400")
	content += fmt.Sprintf("res,err := s.%s.Delete(requestData[\"id\"])\n", obj.GetServiceInterfaceName())
	content += HandleHttpErrorGen("400")
	content += fmt.Sprintf("resData,err := res.Marshal()\n")
	content += HandleHttpErrorGen("500")
	content += fmt.Sprintf("fmt.Fprint(w,string(resData))\n")

	code := FunctionGen(
		httpService.Name,
		"s",
		"DeleteHandler",
		[]string{"w http.ResponseWriter", " r *http.Request"},
		[]string{},
		content,
	)
	return code
}

func HttpListGen(httpService *ApiObject, service *ApiObject, obj *ApiObject) string {
	content := fmt.Sprintf("res,err := s.%s.List()\n", obj.GetServiceInterfaceName())
	content += HandleHttpErrorGen("400")
	content += fmt.Sprintf("resData,err := json.Marshal(res)\n")
	content += HandleHttpErrorGen("500")
	content += fmt.Sprintf("fmt.Fprint(w,string(resData))\n")

	code := FunctionGen(
		httpService.Name,
		"s",
		"ListHandler",
		[]string{"w http.ResponseWriter", " r *http.Request"},
		[]string{},
		content,
	)
	return code
}

func AddSubrouterGen(httpService *ApiObject, service *ApiObject, obj *ApiObject) string {
	content := "router := r.PathPrefix(s.EndpointName).Subrouter().StrictSlash(false)\nrouter.HandleFunc(\"/\", s.GetHandler).Methods(\"GET\", \"OPTIONS\")\nrouter.HandleFunc(\"/list\", s.ListHandler).Methods(\"GET\", \"OPTIONS\")\nrouter.HandleFunc(\"/\", s.CreateHandler).Methods(\"POST\", \"OPTIONS\")\nrouter.HandleFunc(\"/update\", s.UpdateHandler).Methods(\"POST\", \"OPTIONS\")\nrouter.HandleFunc(\"/delete\", s.DeleteHandler).Methods(\"POST\", \"OPTIONS\")"
	code := FunctionGen(
		httpService.Name,
		"s",
		"AddSubrouter",
		[]string{"r *mux.Router"},
		[]string{},
		content,
	)
	return code
}
