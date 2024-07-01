package gitlabapi

import (
	"fmt"
	"log"

	"github.com/xanzy/go-gitlab"
)

// TriggerPipeline triggers a pipeline for a given project and branch.
// It returns an error if the pipeline creation fails.
func CreateMerge(client *gitlab.Client, projectID int, sourceBranch, targetBranch string) error {
	// Create a new pipeline
	title := fmt.Sprintf("Merge %s into %s", sourceBranch, targetBranch)
	mergeRequest, _, err := client.MergeRequests.CreateMergeRequest(projectID, &gitlab.CreateMergeRequestOptions{
		SourceBranch: gitlab.String(sourceBranch),
		TargetBranch: gitlab.String(targetBranch),
		Title:        gitlab.String(title),
	})
	if err != nil {
		log.Fatalf("Failed to create merge request: %v", err)
	}

	fmt.Printf("Merge request created successfully: %s\n", mergeRequest.WebURL)

	return nil
}
