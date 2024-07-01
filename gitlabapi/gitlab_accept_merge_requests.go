package gitlabapi

import (
	"fmt"
	"log"
	"strings"

	"github.com/xanzy/go-gitlab"
)

// getPipelineStatus retrieves the pipeline status for the given branch.
func getPipelineStatus(client *gitlab.Client, projectID int, branch string) (string, error) {
	pipelines, _, err := client.Pipelines.ListProjectPipelines(projectID, &gitlab.ListProjectPipelinesOptions{
		Ref: gitlab.String(branch),
	})
	if err != nil {
		return "", err
	}

	if len(pipelines) > 0 {
		return pipelines[0].Status, nil
	}

	return "none", nil
}

// acceptMergeRequest accepts the merge request with the given ID.
func acceptMergeRequest(client *gitlab.Client, projectID int, mergeRequestIID int) error {
	_, _, err := client.MergeRequests.AcceptMergeRequest(projectID, mergeRequestIID, &gitlab.AcceptMergeRequestOptions{
		MergeWhenPipelineSucceeds: gitlab.Bool(true),
	})
	return err
}

// AcceptMergeRequests accepts merge requests based on the pipeline status and branch name.
func AcceptMergeRequests(client *gitlab.Client, projectID int, branchPrefix string) error {
	// List all merge requests for the project
	mergeRequests, _, err := client.MergeRequests.ListProjectMergeRequests(projectID, &gitlab.ListProjectMergeRequestsOptions{
		State: gitlab.String("opened"),
	})
	if err != nil {
		return err
	}

	// Check if there are any merge requests with the specified branch prefix
	hasMatchingMergeRequest := false
	for _, mr := range mergeRequests {
		if strings.HasPrefix(mr.SourceBranch, branchPrefix) {
			hasMatchingMergeRequest = true
			break
		}
	}

	// If there are no matching merge requests, skip the project
	if !hasMatchingMergeRequest {
		log.Printf("No matching merge requests found for project %d. Skipping project.\n", projectID)
		return nil
	}

	// Iterate over the merge requests
	for _, mr := range mergeRequests {
		// Check if the merge request's source branch matches the provided prefix
		if !strings.HasPrefix(mr.SourceBranch, branchPrefix) {
			continue
		}

		// Get the pipeline status for the merge request's source branch
		pipelineStatus, err := getPipelineStatus(client, projectID, mr.SourceBranch)
		if err != nil {
			log.Printf("Error getting pipeline status for branch %s: %v", mr.SourceBranch, err)
			continue
		}

		// If the pipeline is successful, accept the merge request
		if pipelineStatus == "success" {
			err := acceptMergeRequest(client, projectID, mr.IID)
			if err != nil {
				log.Printf("Error accepting merge request %d: %v", mr.IID, err)
				continue
			}
			fmt.Printf("Merge request %d has been accepted.\n", mr.IID)
		} else {
			fmt.Printf("Pipeline for branch %s has status %s. Skipping merge request %d.\n", mr.SourceBranch, pipelineStatus, mr.IID)
		}
	}

	return nil
}