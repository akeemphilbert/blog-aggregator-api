@wip
Feature: View blog posts
  
  This covers viewing posts from various bloggers as well as filtering posts 
  
  Background:
    
    Given a user Marcus
    And Marcus has permissions to view blog posts
    And the aggregator has blogs
    | id | title        | url                      | feedUrl                             | author | lastUpdated                      | 
    | 1  | Akeem's Blog | https://ak33m.com        | https://ak33m.com/index.xml         | Akeem  | Sat, 27 Mar 2021 17:05:53 -0400  |
    | 2  | Blog 2       | https://blog.example.org | https://blog.example.org/index.xml  | John   | Sat, 12 Jun 2021 17:05:53 -0400  |
    And the aggregator has posts
    | id | blogId | title             | description                    | tags                            | publishDate                     | views |
    | 1  | 1      | It's Over Magento | Magento Love Hate Relationship | php,e-commerce,magento          | Sat, 27 Mar 2021 17:05:53 -0400 | 1     |
    | 2  | 1      | Viro React        | Viro React - Augmented Reality | react, ar                       | Mon, 13 Apr 2020 15:57:22 -0400 | 2     |
    | 3  | 2      | Lorem Ipsum.      | Dolor                          | ar                              | Sat, 12 Jun 2021 15:57:22 -0400 | 3     |                        
    
  Scenario: View posts of a specific blog
    When Marcus selects a blog with id "1"
    Then Marcus should see a list of blog posts 
    | id | blogId | title             | description                    | tags                            | publishDate                     | views |
    | 1  | 1      | It's Over Magento | Magento Love Hate Relationship | php,e-commerce,magento          | Sat, 27 Mar 2021 17:05:53 -0400 | 1     |
    | 2  | 1      | Viro React        | Viro React - Augmented Reality | react, ar                       | Mon, 13 Apr 2020 15:57:22 -0400 | 2     | 
  
  Scenario: View posts by category
    When Marcus selects a category "ar"
    Then Marcus should see a list of blog posts 
    | id | blogId | title             | description                    | tags                            | publishDate                     | views |
    | 2  | 1      | Viro React        | Viro React - Augmented Reality | react, ar                       | Mon, 13 Apr 2020 15:57:22 -0400 | 2     | 
    | 3  | 2      | Lorem Ipsum.      | Dolor                          | ar                              | Sat, 12 Jun 2021 15:57:22 -0400 | 3     |     
    
  Scenario: View most recent posts
    Given The current date is "07/10/21"
    When Marcus views recent posts
    Then Marcus should see posts 30 days from the current date 
    And Marcus should see a list of blog posts 
    | id | blogId | title             | description                    | tags                            | publishDate                     | views |
    | 3  | 2      | Lorem Ipsum.      | Dolor                          | ar                              | Sat, 12 Jun 2021 15:57:22 -0400 | 2     |     
  
  Scenario: View posts by highest views
    When Marcus views posts by highest views
    Then Marcus should see a list of blog posts 
    | id | blogId | title             | description                    | tags                            | publishDate                     | views |
    | 3  | 2      | Lorem Ipsum.      | Dolor                          | ar                              | Sat, 12 Jun 2021 15:57:22 -0400 | 3     |
    | 2  | 1      | Viro React        | Viro React - Augmented Reality | react, ar                       | Mon, 13 Apr 2020 15:57:22 -0400 | 2     |
    | 1  | 1      | It's Over Magento | Magento Love Hate Relationship | php,e-commerce,magento          | Sat, 27 Mar 2021 17:05:53 -0400 | 1     |
    
  