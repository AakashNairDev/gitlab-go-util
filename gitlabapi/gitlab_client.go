package gitlabapi

import (
	"log"
	"os"
	"time"

	"github.com/xanzy/go-gitlab"
)

const (
	maxRetries = 5
	retryDelay = 2 * time.Second
	timeout    = 10 * time.Second
)

// NewGitLabClient creates a new GitLab client with the specified token and URL.
func NewGitLabClient() (*gitlab.Client, error) {
	gitlabToken := os.Getenv("GITLAB_TOKEN")
	if gitlabToken == "" {
		log.Fatal("GITLAB_TOKEN environment variable is not set")
	}

	gitlabURL := os.Getenv("GITLAB_URL")
	if gitlabURL == "" {
		log.Fatal("GITLAB_URL environment variable is not set")
	}

	git, err := gitlab.NewClient(gitlabToken, gitlab.WithBaseURL(gitlabURL))
	if err != nil {
		return nil, err
	}

	return git, nil
}

// Retry is a helper function to retry failed operations.
func Retry(attempts int, sleep time.Duration, f func() error) error {
	if err := f(); err != nil {
		if attempts--; attempts > 0 {
			log.Printf("Retrying after error: %s\n", err)
			time.Sleep(sleep)
			return Retry(attempts, sleep, f)
		}
		return err
	}
	return nil
}

// RateLimiter is a simple rate limiter that can be used to limit the rate at which
// API calls are made to GitLab.
type RateLimiter struct {
	rate     time.Duration
	throttle chan time.Time
}

// NewRateLimiter creates a new rate limiter with the specified rate.
func NewRateLimiter(rate time.Duration) *RateLimiter {
	r := &RateLimiter{
		rate:     rate,
		throttle: make(chan time.Time, 1),
	}
	r.throttle <- time.Now()
	return r
}

// Wait blocks until enough time has passed since the last call to Wait to respect the rate limit.
func (r *RateLimiter) Wait() {
	t := <-r.throttle
	elapsed := time.Since(t)
	if elapsed < r.rate {
		time.Sleep(r.rate - elapsed)
	}
	r.throttle <- time.Now()
}

// UseRateLimiter is a helper function that wraps a function with a rate limiter.
func UseRateLimiter(limiter *RateLimiter, f func() error) error {
	limiter.Wait()
	return f()
}
