//go:generate moq -pkg api_test -out projectionmock_test.go . Projection
package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	blogaggregatormodule "github.com/wepala/blog-aggregator-module"
	"github.com/wepala/weos"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Projection interface {
	weos.Projection
	GetBlogByID (id string) (*Blog, error)
	GetBlogByURL(url string) (*Blog, error)
	GetPosts (page int, limit int, query string, sortOptions map[string]string, filterOptions map[string]interface{}) ([]*Post, int64, error)
}

type Blog struct {
	gorm.Model
	ID string `gorm:"primarykey"`
	Title string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	URL string `json:"url,omitempty"`
	FeedURL string`json:"feedUrl,omitempty"`
	Authors []*Author `json:"authors,omitempty"`
	Posts []*Post `json:"posts,omitempty"`
}

type Author struct {
	gorm.Model
	Name string
	Email string
	BlogID string `json:"blogId"`
}

type Post struct {
	gorm.Model
	ID string `gorm:"primarykey"`
	Title string
	Description string
	Content string
	BlogID string `json:"blogId"`
	Categories []*Category `json:"categories,omitempty" gorm:"many2many:post_categories;"`
	PublishDate time.Time
	Views int
}

type Category struct {
	gorm.Model
	Title string
	Description string
	Posts []*Post `json:"posts,omitempty" gorm:"many2many:post_categories;"`
}

type GORMProjection struct {
	db *gorm.DB
	logger weos.Log
	migrationFolder string
}

func (p *GORMProjection) Persist( entities []weos.Entity) error {
	return nil
}

func (p *GORMProjection) Remove (entities []weos.Entity) error {
	return nil
}

func (p *GORMProjection) GetBlogByID (id string) (*Blog, error) {
	return nil,nil
}

func (p *GORMProjection) GetBlogs () ([]*Blog, error) {
	return nil,nil
}

func (p *GORMProjection) GetBlogByURL(url string) (*Blog, error) {
	var blog *Blog
	if err := p.db.Debug().Preload(clause.Associations).First(&blog,"url = ? OR feed_url = ?",url, url).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil,fmt.Errorf("blog '%s' not found",url)
		}
	}
	if p.db.Error != nil {
		return nil, p.db.Error
	}
	return blog, nil
}
//GetPosts get all the posts in the aggregator
func (p *GORMProjection) GetPosts (page int, limit int, query string, sortOptions map[string]string, filterOptions map[string]interface{}) ([]*Post, int64, error) {
	var posts []*Post
	var count int64
	result := p.db.Debug().Preload("Categories").Scopes(filter(filterOptions),paginate(page,limit),sort(sortOptions)).Find(&posts).Offset(-1).Distinct("posts.id").Count(&count)
	return posts,count,result.Error
}

func sort(order map[string]string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		for key,value := range order {
			//only support certain values since GORM doesn't protect the order function https://gorm.io/docs/security.html#SQL-injection-Methods
			if (value != "asc" && value != "desc" && value != "") || (key != "views" && key != "publishDate") {
				return db
			}
			db.Order(key+" "+value)
		}
		
		return db
	}
}

func category(categoryValue interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if category, ok := categoryValue.(string); ok {
			db.Joins("left join post_categories on post_id = posts.id").Joins("left join categories on categories.id = post_categories.category_id").Where("categories.title = ?",category)
		}
		return db
	}
}

func publishDate(startDateValue interface{}, endDateValue interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if startDateValue != nil && endDateValue != nil {
			var startDate time.Time
			var endDate time.Time
			var err error
			if sdv, ok := startDateValue.(string); ok {
				startDate, err = time.Parse("01/02/06",sdv)
				if err != nil {
					return db
				}
			}
	
			if edv, ok := endDateValue.(string); ok {
				endDate, err = time.Parse("01/02/06",edv)
				endDate = time.Date(endDate.Year(),endDate.Month(),endDate.Day(),23,59,59,endDate.Nanosecond(),endDate.Location())
				if err != nil {
					return db
				}
			}
	
			db.Where("publish_date BETWEEN ? AND ?", startDate, endDate)
		}
		
		return db
	}
}

func filter(filter map[string]interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if filter != nil {
			if categoryValue,ok := filter["category"];ok {
				db.Scopes(category(categoryValue))
				delete(filter,"category")
			}

			var startDateValue interface{}
			var endDateValue interface{}
			var ok bool

			if startDateValue,ok = filter["start_date"];ok {
				delete(filter,"start_date")
			}

			if endDateValue,ok = filter["end_date"];ok {
				delete(filter,"end_date")
			}
			if startDateValue != nil && endDateValue != nil {
				db.Scopes(publishDate(startDateValue,endDateValue))
			}
			return db.Where(filter)
		}
		return db
	}
}

func paginate(page int, limit int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		actualLimit := limit
		actualPage := page
		if actualLimit == 0 {
			actualLimit = -1
		}
		if actualPage == 0 {
			actualPage = 1
		}
		return db.Offset((page - 1) * limit).Limit(actualLimit)
	}
}


func (p *GORMProjection) GetEventHandler() weos.EventHandler {
	return func (event weos.Event) {
		switch event.Type {
		case blogaggregatormodule.BLOG_ADDED:
			var blog *Blog
			err := json.Unmarshal(event.Payload,&blog)
			if err != nil {
				p.logger.Errorf("error unmarshalling event '%s'",err)
			}
			db := p.db.Create(blog)
			if db.Error != nil {
				p.logger.Errorf("error creating blog '%s'", err)
			}
		case blogaggregatormodule.BLOG_UPDATED:
			var blog *Blog
			err := json.Unmarshal(event.Payload,&blog)
			if err != nil {
				p.logger.Errorf("error unmarshalling event '%s'",err)
			}
			blog.ID = event.Meta.EntityID
			db := p.db.Model(blog).Updates(blog)
			if db.Error != nil {
				p.logger.Errorf("error updating blog '%s'", err)
			}
		case blogaggregatormodule.AUTHOR_CREATED:
			var author *Author
			err := json.Unmarshal(event.Payload,&author)
			if err != nil {
				p.logger.Errorf("error unmarhsalling event '%s'",err)
			}
			db := p.db.Create(author)
			if db.Error != nil {
				p.logger.Errorf("error creating author '%s'",err)
			}
		case blogaggregatormodule.POST_CREATED:
			var post *Post
			err := json.Unmarshal(event.Payload,&post)
			if err != nil {
				p.logger.Errorf("error unmarshalling event '%s'",err)
			}
			db := p.db.Create(post)
			if db.Error != nil {
				p.logger.Errorf("error creating post '%s'",err)
			}
		}
	}
}
//runs migrations
func (p *GORMProjection) Migrate(ctx context.Context) error {
	err := p.db.AutoMigrate(&Blog{},&Post{},&Author{},&Category{})
	if err != nil {
		return err
	}

	return nil
}

func NewProjection(application weos.Application) (*GORMProjection, error) {
	projection := &GORMProjection{
		db: application.DB(),
		logger: application.Logger(),
	}
	application.AddProjection(projection)
	return projection, nil
}

