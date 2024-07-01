module main

go 1.22.2

replace gitlabapi/gitlabapi => ./gitlabapi

require (
	github.com/xanzy/go-gitlab v0.105.0
	gitlabapi/gitlabapi v0.0.0-00010101000000-000000000000
)

require (
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-retryablehttp v0.7.6 // indirect
	golang.org/x/oauth2 v0.20.0 // indirect
	golang.org/x/time v0.5.0 // indirect
)
