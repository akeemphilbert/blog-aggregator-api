package api

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	blogaggregatormodule "github.com/wepala/blog-aggregator-module"
	"github.com/wepala/weos"
	weoscontroller "github.com/wepala/weos-controller"
)

type API struct {
	weoscontroller.API
	application weos.Application
	Log weos.Log
	DB *sql.DB
}

func (a *API) AddBlog(e echo.Context) error {
	e.Echo().Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"*"},
		AllowMethods: []string{http.MethodPut},
	}))
	return e.JSON(http.StatusCreated, "Blog Added")
}

func (a *API) Initialize() error {
	var err error
	//initialize app
	a.application, err = weos.NewApplicationFromConfig(a.Config.ApplicationConfig, a.Log, a.DB, nil, nil)
	if err != nil {
		return err
	}
	//enable module
	err = blogaggregatormodule.Initialize(a.application)
	if err != nil {
		return err
	}
	//run fixtures
	err = a.application.Migrate(context.Background())
	if err != nil {
		return err
	}
	//set log level to debug
	a.EchoInstance().Logger.SetLevel(log.DEBUG)
	return nil
}

func New(port *string, apiConfig string) {
	e := echo.New()
	weoscontroller.Initialize(e,&API{},apiConfig)
	e.Logger.Fatal(e.Start(":"+*port))
}