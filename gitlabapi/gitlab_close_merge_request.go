package gitlabapi

import (
	"log"

	"github.com/xanzy/go-gitlab"
)

// TriggerPipeline triggers a pipeline for a given project and branch.
// It returns an error if the pipeline creation fails.
func CloseMerge(client *gitlab.Client, projectID int, sourceBranch string) error {

	state := "opened"
	listOptions := &gitlab.ListProjectMergeRequestsOptions{
		SourceBranch: gitlab.String(sourceBranch),
		State:        gitlab.String(state),
	}

	mergeRequests, _, err := client.MergeRequests.ListProjectMergeRequests(projectID, listOptions)
	if err != nil {
		log.Fatalf("Failed to list merge requests: %v", err)
	}

	// Check and close merge requests with no changes
	for _, mr := range mergeRequests {
		// Fetch the diff between the source and target branches
		compare, _, err := client.Repositories.Compare(projectID, &gitlab.CompareOptions{
			From: &mr.SourceBranch,
			To:   &mr.TargetBranch,
		})
		if err != nil {
			log.Printf("Failed to fetch diff for MR %d: %v", mr.IID, err)
			continue
		}

		// Check if the diff is empty
		if len(compare.Diffs) == 0 {
			// Close the merge request
			_, err := client.MergeRequests.DeleteMergeRequest(projectID, mr.IID)
			if err != nil {
				log.Printf("Failed to close merge request %d: %v", mr.IID, err)
			} else {
				log.Printf("Merge request %d has no changes and has been closed", mr.IID)
			}
		} else {
			log.Printf("Merge request %d has changes and will not be closed", mr.IID)
		}
	}

	return nil
}
