package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	blogaggregatormodule "github.com/wepala/blog-aggregator-module"
	"github.com/wepala/weos"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Blog struct {
	gorm.Model
	ID          string    `gorm:"primarykey"`
	Title       string    `json:"title,omitempty"`
	Description string    `json:"description,omitempty"`
	URL         string    `json:"url,omitempty"`
	FeedURL     string    `json:"feedUrl,omitempty"`
	Authors     []*Author `json:"authors,omitempty"`
	Posts       []*Post   `json:"posts,omitempty"`
}

type Author struct {
	gorm.Model
	Name   string
	Email  string
	BlogID string `json:"blogId"`
}

type Post struct {
	gorm.Model
	Title       string
	Description string
	Content     string
	BlogID      string `json:"blogId"`
}

type GORMProjection struct {
	db              *gorm.DB
	logger          weos.Log
	migrationFolder string
}

func (p *GORMProjection) Persist(entities []weos.Entity) error {
	return nil
}

func (p *GORMProjection) Remove(entities []weos.Entity) error {
	return nil
}

func (p *GORMProjection) GetBlogByID(id string) (*Blog, error) {
	var blog *Blog
	if err := p.db.Debug().Preload(clause.Associations).First(&blog, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("blog '%s' not found", id)
		}
	}
	if p.db.Error != nil {
		return nil, p.db.Error
	}

	return blog, nil
}

func (p *GORMProjection) GetBlogs() ([]*Blog, error) {
	return nil, nil
}

func (p *GORMProjection) GetBlogByURL(url string) (*Blog, error) {
	var blog *Blog
	if err := p.db.Debug().Preload(clause.Associations).First(&blog, "url = ? OR feed_url = ?", url, url).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("blog '%s' not found", url)
		}
	}
	if p.db.Error != nil {
		return nil, p.db.Error
	}
	return blog, nil
}

func (p *GORMProjection) GetPostsByCategory(id string) ([]*Post, error) {
	return nil, nil
}

func (p *GORMProjection) GetPostsByBlog(id string) ([]*Post, error) {
	return nil, nil
}

func (p *GORMProjection) GetPostsByAuthor(id string) ([]*Post, error) {
	return nil, nil
}

//GetPosts get all the posts in the aggregator
func (p *GORMProjection) GetPosts() ([]*Post, error) {
	return nil, nil
}

func (p *GORMProjection) GetEventHandler() weos.EventHandler {
	return func(event weos.Event) {
		switch event.Type {
		case blogaggregatormodule.BLOG_ADDED:
			var blog *Blog
			err := json.Unmarshal(event.Payload, &blog)
			if err != nil {
				p.logger.Errorf("error unmarshalling event '%s'", err)
			}
			db := p.db.Create(blog)
			if db.Error != nil {
				p.logger.Errorf("error creating blog '%s'", err)
			}
		case blogaggregatormodule.BLOG_UPDATED:
			var blog *Blog
			err := json.Unmarshal(event.Payload, &blog)
			if err != nil {
				p.logger.Errorf("error unmarshalling event '%s'", err)
			}
			blog.ID = event.Meta.EntityID
			db := p.db.Model(blog).Updates(blog)
			if db.Error != nil {
				p.logger.Errorf("error updating blog '%s'", err)
			}
		case blogaggregatormodule.AUTHOR_CREATED:
			var author *Author
			err := json.Unmarshal(event.Payload, &author)
			if err != nil {
				p.logger.Errorf("error unmarhsalling event '%s'", err)
			}
			db := p.db.Create(author)
			if db.Error != nil {
				p.logger.Errorf("error creating author '%s'", err)
			}
		case blogaggregatormodule.POST_CREATED:
			var post *Post
			err := json.Unmarshal(event.Payload, &post)
			if err != nil {
				p.logger.Errorf("error unmarshalling event '%s'", err)
			}
			db := p.db.Create(post)
			if db.Error != nil {
				p.logger.Errorf("error creating post '%s'", err)
			}
		}
	}
}

//runs migrations
func (p *GORMProjection) Migrate(ctx context.Context) error {
	err := p.db.AutoMigrate(&Blog{}, &Post{}, &Author{})
	if err != nil {
		return err
	}

	return nil
}

func NewProjection(application weos.Application) (*GORMProjection, error) {
	projection := &GORMProjection{
		db:     application.DB(),
		logger: application.Logger(),
	}
	application.AddProjection(projection)
	return projection, nil
}
