package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	blogaggregatormodule "github.com/wepala/blog-aggregator-module"
	"github.com/wepala/weos"
	weoscontroller "github.com/wepala/weos-controller"
)

type API struct {
	weoscontroller.API
	Application weos.Application
	Log weos.Log
	DB *sql.DB
	Client *http.Client
	projection *GORMProjection
}

func (a *API) AddBlog(e echo.Context) error {
	e.Echo().Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"*"},
		AllowMethods: []string{http.MethodPut},
	}))
	var blogAddRequest *blogaggregatormodule.AddBlogRequest
	err := json.NewDecoder(e.Request().Body).Decode(&blogAddRequest)
	if err != nil {
		return err
	}
	a.Application.Dispatcher().Dispatch(e.Request().Context(),blogaggregatormodule.AddBlogCommand(blogAddRequest.Url))
	return e.JSON(http.StatusCreated, "Blog Added")
}

func (a *API) Initialize() error {
	var err error
	//initialize app
	if a.Client == nil {
		a.Client = &http.Client{
			Timeout: time.Second*10,
		}
	}
	a.Application, err = weos.NewApplicationFromConfig(a.Config.ApplicationConfig, a.Log, a.DB, a.Client, nil)
	if err != nil {
		return err
	}
	//setup projections
	a.projection, err = NewProjection(a.Application)
	if err != nil {
		return err
	}
	//enable module
	err = blogaggregatormodule.Initialize(a.Application)
	if err != nil {
		return err
	}
	//run fixtures
	err = a.Application.Migrate(context.Background())
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