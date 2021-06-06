package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/cucumber/godog"
	"github.com/cucumber/messages-go/v10"
	"github.com/labstack/echo/v4"
	api "github.com/wepala/blog-aggregator-api/src"
	blogaggregatormodule "github.com/wepala/blog-aggregator-module"
	"github.com/wepala/go-testhelpers"
	weoscontroller "github.com/wepala/weos-controller"
)

type TestBlog struct
{
	Title string
	URL string
	FeedLink string
}

type TestUser struct
{
	Name string
	Site string
	IsLoggedIn bool
	Blog *TestBlog
}

type FeedItem struct {
	Title string 
	Description string
	Link string
	Category string
	PublishDate string
}

var testUsers map[string]*TestUser
var testBlogs map[string]*TestBlog
var testBlog *TestBlog
var testBlogPage string
var testFeed string
var err error
var blogAPI *api.API
var request interface{}
var endpoint string //the endpoint for the request
var method string//the method of the request
var e *echo.Echo
var response *http.Response
var createdBlog *api.Blog

func aPingbackUrlShouldBeGenerated() error {
	return godog.ErrPending
}

func aUserNamed(arg1 string) error {
	testUsers[arg1] = &TestUser{
		Name: arg1,
	}
	return err
}

func anAuthorShouldBeCreatedForEachAuthorInTheFeed() error {
	return godog.ErrPending
}

func anErrorScreenShouldBeShown(arg1 string) error {
	return godog.ErrPending
}

func followsTheBlog(arg1, arg2 string) error {
	return nil
}

func hasABlog(arg1, arg2 string) error {
	if user,ok := testUsers[arg1]; ok {
		user.Blog = &TestBlog{
			URL: arg2,
		}
		testBlogs[arg2] = user.Blog
		testBlog = user.Blog
		return err
	}
	err = fmt.Errorf("user %s not defined",arg1)
	return err
}

func hitsTheSubmitButton(arg1 string) error {
	reqBytes, err := json.Marshal(request)
	if err != nil {
		return err
	}
	body := bytes.NewReader(reqBytes)
	req := httptest.NewRequest(method,endpoint,body)
	req = req.WithContext(context.TODO())
	req.Close = true
	rw := httptest.NewRecorder()
	e.ServeHTTP(rw,req)
	response = rw.Result()
	defer response.Body.Close()

	return err
}

func isLoggedIn(arg1 string) error {
	if user,ok := testUsers[arg1]; ok {
		user.IsLoggedIn = true
		return err
	}
	
	err =  fmt.Errorf("user %s not defined",arg1)
	return err
}

func isLoggedInWithGoogle(arg1 string) error {
	if user,ok := testUsers[arg1]; ok {
		user.IsLoggedIn = true
		return err
	}
	
	err =  fmt.Errorf("user %s not defined",arg1)
	return err
}

func isNotLoggedIn(arg1 string) error {
	if user,ok := testUsers[arg1]; ok {
		user.IsLoggedIn = false
		return nil
	}
	
	return fmt.Errorf("user %s not defined",arg1)
}

func isOnTheBlogSubmitScreen(arg1 string) error {
	request = &blogaggregatormodule.AddBlogRequest{}
	method = "PUT"
	endpoint = "/blog"
	return nil
}

func postsShouldBeCreatedForEachPost() error {
	return godog.ErrPending
}

func profilesForTheBlogAuthorsShouldBeCreated() error {
	if createdBlog == nil {
		return fmt.Errorf("blog was not created by a previous step")
	}
	if len(createdBlog.Authors) == 0 {
		return fmt.Errorf("expected there to be authors with blog")
	}
	return err
}

func shouldBeRedirectedToTheProfilePageForThatBlog(arg1 string) error {
	return nil
}

func successfullyCompletesTheCaptcha(arg1 string) error {
	return nil
}

func successfullySubmitsAFeed(arg1 string) error {
	return godog.ErrPending
}

func theAggregatorSupportsAtomFeedsAsWellAsRssFeeds() error {
	return nil
}

func theBlogDetailsStoredInTheAggregator() error {
	return godog.ErrPending
}

