package main

import (
	"electric-bicycle/controllers"
	"electric-bicycle/models"
	"github.com/kataras/iris/v12"
	"log"
)

func main() {

	//log.Print(reflect.TypeOf(time.Now()))

	engine, err := models.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	defer engine.Close()

	app := iris.Default()
	app.Post("/distance/add", controllers.DistanceNew)
	app.Post("/distance/delete", controllers.DistanceDelete)
	app.Post("/distance/update", controllers.DistanceUpdate)
	app.Get("/distance/list", controllers.DistanceList)
	app.Run(iris.Addr(":8080"))
}
