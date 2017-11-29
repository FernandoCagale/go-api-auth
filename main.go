package main

import (
	"time"

	"github.com/FernandoCagale/go-api-auth/src/checker"
	"github.com/FernandoCagale/go-api-auth/src/config"
	"github.com/FernandoCagale/go-api-auth/src/datastore"
	"github.com/FernandoCagale/go-api-auth/src/handlers"
	"github.com/FernandoCagale/go-api-auth/src/lib"
	"github.com/jinzhu/gorm"

	log "github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	var db *gorm.DB
	env := config.LoadEnv()
	app := echo.New()

	go bindDatastore(app, db, env.DatastoreURL)

	app.Use(middleware.Logger())

	defer db.Close()

	checkers := map[string]checker.Checker{
		"api":      checker.NewApi(),
		"postgres": checker.NewPostgres(env.DatastoreURL),
	}

	healthzHandler := handlers.NewHealthzHandler(checkers)
	app.GET("/health", healthzHandler.HealthzIndex)

	userHandler := handlers.NewUserHandler()

	app.POST("/login", userHandler.Login)
	app.GET("/user", userHandler.GetAll)
	app.POST("/user", userHandler.Save)

	group := app.Group("/v1")

	group.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningMethod: "HS512",
		SigningKey:    []byte("secret"),
	}))

	group.GET("/user/:id", userHandler.Get)
	group.PUT("/user/:id", userHandler.Update)
	group.DELETE("/user/:id", userHandler.Delete)

	app.Logger.Fatal(app.Start(":" + env.Port))
}

func bindDatastore(app *echo.Echo, db *gorm.DB, url string) {
	for {
		db, err := datastore.New(url)
		failOnError(err, "Failed to init dababase connection!")
		if err == nil {
			app.Use(lib.BindDb(db))
			break
		}
		time.Sleep(time.Second * 5)
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Info(msg)
	}
}
