package main

import (
	"apigen"
	npc "apigen/cmd/service"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {

	//ExampleNpcServiceRun()
	NpcExampleServiceGen()
}

func NpcExampleServiceGen() {
	reviewFields := []*apigen.FieldProperties{
		{
			Name: "Stars",
			Type: "int",
		},
		{
			Name: "Review",
			Type: "string",
		},
	}
	reviewObj := &apigen.ApiObject{Name: "Review", Fields: reviewFields}

	fields := []*apigen.FieldProperties{
		{
			Name: "Points",
			Type: "int",
		},
		{
			Name: "Price",
			Type: "string",
		},
		{
			Name:     "Games",
			Complex:  true,
			Object:   reviewObj,
			Repeated: true,
		},
	}
	gameObj := &apigen.ApiObject{Name: "Game", Fields: fields}

	fields = []*apigen.FieldProperties{
		{
			Name: "Money",
			Type: "int",
		},
		{
			Name: "Id",
			Type: "string",
		},
		{
			Name: "Name",
			Type: "string",
		},
		{
			Name: "Rizz",
			Type: "string",
		},

		{
			Name:     "Bots",
			Type:     "string",
			Repeated: true,
		},

		{
			Name:     "Games",
			Complex:  true,
			Object:   gameObj,
			Repeated: true,
		},
	}
	npcObj := &apigen.ApiObject{Name: "Npc", Fields: fields}

	npcObj.GetLocalStoreName()

	GenerateCode(npcObj)

}

func GenerateCode(obj *apigen.ApiObject) {
	generatedPath := filepath.Join(".", "cmd", "service")
	err := os.MkdirAll(generatedPath, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	f, err := os.Create(filepath.Join(generatedPath, "service.go"))
	if err != nil {
		log.Fatal(err)
	}
	data := []byte(apigen.CreateService(obj))
	f.Write(data)
}

func ExampleNpcServiceRun() {
	npcService := &npc.NpcLocalService{}
	_, err := npcService.Create(&npc.Npc{
		Money: 10,
		Id:    "1",
		Name:  "npc1",
		Rizz:  "none",
		Bots:  nil,
		Games: nil,
	})
	if err != nil {
		log.Fatal(err)
	}

	npcHttpService := npc.NpcHttpService{
		NpcService:   npcService,
		EndpointName: "/npc",
	}

	r := mux.NewRouter().StrictSlash(false)
	npcHttpService.AddSubrouter(r)

	fmt.Println(r)

	http.Handle("/", r)

	http.ListenAndServe(":8000", nil)
}
