package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
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
		t.Fatalf("error setting up request %s", err)
	}
	body := bytes.NewReader(reqBytes)
	req := httptest.NewRequest("PUT", "/body", body)
	req = req.WithContext(context.TODO())
	req.Close = true
	recorder := httptest.NewRecorder()
	blogAPI.AddBlog(e.NewContext(req, recorder))

	if len(dispatcher.DispatchCalls()) == 0 {
		t.Error("expected a command to be dispatched")
	}
}

func TestGetPosts(t *testing.T) {
	e := echo.New()

	mockPosts := []*api.Post{
		{
			Title: "Post 1",
		},
	}

	mockPage := 1
	mockLimit := 5
	mockBlogId := "abcdef"
	mockCategory := "testing"
	mockStartDate := "07/10/21"
	mockEndDate := "06/10/21"
	var mockPostsResult []*api.Post

	mockProjection := &ProjectionMock{
		GetPostsFunc: func(page, limit int, query string, sortOptions map[string]string, filterOptions map[string]interface{}) ([]*api.Post, int64, error) {
			var filterOption interface{}
			var sortOption string
			var ok bool

			if page != mockPage {
				t.Fatalf("expected page to be %d, got %d", mockPage, page)
			}

			if limit != mockLimit {
				t.Fatalf("expected limit to be %d, got %d", mockLimit, limit)
			}

			//check filter options are set correctly
			if filterOption, ok = filterOptions["blog_id"]; !ok {
				t.Fatal("expected the filter option 'blog_id' to be set")
			}
			if blogId, ok := filterOption.(string); ok {
				if blogId != mockBlogId {
					t.Errorf("expected the blog_id filter value to be '%s', got '%s'", mockBlogId, blogId)
				}
			}

			if filterOption, ok = filterOptions["category"]; !ok {
				t.Fatalf("expected the filter option 'category' to be set")
			}
			if category, ok := filterOption.(string); ok {
				if category != mockCategory {
					t.Errorf("expected the category filter value to be '%s', got '%s'", mockCategory, category)
				}
			}

			if filterOption, ok = filterOptions["start_date"]; !ok {
				t.Fatalf("expected the filter option 'start_date' to be set")
			}
			if startDate, ok := filterOption.(string); ok {
				if startDate != mockStartDate {
					t.Errorf("expected the start_date filter value to be '%s', got '%s'", mockStartDate, startDate)
				}
			}

			if filterOption, ok = filterOptions["end_date"]; !ok {
				t.Fatalf("expected the filter option 'end_date' to be set")
			}
			if endDate, ok := filterOption.(string); ok {
				if endDate != mockEndDate {
					t.Errorf("expected the start_date filter value to be '%s', got '%s'", mockEndDate, endDate)
				}
			}
			//check sort options
			if sortOption, ok = sortOptions["views"]; !ok {
				t.Fatalf("expected a sort option view to be set")
			}

			if sortOption != "desc" {
				t.Errorf("expected the view sort to be '%s', got '%s'", "desc", sortOption)
			}

			mockPostsResult = mockPosts[(page-1)*limit : api.Min(limit*page, len(mockPosts))]
			return mockPostsResult, int64(len(mockPosts)), nil
		},
	}

	application := &ApplicationMock{
		ProjectionsFunc: func() []weos.Projection {
			return []weos.Projection{mockProjection}
		},
	}
	blogAPI := &api.API{
		Application: application,
	}
	req := httptest.NewRequest("GET", fmt.Sprintf("/posts?page=%d&limit=%d&blog_id=%s&category=%s&start_date=%s&end_date=%s&views=desc", mockPage, mockLimit, mockBlogId, mockCategory, mockStartDate, mockEndDate), nil)
	req = req.WithContext(context.TODO())
	req.Close = true
	recorder := httptest.NewRecorder()
	blogAPI.GetPosts(e.NewContext(req, recorder))

	if len(mockProjection.GetPostsCalls()) == 0 {
		t.Error("expected GetPosts to be called")
	}
	//check response code
	if recorder.Code != 200 {
		t.Errorf("expected response code to be %d, got %d", 200, recorder.Code)
	}
	//check response body is a postlist
	var postList *api.PostList
	json.NewDecoder(recorder.Body).Decode(&postList)
	if postList == nil {
		t.Fatal("expected post list response")
	}

	if postList.Total != int64(len(mockPosts)) {
		t.Errorf("expected the total posts to be %d, got %d", len(mockPosts), postList.Total)
	}

	if postList.Page != mockPage {
		t.Errorf("expected the page to be %d, got %d", mockPage, postList.Page)
	}

	if postList.Limit != mockLimit {
		t.Errorf("expected the limit to be %d, got %d", mockLimit, postList.Limit)
	}

	if len(postList.Items) != len(mockPostsResult) {
		t.Fatalf("expected %d items to be returned, got %d", len(mockPostsResult), len(postList.Items))
	}

	for i, _ := range mockPostsResult {
		if postList.Items[i].Title != mockPostsResult[i].Title {
			t.Errorf("expected post in position %d title to be %s, got '%s'", i, mockPostsResult[i].Title, postList.Items[i].Title)
		}
	}

}

