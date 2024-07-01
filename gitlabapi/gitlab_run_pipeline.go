package gitlabapi

import (
	"fmt"
	"time"

	"github.com/xanzy/go-gitlab"
)

const (
	rateLimit = 5 * time.Second // Rate limit set to 5 seconds
)

// TriggerPipeline triggers a pipeline for a given project and branch.
// It returns an error if the pipeline creation fails.
func TriggerPipeline(client *gitlab.Client, projectID int, branch string) error {
	// Create a new pipeline
	pipeline, _, err := client.Pipelines.CreatePipeline(projectID, &gitlab.CreatePipelineOptions{
		Ref: gitlab.String(branch),
	})
	if err != nil {
		return fmt.Errorf("failed to create pipeline for project: %v", err)
	}

	fmt.Printf("Pipeline created with ID: %d\n", pipeline.ID)

	// Sleep to avoid hitting the rate limit
	time.Sleep(rateLimit)

	return nil // Return nil if the pipeline was created successfully
}
