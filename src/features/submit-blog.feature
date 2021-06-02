Feature: Submit blog to aggregator 
  
  Background:
    Given a user named "Akeem"
    And "Akeem" has a blog "https://ak33m.com"
    And the blog has a link to a feed "https://ak33m.com/index.xml"
    And "Akeem" is logged in
    And the feed has posts
    | title       | content             | publish date                    |
    | Post 1      | Lorem Ipsum         | Sat, 27 Mar 2021 17:05:53 -0400 | 
    | Post 2      | Lorem Ipsum         | Sat, 27 Mar 2021 17:05:53 -0400 | 
    And a user named "Francis"
    And "Francis" follows the blog "https://ak33m.com"
    And "Francis" is not logged in 
    And the aggregator supports atom feeds as well as rss feeds
    
  Scenario: Submitting a blog website to aggregator as non logged in user
    Given "Francis" is on the blog submit screen
    And "Francis" successfully completes the captcha
    And the url "https://ak33m.com" is entered 
    When "Francis" hits the submit button 
    Then the blog should be added to the aggregator
    And profiles for the blog authors should be created
    And the blog posts from the feed should be added to the aggregator
    And "Francis" should be redirected to the profile page for that blog 
    
  Scenario: Submitting a blog website to aggregator as logged in user
    Given "Akeem" is on the blog submit screen
    And "Akeem" is logged in with google
    When the url "https://ak33m.com" is entered 
    And "Akeem" hits the submit button 
    Then the blog should be added to the aggregator
    And profiles for the blog authors should be created
    And the blog posts from the feed should be added to the aggregator
    And "Francis" should be redirected to the profile page for that blog 
    
  Scenario: Submitting a blog feed url to aggregator 
    Given "Akeem" is on the blog submit screen
    When the url "https://ak33m.com/index.xml" is entered 
    And "Akeem" hits the submit button
    Then the blog should be added to the aggregator
    And profiles for the blog authors should be created
    And the blog posts from the feed should be added to the aggregator
    And "Akeem" should be redirected to the profile page for that blog 
    
  Scenario: Submitting a website that is not a blog to aggregator 
    Given "Akeem" is on the blog submit screen
    When the url "https://google.com" is entered 
    And "Akeem" hits the submit button
    Then an error screen should be shown "Akeem"
    
    
  Scenario: Processing blog submitted by non logged in user 
    When "Francis" successfully submits a feed 
    Then the feed details should be extracted
    And the blog details stored in the aggregator 
    And an author should be created for each author in the feed
    And posts should be created for each post
    And a pingback url should be generated 
    
    