// pagination.go
package gitlabapi

import (
	"github.com/xanzy/go-gitlab"
)

// Paginate makes multiple API calls to handle pagination for the given function.
func Paginate(client *gitlab.Client, groupID int, f func(options gitlab.ListOptions) (*gitlab.Response, error)) error {
	page := 1
	perPage := 20

	for {
		response, err := f(gitlab.ListOptions{
			Page:    page,
			PerPage: perPage,
		})
		if err != nil {
			return err
		}

		if response.NextPage == 0 {
			break
		}

		page = response.NextPage
	}

	return nil
}