func theBlogHasALinkToAFeed(arg1 string) error {
	testBlogPage = fmt.Sprintf(`<!DOCTYPE html><html lang="en" data-theme=""><head><title> Akeem Philbert | Akeem Philbert&#39;s Blog </title><meta charset="utf-8"><meta name="generator" content="Hugo 0.82.0" /><meta name="viewport" content="width=device-width,initial-scale=1,viewport-fit=cover"><meta name="description" content="">
		
		<link rel="stylesheet"
			  href="https://ak33m.com/css/style.min.2277e4d1f5f913138c1883033695f7a9779a2dcdc66ae94d514bd151bebd8f78.css"
			  integrity="sha256-Infk0fX5ExOMGIMDNpX3qXeaLc3GaulNUUvRUb69j3g="
			  crossorigin="anonymous"
			  type="text/css">
		
		<link rel="stylesheet"
			href="https://ak33m.com/css/markupHighlight.min.f798cbda9aaa38f89eb38be6414bd082cfd71a6780375cbf67b6d2fb2b96491e.css"
			integrity="sha256-95jL2pqqOPies4vmQUvQgs/XGmeAN1y/Z7bS&#43;yuWSR4="
			crossorigin="anonymous"
			type="text/css">
		
		<link rel="stylesheet" 
		href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/5.15.1/css/all.min.css" 
		integrity="sha512-+4zCK9k+qNFUR5X+cKL9EIR+ZOhtIloNl9GIKS57V1MyNsYpYcUrUeQc9vNfzsWfV28IaLL3i96P9sdNyeRssA==" 
		crossorigin="anonymous" />
	
		
		<link rel="shortcut icon" href="https://ak33m.com/favicon.ico" type="image/x-icon">
		<link rel="apple-touch-icon" sizes="180x180" href="https://ak33m.com/apple-touch-icon.png">
		<link rel="icon" type="image/png" sizes="32x32" href="https://ak33m.com/favicon-32x32.png">
		<link rel="icon" type="image/png" sizes="16x16" href="https://ak33m.com/favicon-16x16.png">
	
		<link rel="canonical" href="https://ak33m.com/">
	
		
		<link rel="alternate" type="application/rss+xml" href="%s" title="Akeem Philbert's Blog" />
		

	</head>
	<body>
	</body>
	
	</html>
	`,arg1)
	return nil
}

func theBlogPostsFromTheFeedShouldBeAddedToTheAggregator() error {
	if createdBlog == nil {
		return fmt.Errorf("blog was not created by a previous step")
	}
	if len(createdBlog.Posts) == 0 {
		return fmt.Errorf("expected there to be posts with blog")
	}
	return err
}

func theBlogShouldBeAddedToTheAggregator() error {
	//check that the status code is correct
	if response.StatusCode != http.StatusCreated {
		return fmt.Errorf("expected the status code to be %d, got %d",http.StatusCreated,response.StatusCode)
	}
	//check that the blog was added correctly to the projection
	projections := blogAPI.Application.Projections()
	if len(projections) == 0 {
		return fmt.Errorf("there are no projections configured")
	}
	projection := projections[0].(*api.GORMProjection)
	createdBlog, err = projection.GetBlogByURL(request.(*blogaggregatormodule.AddBlogRequest).Url)
	if err != nil {
		return err
	}

	if createdBlog == nil {
		return fmt.Errorf("blog with urls '%s' does not exist",request.(*blogaggregatormodule.AddBlogRequest).Url)
	}

	if createdBlog.URL != testBlog.URL {
		return fmt.Errorf("expected blog url to be %s, got %s",testBlog.URL,createdBlog.URL)
	}
	return err
}

func theFeedDetailsShouldBeExtracted() error {
	return godog.ErrPending
}

func theFeedHasPosts(arg1 *messages.PickleStepArgument_PickleTable) error {
	var err error
	testFeed = `<?xml version="1.0" encoding="windows-1252"?><rss version="2.0">
	  <channel>
		<title>%s</title>
		<link>https://ak33m.com</link>
		<description>Recent content on Akeem Philbert&#39;s Blog</description>
		<category domain="www.dmoz.com">Computers/Software/Internet/Site Management/Content Management</category>
		<copyright>Copyright 2021 Some Site</copyright>
		<docs>http://blogs.law.harvard.edu/tech/rss</docs>
		<language>en-us</language>
		<lastBuildDate>Tue, 19 Oct 2004 13:39:14 -0400</lastBuildDate>
		<itunes:author>Sojourner Truth</itunes:author>
		<pubDate>Tue, 19 Oct 2004 13:38:55 -0400</pubDate>
		<generator>FeedForAll Beta1 (0.0.1.8)</generator>
		<image>
		  <url>http://www.feedforall.com/ffalogo48x48.gif</url>
		  <title>FeedForAll Sample Feed</title>
		  <link>http://www.feedforall.com/industry-solutions.htm</link>
		  <description>FeedForAll Sample Feed</description>
		  <width>48</width>
		  <height>48</height>
		</image>
		%s
	  </channel>
	</rss>`
	//TODO loop through the table and add feed item to the feed 
	items := ""
	itemColumns := make([]string,len(arg1.Rows[0].Cells))
	for i,_ := range arg1.Rows {
		if i == 0 {
			for j,column := range arg1.Rows[i].Cells {
				itemColumns[j] = column.Value
			}
		} else {
			feedItem := &FeedItem{}
			for j,column := range arg1.Rows[i].Cells {
				if itemColumns[j] == "title" {
					feedItem.Title = column.Value
				}

				if itemColumns[j] == "content" {
					feedItem.Description = column.Value
				}

				if itemColumns[j] == "publish date" {
					feedItem.PublishDate = column.Value
				}
			}
			
			items = items + fmt.Sprintf(`<item>
			<title>%s</title>
			<description>%s</description>
			<link>%s</link>
			<pubDate>%s</pubDate>
		  </item>`,feedItem.Title,feedItem.Link, feedItem.Description,feedItem.PublishDate)

		}
	}


	testFeed = fmt.Sprintf(testFeed,testBlog.Title,items)
	return err
}

