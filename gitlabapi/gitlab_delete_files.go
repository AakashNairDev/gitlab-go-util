// gitlabapi/gitlabapi.go

package gitlabapi

import (
	"fmt"
	"log"
	"regexp"

	"github.com/xanzy/go-gitlab"
)

// DeleteCarFilesAndCreateMergeRequest deletes .car files from the specified project
// and creates a merge request with the deletions, if any .car files are found.
func DeleteCarFilesAndCreateMergeRequest(client *gitlab.Client, projectID int) error {
	// Checkout to a new feature branch if it doesn't exist
	branchName := "feature/delete-car-files"
	ref, _, err := client.Branches.GetBranch(projectID, branchName)
	if err != nil || ref == nil {
		// Create a new branch if it doesn't exist
		_, _, err = client.Branches.CreateBranch(projectID, &gitlab.CreateBranchOptions{
			Branch: gitlab.String(branchName),
			Ref:    gitlab.String("develop"),
		})
		if err != nil {
			return fmt.Errorf("failed to create branch for project %d: %v", projectID, err)
		}
	}

	// List all files in the project
	tree, _, err := client.Repositories.ListTree(projectID, &gitlab.ListTreeOptions{
		Ref:       gitlab.String(branchName),
		Recursive: gitlab.Bool(true),
	})
	if err != nil {
		return fmt.Errorf("failed to list files for project %d: %v", projectID, err)
	}

	// Filter files to find those that match the regex pattern
	carRegex := regexp.MustCompile(`\.car$`)
	deletedFiles := false
	for _, item := range tree {
		if carRegex.MatchString(item.Path) && item.Type == "blob" {
			// Delete the matched file
			_, err := client.RepositoryFiles.DeleteFile(projectID, item.Path, &gitlab.DeleteFileOptions{
				Branch:        gitlab.String(branchName),
				CommitMessage: gitlab.String("Delete car file"),
			})
			if err != nil {
				log.Printf("Failed to delete file %s for project %d: %v\n", item.Path, projectID, err)
				continue
			}
			deletedFiles = true
		}
	}

	// If no .car files were deleted, skip the project
	if !deletedFiles {
		log.Printf("No .car files found for project %d, skipping project\n", projectID)
		return nil
	}

	// Commit the changes
	_, _, err = client.Commits.CreateCommit(projectID, &gitlab.CreateCommitOptions{
		Branch:        gitlab.String(branchName),
		CommitMessage: gitlab.String("Delete car files"),
	})
	if err != nil {
		return fmt.Errorf("failed to commit changes for project %d: %v", projectID, err)
	}

	// Create a merge request
	targetBranch := "develop" // The branch you want to merge into
	title := fmt.Sprintf("Merge request to delete car files")
	_, _, err = client.MergeRequests.CreateMergeRequest(projectID, &gitlab.CreateMergeRequestOptions{
		SourceBranch: gitlab.String(branchName),
		TargetBranch: gitlab.String(targetBranch),
		Title:        gitlab.String(title),
	})
	if err != nil {
		return fmt.Errorf("failed to create merge request for project %d: %v", projectID, err)
	}

	return nil
}
