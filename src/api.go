package api

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
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
	blogAddRequest := &blogaggregatormodule.AddBlogRequest{e.FormValue("url")}
	err := a.Application.Dispatcher().Dispatch(e.Request().Context(),blogaggregatormodule.AddBlogCommand(blogAddRequest.Url))
	if err != nil {
		return weoscontroller.NewControllerError("Error creating blog",err,0)
	}
	return e.JSON(http.StatusCreated, "Blog Added")
}
//Get list of authors
func (a *API) GetAuthors(e echo.Context) error {
	page,_ := strconv.Atoi(e.QueryParam("page"))
	limit,_ := strconv.Atoi(e.QueryParam("limit"))

	filters := make(map[string]interface{})
	sorts := make(map[string]string)

	authors,_, err := a.projection.GetAuthors(page,limit,"",sorts,filters)

	if err != nil {
		return weoscontroller.NewControllerError("Error getting authors",err,0)
	}
	return e.JSON(http.StatusOK, authors)
}
//Get list of posts. 
func (a *API) GetPosts (e echo.Context) error {
	//initialize projection params
	var lastError error
	var page int
	var limit int
	filters := make(map[string]interface{})
	sorts := make(map[string]string)
	//parse query parameters
	page, _ = strconv.Atoi(e.QueryParam("page"))
	limit, _ = strconv.Atoi(e.QueryParam("limit"))
	//parse sort parameters
	if viewsSort := e.QueryParam("views"); viewsSort != "" {
		sorts["views"] = viewsSort
	}
	//parse query parameters
	if blogId:=e.QueryParam("blog_id");blogId != "" {
		filters["blog_id"] = blogId
	}

	if category:=e.QueryParam("category");category != "" {
		filters["category"] = category
	}

	startDate := e.QueryParam("start_date")
	endDate := e.QueryParam("end_date")

	if startDate != "" && endDate != "" {
		filters["start_date"] = startDate
		filters["end_date"] = endDate
	}

	if page == 0 {
		page = 1
	}

	for _,projection := range a.Application.Projections() {
		posts, count, err := projection.(Projection).GetPosts(page,limit,"",sorts,filters)
		if err == nil {
			return e.JSON(http.StatusOK,&PostList{
				Page: page,
				Limit: limit,
				Total: count,
				Items: posts,
			})
		} else {
			lastError = err
		}
	}
	return lastError
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