func theUrlIsEntered(arg1 string) error {
	request.(*blogaggregatormodule.AddBlogRequest).Url = arg1
	return nil
}

func reset(*godog.Scenario) {
	testBlog = nil
	testUsers = make(map[string]*TestUser)
	testBlogs = make(map[string]*TestBlog)
	err = nil
	createdBlog = nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	err = os.Remove("test.db")//TODO hack to reset the database between runs
	e = echo.New()
	blogAPI = &api.API{}
	blogDataFetched := 0
	blogAPI.Client = testhelpers.NewTestClient(func(req *http.Request) *http.Response {
		blogDataFetched += 1
		//thi is fetching the blog page 
		if blogDataFetched == 1 {
			resp := testhelpers.NewBytesResponse(200,[]byte(testBlogPage))
			resp.Header.Set("Content-Type", "text/html")
			return resp
		}

		resp := testhelpers.NewBytesResponse(200,[]byte(testFeed))
		resp.Header.Set("Content-Type", "application/rss+xml")
		return resp
	})
	weoscontroller.Initialize(e,blogAPI,"../api.yaml")

	ctx.BeforeScenario(reset)
	ctx.Step(`^a pingback url should be generated$`, aPingbackUrlShouldBeGenerated)
	ctx.Step(`^a user named "([^"]*)"$`, aUserNamed)
	ctx.Step(`^an author should be created for each author in the feed$`, anAuthorShouldBeCreatedForEachAuthorInTheFeed)
	ctx.Step(`^an error screen should be shown "([^"]*)"$`, anErrorScreenShouldBeShown)
	ctx.Step(`^"([^"]*)" follows the blog "([^"]*)"$`, followsTheBlog)
	ctx.Step(`^"([^"]*)" has a blog "([^"]*)"$`, hasABlog)
	ctx.Step(`^"([^"]*)" hits the submit button$`, hitsTheSubmitButton)
	ctx.Step(`^"([^"]*)" is logged in$`, isLoggedIn)
	ctx.Step(`^"([^"]*)" is logged in with google$`, isLoggedInWithGoogle)
	ctx.Step(`^"([^"]*)" is not logged in$`, isNotLoggedIn)
	ctx.Step(`^"([^"]*)" is on the blog submit screen$`, isOnTheBlogSubmitScreen)
	ctx.Step(`^posts should be created for each post$`, postsShouldBeCreatedForEachPost)
	ctx.Step(`^profiles for the blog authors should be created$`, profilesForTheBlogAuthorsShouldBeCreated)
	ctx.Step(`^"([^"]*)" should be redirected to the profile page for that blog$`, shouldBeRedirectedToTheProfilePageForThatBlog)
	ctx.Step(`^"([^"]*)" successfully completes the captcha$`, successfullyCompletesTheCaptcha)
	ctx.Step(`^"([^"]*)" successfully submits a feed$`, successfullySubmitsAFeed)
	ctx.Step(`^the aggregator supports atom feeds as well as rss feeds$`, theAggregatorSupportsAtomFeedsAsWellAsRssFeeds)
	ctx.Step(`^the blog details stored in the aggregator$`, theBlogDetailsStoredInTheAggregator)
	ctx.Step(`^the blog has a link to a feed "([^"]*)"$`, theBlogHasALinkToAFeed)
	ctx.Step(`^the blog posts from the feed should be added to the aggregator$`, theBlogPostsFromTheFeedShouldBeAddedToTheAggregator)
	ctx.Step(`^the blog should be added to the aggregator$`, theBlogShouldBeAddedToTheAggregator)
	ctx.Step(`^the feed details should be extracted$`, theFeedDetailsShouldBeExtracted)
	ctx.Step(`^the feed has posts$`, theFeedHasPosts)
	ctx.Step(`^the url "([^"]*)" is entered$`, theUrlIsEntered)
}

func TestSubmitBlog(t *testing.T) {
	status := godog.TestSuite{
		Name: "Submit Blog Feature Test",
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format: "pretty",
		},
	}.Run()
	if status != 0 {
		t.Errorf("there was an error running tests, exit code %d",status)
	}
}