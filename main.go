package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"gitlabapi/gitlabapi" // Make sure this is the correct import path for your gitlabapi package

	"github.com/xanzy/go-gitlab"
)

const (
	rateLimit = 5 * time.Second // Rate limit set to 5 seconds
	perPage   = 20              // Number of projects to process per page
)

func main() {
	client, err := gitlabapi.NewGitLabClient()
	if err != nil {
		log.Fatalf("Failed to create GitLab client: %v", err)
	}

	// Create a rate limiter with the desired rate limit
	limiter := gitlabapi.NewRateLimiter(rateLimit)

	// Prompt for the group name
	fmt.Print("Enter the group name: ")
	reader := bufio.NewReader(os.Stdin)
	groupName, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Failed to read group name: %v", err)
	}
	groupName = strings.TrimSpace(groupName)

	// Get the group ID by searching for the group
	var groups []*gitlab.Group
	err = gitlabapi.UseRateLimiter(limiter, func() error {
		var err error
		groups, _, err = client.Groups.ListGroups(&gitlab.ListGroupsOptions{
			Search: gitlab.String(groupName),
		})
		return err
	})
	if err != nil {
		log.Fatalf("Failed to search for group: %v", err)
	}

	if len(groups) == 0 {
		log.Fatalf("No group found with name: %s", groupName)
	}

	groupID := groups[0].ID

	// Prompt for the action to perform
	fmt.Print("Enter action (create-gitignore/accept-merge-request/delete-car-files/create-branch/trigger-pipeline/create-mr/close-mr/change-project-rules): ")
	action, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Failed to read action: %v", err)
	}
	action = strings.TrimSpace(action)

	// Declare variable here
	var branchPrefix string
	var refBranch string
	var newBranch string
	var triggerBranch string
	var sourceBranch string
	var targetBranch string
	var closeBranch string
	var newRegex string

	if action == "accept-merge-request" {
		// Prompt for the branch name prefix
		fmt.Print("Enter the branch name prefix: ")
		branchPrefix, err = reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Failed to read branch name prefix: %v", err)
		}
		branchPrefix = strings.TrimSpace(branchPrefix)
	}

	if action == "close-mr" {
		// Prompt for the branch name prefix
		fmt.Print("Enter the branch name prefix: ")
		closeBranch, err = reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Failed to read branch name prefix: %v", err)
		}
		closeBranch = strings.TrimSpace(closeBranch)
	}

	if action == "create-mr" {
		// Prompt for the source branch name prefix
		fmt.Print("Enter the Source branch name prefix: ")
		sourceBranch, err = reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Failed to read branch name prefix: %v", err)
		}
		sourceBranch = strings.TrimSpace(sourceBranch)

		// Prompt for the target branch name prefix
		fmt.Print("Enter the Source branch name prefix: ")
		targetBranch, err = reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Failed to read branch name prefix: %v", err)
		}
		targetBranch = strings.TrimSpace(targetBranch)
	}

	if action == "change-project-rules" {

		// Prompt for the target branch name prefix
		fmt.Print("Enter future regex: ")
		newRegex, err = reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Failed to read branch name prefix: %v", err)
		}
		newRegex = strings.TrimSpace(targetBranch)

	}

	if action == "trigger-pipeline" {
		//Promy for the branch
		fmt.Print("Enter the branch name to trigger pipeline: ")
		triggerBranch, err = reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Failed to read branch name: %v", err)
		}
		triggerBranch = strings.TrimSpace(triggerBranch)
	}

	// Prompt for the reference branch and new branch names if the action is "create-branch"
	if action == "create-branch" {
		// Prompt for the reference branch
		fmt.Print("Enter the reference branch: ")
		refBranch, err = reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Failed to read reference branch: %v", err)
		}
		refBranch = strings.TrimSpace(refBranch)

		// Prompt for the new branch
		fmt.Print("Enter the new branch: ")
		newBranch, err = reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Failed to read new branch: %v", err)
		}
		newBranch = strings.TrimSpace(newBranch)
	}

	// Process projects in chunks of 20
	page := 1
	for {
		var projects []*gitlab.Project
		err = gitlabapi.UseRateLimiter(limiter, func() error {
			var err error
			projects, _, err = gitlabapi.ListProjects(client, groupID, page, perPage)
			return err
		})
		if err != nil {
			log.Fatalf("Failed to list projects: %v", err)
		}

		for _, project := range projects {
			fmt.Printf("Processing project ID: %d, Name: %s\n", project.ID, project.Name)

			// Apply the selected action to each project in the chunk
			switch action {
			case "change-project-rules":

				err = gitlabapi.UseRateLimiter(limiter, func() error {
					return gitlabapi.ChangeProjectRules(client, project.ID, project.Name, newRegex)
				})
				if err != nil {
					log.Printf("Failed to create branch and protect for project %s: %v\n", project.Name, err)
				}
			case "close-mr":
				err = gitlabapi.UseRateLimiter(limiter, func() error {
					return gitlabapi.CloseMerge(client, project.ID, closeBranch)
				})
				if err != nil {
					log.Printf("Failed to create branch and protect for project %s: %v\n", project.Name, err)
				}
			case "create-mr":
				err = gitlabapi.UseRateLimiter(limiter, func() error {
					return gitlabapi.CreateMerge(client, project.ID, sourceBranch, targetBranch)
				})
				if err != nil {
					log.Printf("Failed to create branch and protect for project %s: %v\n", project.Name, err)
				}
			case "trigger-pipeline":
				err = gitlabapi.UseRateLimiter(limiter, func() error {
					return gitlabapi.TriggerPipeline(client, project.ID, triggerBranch)
				})
				if err != nil {
					log.Printf("Failed to create branch and protect for project %s: %v\n", project.Name, err)
				}
			case "create-branch":
				err = gitlabapi.UseRateLimiter(limiter, func() error {
					return gitlabapi.CreateBranchAndProtect(client, project.ID, refBranch, newBranch)
				})
				if err != nil {
					log.Printf("Failed to create branch and protect for project %s: %v\n", project.Name, err)
				}
			case "create-gitignore":
				err = gitlabapi.UseRateLimiter(limiter, func() error {
					branchName := "feature/add-gitignore"
					ignorePath := "assets/gitignore"
					return gitlabapi.CreateBranchAndIgnore(client, project.ID, branchName, ignorePath)
				})
				if err != nil {
					log.Printf("Failed to create branch and .gitignore for project %s: %v\n", project.Name, err)
					// Continue to the next project after an error
					continue
				}
			case "accept-merge-request":
				err = gitlabapi.UseRateLimiter(limiter, func() error {
					return gitlabapi.AcceptMergeRequests(client, project.ID, branchPrefix)
				})
				if err != nil {
					log.Printf("Failed to accept merge requests for project %s: %v\n", project.Name, err)
					// Continue to the next project after an error
					continue
				}
			case "delete-car-files":
				err = gitlabapi.UseRateLimiter(limiter, func() error {
					return gitlabapi.DeleteCarFilesAndCreateMergeRequest(client, project.ID)
				})
				if err != nil {
					log.Printf("Failed to delete .car files and create merge requests for project %s: %v\n", project.Name, err)
					// Continue to the next project after an error
					continue
				}
			default:
				log.Fatalf("Invalid action: %s", action)
			}
		}

		// Break the loop if there are no more projects
		if len(projects) < perPage {
			break
		}

		// Move to the next page
		page++
	}
}
