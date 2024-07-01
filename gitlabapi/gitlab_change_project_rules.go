package gitlabapi

import (
	"fmt"
	"log"

	"github.com/xanzy/go-gitlab"
)

// ChangeProjectRules changes the push rules for the specified project.
func ChangeProjectRules(client *gitlab.Client, projectID int, projectName, newRegex string) error {
	_, _, err := client.Projects.EditProjectPushRule(projectID, &gitlab.EditProjectPushRuleOptions{
		BranchNameRegex: &newRegex,
	})
	if err != nil {
		return fmt.Errorf("failed to update push rule for project %s: %v", projectName, err)
	}

	log.Printf("Updated push rule for project %s", projectName)
	return nil
}