func TestGetBlogsError(t *testing.T) {
	e := echo.New()

	mockPage := 1
	mockLimit := 5
	mockError := "WHERE conditions required"

	mockProjection := &ProjectionMock{
		GetPostsFunc: func(page, limit int, query string, sortOptions map[string]string, filterOptions map[string]interface{}) ([]*api.Post, int64, error) {
			if page != mockPage {
				t.Fatalf("expected page to be %d, got %d", mockPage, page)
			}

			if limit != mockLimit {
				t.Fatalf("expected limit to be %d, got %d", mockLimit, limit)
			}

			return nil, 0, errors.New(mockError)
		},
	}

	application := &ApplicationMock{
		ProjectionsFunc: func() []weos.Projection {
			return []weos.Projection{mockProjection}
		},
	}
	blogAPI := &api.API{
		Application: application,
	}
	req := httptest.NewRequest("GET", fmt.Sprintf("/posts?page=%d&limit=%d", mockPage, mockLimit), nil)
	req = req.WithContext(context.TODO())
	req.Close = true
	recorder := httptest.NewRecorder()
	err := blogAPI.GetPosts(e.NewContext(req, recorder))

	if len(mockProjection.GetPostsCalls()) == 0 {
		t.Error("expected GetPosts to be called")
	}

	if err == nil {
		t.Fatalf("expected an error")
	}

	// //check response code
	// if recorder.Code != 500 {
	// 	t.Errorf("expected response code to be %d, got %d",500,recorder.Code)
	// }

	// //confirm error response
	// var errorResponse *api.ErrorResponse
	// json.NewDecoder(recorder.Body).Decode(&errorResponse)

	// if errorResponse == nil {
	// 	t.Fatal("expected error response")
	// }

	// if errorResponse.Message == "" {
	// 	t.Error("an error message must be set")
	// }

	// if errorResponse.Code == "" {
	// 	t.Error("an error code must be set")
	// }

	// if errorResponse.Message == mockError {
	// 	t.Errorf("the raw error should not be returned to the user")
	// }

}

func TestGetCategories(t *testing.T) {
	e := echo.New()

	mockCategories := []*api.Category{
		{
			Title: "Category 1",
		},
	}

	mockPage := 1
	mockLimit := 5
	var mockCategoriesResult []*api.Category

	mockProjection := &ProjectionMock{
		GetCategoriesFunc: func(page, limit int, sortOptions map[string]string, filterOptions map[string]interface{}) ([]*api.Category, int64, error) {
			var sortOption string
			var ok bool

			if page != mockPage {
				t.Fatalf("expected page to be %d, got %d", mockPage, page)
			}

			if limit != mockLimit {
				t.Fatalf("expected limit to be %d, got %d", mockLimit, limit)
			}
			//check sort options
			if sortOption, ok = sortOptions["views"]; !ok {
				t.Fatalf("expected a sort option view to be set")
			}

			if sortOption != "desc" {
				t.Errorf("expected the view sort to be '%s', got '%s'", "desc", sortOption)
			}

			mockCategoriesResult = mockCategories[(page-1)*limit : api.Min(limit*page, len(mockCategories))]
			return mockCategoriesResult, int64(len(mockCategories)), nil
		},
	}

	application := &ApplicationMock{
		ProjectionsFunc: func() []weos.Projection {
			return []weos.Projection{mockProjection}
		},
	}
	blogAPI := &api.API{
		Application: application,
	}
	req := httptest.NewRequest("GET", fmt.Sprintf("/categories?page=%d&limit=%d&views=desc", mockPage, mockLimit), nil)
	req = req.WithContext(context.TODO())
	req.Close = true
	recorder := httptest.NewRecorder()
	blogAPI.GetCategories(e.NewContext(req, recorder))

	if len(mockProjection.GetCategoriesCalls()) == 0 {
		t.Error("expected GetPosts to be called")
	}
	//check response code
	if recorder.Code != 200 {
		t.Errorf("expected response code to be %d, got %d", 200, recorder.Code)
	}
	//check response body is a postlist
	var categoryList *api.CategoryList
	json.NewDecoder(recorder.Body).Decode(&categoryList)
	if categoryList == nil {
		t.Fatal("expected post list response")
	}

	if categoryList.Total != int64(len(mockCategories)) {
		t.Errorf("expected the total posts to be %d, got %d", len(mockCategories), categoryList.Total)
	}

	if categoryList.Page != mockPage {
		t.Errorf("expected the page to be %d, got %d", mockPage, categoryList.Page)
	}

	if categoryList.Limit != mockLimit {
		t.Errorf("expected the limit to be %d, got %d", mockLimit, categoryList.Limit)
	}

	if len(categoryList.Items) != len(mockCategoriesResult) {
		t.Fatalf("expected %d items to be returned, got %d", len(mockCategoriesResult), len(categoryList.Items))
	}

	for i, _ := range mockCategoriesResult {
		if categoryList.Items[i].Title != mockCategoriesResult[i].Title {
			t.Errorf("expected category in position %d title to be %s, got '%s'", i, mockCategoriesResult[i].Title, categoryList.Items[i].Title)
		}
	}

}
