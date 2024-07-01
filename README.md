# GitLab Bulk Actions Script

This script allows you to perform various bulk actions on a GitLab project directly from the command line. The actions include creating a `.gitignore` file, managing merge requests, deleting files, creating branches, triggering pipelines, and changing project rules.

## Prerequisites

- Go (Golang) installed on your machine.
- GitLab API token with necessary permissions.
- GitLab project ID.

## Installation

1. Clone this repository to your local machine:
   ```sh
   git clone https://github.com/yourusername/your-repo.git
   cd your-repo

Usage
Run the script with the desired action and provide necessary arguments. Here are the available actions:

create-gitignore

accept-merge-request

delete-files

create-branch

trigger-pipeline

create-mr

close-mr

change-project-rules
