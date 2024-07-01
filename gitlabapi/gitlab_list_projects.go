package gitlabapi

import (
	"github.com/xanzy/go-gitlab"
)

// ListProjects lists projects in the specified group with pagination.
func ListProjects(client *gitlab.Client, groupID, page, perPage int) ([]*gitlab.Project, *gitlab.Response, error) {
	projects, resp, err := client.Groups.ListGroupProjects(groupID, &gitlab.ListGroupProjectsOptions{
		ListOptions: gitlab.ListOptions{
			Page:    page,
			PerPage: perPage,
		},
	})
	if err != nil {
		return nil, nil, err
	}

	return projects, resp, nil
}
