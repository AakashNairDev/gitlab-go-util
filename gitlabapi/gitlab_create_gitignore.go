package gitlabapi

import (
	"fmt"
	"io/ioutil"
	"strings"
	"log"

	"github.com/xanzy/go-gitlab"
)

// CreateBranchAndIgnore creates a branch and adds or updates a .gitignore file.
// It will not retry if the branch is already available.
func CreateBranchAndIgnore(client *gitlab.Client, projectID int, branchName, ignorePath string) error {
	// Check if the branch already exists
	branches, _, err := client.Branches.ListBranches(projectID, &gitlab.ListBranchesOptions{
		Search: gitlab.String(branchName),
	})
	if err != nil {
		return err
	}

	for _, branch := range branches {
		if strings.EqualFold(branch.Name, branchName) {
			fmt.Printf("Skipping project %d because branch %s already exists\n", projectID, branchName)
			return nil // Skip the project if the branch already exists
		}
	}

	// Create the branch
	err = Retry(maxRetries, retryDelay, func() error {
		_, _, err := client.Branches.CreateBranch(projectID, &gitlab.CreateBranchOptions{
			Branch: gitlab.String(branchName),
			Ref:    gitlab.String("develop"), // Use develop as the reference branch
		})
		return err
	})
	if err != nil {
		return err
	}

	// Read the .gitignore file content
	gitignoreContent, err := ioutil.ReadFile(ignorePath)
	if err != nil {
		return err
	}

	// Check if the .gitignore file already exists on the feature branch
	_, _, err = client.RepositoryFiles.GetFile(projectID, ".gitignore", &gitlab.GetFileOptions{
		Ref: gitlab.String(branchName),
	})

	var action gitlab.FileActionValue
	if err == nil {
		// The .gitignore file already exists, so we update it
		action = gitlab.FileUpdate
	} else {
		// The .gitignore file does not exist, so we create it
		action = gitlab.FileCreate
	}

	// Add or update the .gitignore file on the feature branch
	commitAction := &gitlab.CommitActionOptions{
		Action:   gitlab.FileAction(action),
		FilePath: gitlab.String(".gitignore"),
		Content:  gitlab.String(string(gitignoreContent)), // Convert byte slice to string	
	}
	_, _, err = client.Commits.CreateCommit(projectID, &gitlab.CreateCommitOptions{
		Branch:        gitlab.String(branchName),
		CommitMessage: gitlab.String("Add or update .gitignore"),
		Actions:       []*gitlab.CommitActionOptions{commitAction},
	})
	if err != nil {
		log.Fatalf("Failed to add or update .gitignore for project %s: %v", projectID, err)
	}
	fmt.Printf("Added or updated .gitignore for project: %s\n", projectID)	

	// Create a merge request
	targetBranch := "develop" // The branch you want to merge into
	title := fmt.Sprintf("Merge request from %s to %s", branchName, targetBranch)
	err = Retry(maxRetries, retryDelay, func() error {
		_, _, err := client.MergeRequests.CreateMergeRequest(projectID, &gitlab.CreateMergeRequestOptions{
			SourceBranch: gitlab.String(branchName),
			TargetBranch: gitlab.String(targetBranch),
			Title:        gitlab.String(title),
		})
		return err
	})
	if err != nil {
		return err
	}

	return nil
}