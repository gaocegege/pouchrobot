package open

import (
	"strings"

	"github.com/pouchcontainer/pouchrobot/server/utils"

	"github.com/google/go-github/github"
	"github.com/sirupsen/logrus"
)

var (
	// XS is
	XS = 10
	// S is
	S = 40
	// M is
	M = 80
	// L is
	L = 160
	// XL is
	XL = 640
)

// ParseToGeneratePRLabels parses
func ParseToGeneratePRLabels(pr *github.PullRequest) []string {
	var labels []string
	labels = append(labels, ParseToGetPRSize(pr))
	labels = append(labels, ParseTitleToGenerateLabels(pr)...)
	return utils.UniqueElementSlice(labels)
}

// ParseToGetPRSize parses the pr additions and deletions
func ParseToGetPRSize(pr *github.PullRequest) string {
	if pr.Additions == nil || pr.Deletions == nil {
		logrus.Errorf("pull request %d has no field value of Additions or Deletions when ParseToGetPRSize", *(pr.Number))
		return "size/XS"
	}
	if *(pr.Additions)+*(pr.Deletions) <= XS {
		return "size/XS"
	}
	if *(pr.Additions)+*(pr.Deletions) <= S {
		return "size/S"
	}
	if *(pr.Additions)+*(pr.Deletions) <= M {
		return "size/M"
	}
	if *(pr.Additions)+*(pr.Deletions) <= L {
		return "size/L"
	}
	if *(pr.Additions)+*(pr.Deletions) <= XL {
		return "size/XL"
	}
	return "size/XXL"
}

// ParseTitleToGenerateLabels parses
func ParseTitleToGenerateLabels(pr *github.PullRequest) []string {
	if pr.Title == nil {
		logrus.Errorf("pull request %d has no title when ParseTitleToGenerateLabels", *(pr.Number))
		return nil
	}
	var labels []string
	title := pr.Title
	for label, matchedSlice := range utils.TitleMatches {
		for _, pattern := range matchedSlice {
			lowerCaseTitle := strings.ToLower(*title)
			if strings.Contains(lowerCaseTitle, pattern) {
				labels = append(labels, label)
				break
			}
		}
	}
	return labels
}
