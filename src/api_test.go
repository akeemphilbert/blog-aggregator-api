package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	api "github.com/wepala/blog-aggregator-api/src"
	blogaggregatormodule "github.com/wepala/blog-aggregator-module"
	"github.com/wepala/weos"
)

func TestBlogAdd(t *testing.T) {
	e := echo.New()
	dispatcher := &DispatcherMock{
		DispatchFunc: func(ctx context.Context, command *weos.Command) error {
			return nil
		},
	}
	application := &ApplicationMock{
		DispatcherFunc: func() weos.Dispatcher {
			return dispatcher
		},
	}
	blogAPI := &api.API{
		Application: application,
	}

	request := &blogaggregatormodule.AddBlogRequest{
		Url: "https://ak33m.com",
	}

	reqBytes, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("error setting up request %s",err)
	}
	body := bytes.NewReader(reqBytes)
	req := httptest.NewRequest("PUT","/body",body)
	req = req.WithContext(context.TODO())
	req.Close = true
	recorder := httptest.NewRecorder()
	blogAPI.AddBlog(e.NewContext(req,recorder))

	if len(dispatcher.DispatchCalls()) == 0 {
		t.Error("expected a command to be dispatched")
	}
}