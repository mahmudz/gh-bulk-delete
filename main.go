package main

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/cli/go-gh"
)

func deleteRepo(Repo_Name string) {
	fmt.Println("Deleting Repo: " + Repo_Name)

	client, err := gh.RESTClient(nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	response := struct {
		Full_Name string
		Name      string
	}{}

	err = client.Delete("repos/"+Repo_Name, &response)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func main() {
	selectedType := ""
	repoTypeSelectPrompt := &survey.Select{
		Message: "Select Repo Type:",
		Options: []string{
			"All",
			"Public",
			"Private",
			"Sources",
			"Forks",
			"Archived",
			"Mirrors",
			"Templates",
		},
	}
	survey.AskOne(repoTypeSelectPrompt, &selectedType)

	client, err := gh.RESTClient(nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	response := []struct {
		Full_Name string
		Name      string
	}{}

	err = client.Get("user/repos?type="+strings.ToLower(selectedType), &response)
	if err != nil {
		fmt.Println(err)
		return
	}

	repoList := []string{}
	for _, repo := range response {
		repoList = append(repoList, repo.Full_Name)
	}

	selectedRepos := []string{}
	prompt := &survey.MultiSelect{
		Message: "Select Repos:",
		Options: repoList,
	}
	survey.AskOne(prompt, &selectedRepos)

	if len(selectedRepos) == 0 {
		fmt.Println("No Repos Selected")
		return
	}

	confirmation := false
	confirmationPrompt := &survey.Confirm{
		Message: "Are you sure you want to delete these repos?",
	}
	survey.AskOne(confirmationPrompt, &confirmation)

	if confirmation {
		for _, repo := range selectedRepos {
			deleteRepo(repo)
		}
	}
}
