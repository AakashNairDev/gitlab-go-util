package gitlabapi

import (
	"fmt"
	"github.com/xanzy/go-gitlab"
)

// CreateBranchAndProtect creates a new branch from a reference branch and protects it
func CreateBranchAndProtect(client *gitlab.Client, projectID int, refBranch string, newBranch string) error {
	// Check if the reference branch exists
	refExists, err := checkBranchExists(client, projectID, refBranch)
	if err != nil {
		return fmt.Errorf("failed to check if reference branch exists: %v", err)
	}
	if !refExists {
		return fmt.Errorf("reference branch does not exist: %s", refBranch)
	}

	// Check if the new branch already exists
	newExists, err := checkBranchExists(client, projectID, newBranch)
	if err != nil {
		return fmt.Errorf("failed to check if new branch exists: %v", err)
	}
	if newExists {
		fmt.Printf("Skipping project because branch %s already exists\n", newBranch)
		return nil
	}

	// Create the branch
	branch, _, err := client.Branches.CreateBranch(projectID, &gitlab.CreateBranchOptions{
		Branch: gitlab.String(newBranch),
		Ref:    gitlab.String(refBranch),
	})
	if err != nil {
		return fmt.Errorf("failed to create branch: %v", err)
	}

	// Protect the branch
	_, _, err = client.ProtectedBranches.ProtectRepositoryBranches(projectID, &gitlab.ProtectRepositoryBranchesOptions{
		Name: gitlab.String(newBranch),
		PushAccessLevel: gitlab.AccessLevel(gitlab.NoPermissions),
		MergeAccessLevel: gitlab.AccessLevel(gitlab.MaintainerPermissions),
	})
	if err != nil {
		return fmt.Errorf("failed to protect branch: %v", err)
	}

	fmt.Printf("Created branch: %s\n", branch.Name)
	return nil
}

// checkBranchExists checks if a branch exists in the project
func checkBranchExists(client *gitlab.Client, projectID int, branchName string) (bool, error) {
	branches, _, err := client.Branches.ListBranches(projectID, &gitlab.ListBranchesOptions{
		Search: gitlab.String(branchName),
	})
	if err != nil {
		return false, err
	}

	for _, branch := range branches {
		if branch.Name == branchName {
			return true, nil
		}
	}

	return false, nil
}