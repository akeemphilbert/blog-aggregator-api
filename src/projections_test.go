package api_test

import (
	"context"
	"testing"

	api "github.com/wepala/blog-aggregator-api/src"
	"github.com/wepala/weos"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestProjection_GetPosts(t *testing.T) {
	//setup gorm db connection
	//TODO setup a way to test against multiple database
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database '%s'",err)
	}

	logger := &LogMock{
		ErrorFunc: func(args ...interface{}) {},
	}
	
	application := &ApplicationMock{
		DBFunc: func() *gorm.DB {
			return db
		},
		LoggerFunc: func() weos.Log {
			return logger
		},
		AddProjectionFunc: func(projection weos.Projection) error {
			return nil
		},
	}
	
	projection, err := api.NewProjection(application)
	if err != nil {
		t.Fatalf("unexpected error setting up projection '%s'",err)
	}
	//add blogs and blog posts to the database
	projection.Migrate(context.Background())
	//check that the database is called
	if len(application.DBCalls()) == 0 {
		t.Error("expected the db to be called")
	}
	mockBlogs := []*api.Blog{
		{
			ID: "123",
			Title: "Some Blog 1",
		},
		{
			ID: "456",
			Title: "Some Blog 1",
		},
	}
	db.Create(mockBlogs)
	if db.Error != nil {
		t.Fatalf("error setting up mock blogs '%s'",db.Error)
	}
	mockPosts := []*api.Post {
		{
			ID: "1",
			Title: "Post 1",
			BlogID: "123",
		},
		{
			ID: "2",
			Title: "Post 2",
			BlogID: "123",
		},
		{
			ID: "3",
			Title: "Post 3",
			BlogID: "456",
		},
		{
			ID: "4",
			Title: "Post 3",
			BlogID: "123",
		},
		{
			ID: "5",
			Title: "Post 4",
			BlogID: "123",
		},
		{
			ID: "6",
			Title: "Post 5",
			BlogID: "123",
		},
		{
			ID: "7",
			Title: "Post 6",
			BlogID: "123",
		},
		{
			ID: "8",
			Title: "Post 7",
			BlogID: "123",
		},
	}
	db.Create(mockPosts)
	if db.Error != nil {
		t.Fatalf("error setting up mock posts '%s'",db.Error)
	}
	t.Run("get posts by blog", func(t *testing.T) {
		//run get posts
		filters := make(map[string]interface{})
		filters["blog_id"] = "123"
		posts, count, err := projection.GetPosts(2,2,"",nil,filters)
		if err != nil {
			t.Fatalf("unexpected error getting posts '%s'",err)
		}
		if count != 7 {
			t.Errorf("expected the number posts to be returned to be %d, got %d",7,count)
		}

		if len(posts) != 2 {
			t.Fatalf("expected %d posts to be returned, got %d",2,len(posts))
		}

		//check that the first result matches the item in the list having accounted for pagination
		if posts[0].Title != mockPosts[3].Title {
			t.Errorf("expected the post in position %d to have title %s, got '%s'",0,mockPosts[3].Title,posts[0].Title)
		}
	})
	

	
}