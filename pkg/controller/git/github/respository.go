package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/v31/github"
	redhatcopv1alpha1 "github.com/redhat-cop/project-initialize-operator/project-initialize/pkg/apis/redhatcop/v1alpha1"
)

func CheckForGitOpsRepository(teamName string, owner string, token string) (bool, error) {
	// Check if repo exists
	client := InitializeGitHubClient(token)
	repoName := GetTeamRepoName(teamName)

	repos, _, err := client.Repositories.List(context.Background(), owner, nil)
	if err != nil {
		return false, err
	}
	// Search repos to see if the teams repo exists
	for _, repo := range repos {
		if repo.GetName() == repoName {
			return true, nil
		}
	}
	return false, nil
}

func CreateNewRespository(teamName string, token string, templateOwner string, templateRepo string, gitDetails *redhatcopv1alpha1.Git) error {
	client := InitializeGitHubClient(token)

	repoName := GetTeamRepoName(teamName)
	newRepo := getTemplateRequest(repoName, gitDetails.Owner, gitDetails.Private, gitDetails.Desc)
	repo, _, err := client.Repositories.CreateFromTemplate(context.Background(), templateOwner, templateRepo, newRepo)
	if err != nil {
		return err
	}
	log.Info(fmt.Sprintf("Successfully created new repo: %v\n", repo.GetName()))
	return nil
}

func getTemplateRequest(repoName string, owner string, private bool, description string) *github.TemplateRepoRequest {
	return &github.TemplateRepoRequest{
		Name:        &repoName,
		Owner:       &owner,
		Private:     &private,
		Description: &description,
	}
}

func GetTeamRepoName(teamName string) string {
	return fmt.Sprintf("%s-gitops", teamName)
}