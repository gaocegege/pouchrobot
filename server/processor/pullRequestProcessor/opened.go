package pullRequestProcessor

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/pouchcontainer/pouchrobot/server/processor/pullRequestProcessor/open"
	"github.com/pouchcontainer/pouchrobot/server/utils"

	"github.com/google/go-github/github"
)

// ActToPROpened acts a pull request opened event.
func (prp *PullRequestProcessor) ActToPROpened(pr *github.PullRequest) error {
	prp.attachLabels(pr)
	prp.attachComments(pr)
	return nil
}

func (prp *PullRequestProcessor) attachLabels(pr *github.PullRequest) error {
	// attach labels
	labels := open.ParseToGeneratePRLabels(pr)
	if len(labels) == 0 {
		return nil
	}
	return prp.Client.AddLabelsToIssue(*(pr.Number), labels)
}

func (prp *PullRequestProcessor) attachComments(pr *github.PullRequest) error {
	// check pull request whether title is sufficient
	prp.attachTitleComments(pr)

	// check pull request whether description is sufficient
	prp.attachBodyComments(pr)

	// check whether this pull request is signed off
	prp.addSignoffComments(pr)

	// check whether this contributor is the first time contributor
	prp.attachFirstContributionComments(pr)

	return nil
}

func (prp *PullRequestProcessor) attachTitleComments(pr *github.PullRequest) error {
	if pr.Title != nil && len(*(pr.Title)) > 20 {
		return nil
	}

	// attach comment
	body := fmt.Sprintf(utils.PRTitleTooShort, *(pr.User.Login))
	newComment := &github.IssueComment{
		Body: &body,
	}

	return prp.Client.AddCommentToPR(*(pr.Number), newComment)
}

func (prp *PullRequestProcessor) attachBodyComments(pr *github.PullRequest) error {
	if pr.Body != nil && len(*(pr.Body)) > 50 {
		return nil
	}

	body := fmt.Sprintf(utils.PRDescriptionTooShort, *(pr.User.Login))
	newComment := &github.IssueComment{
		Body: &body,
	}

	return prp.Client.AddCommentToPR(*(pr.Number), newComment)
}

func (prp *PullRequestProcessor) addSignoffComments(pr *github.PullRequest) error {
	// check whether commits are following the rules
	commits, err := prp.Client.ListCommits(*(pr.Number))
	if err != nil {
		return err
	}

	needSignoff := false
	for _, commit := range commits {
		if commit.Commit != nil && !dcoRegex.MatchString(*commit.Commit.Message) {
			needSignoff = true
			break
		}
	}

	if !needSignoff {
		return nil
	}

	body := fmt.Sprintf(utils.PRNeedsSignOff, *(pr.User.Login))
	newComment := &github.IssueComment{
		Body: &body,
	}

	return prp.Client.AddCommentToPR(*(pr.Number), newComment)
}

// attachFirstContributionComments attaches a first contributor comments when
// it is the first time for author to contribute.
func (prp *PullRequestProcessor) attachFirstContributionComments(pullRequest *github.PullRequest) error {
	// since webhook pull requests are different from raw pull request from GET api,
	// we need to get a brand new pull request from GitHub.
	pr, err := prp.Client.GetSinglePR(*(pullRequest.Number))
	if err != nil {
		return err
	}
	// check whether this is the first contributor of the committer
	if pr.AuthorAssociation == nil {
		return nil
	}

	logrus.Infof("Author in pr %d is %s", *(pr.Number), *(pr.AuthorAssociation))
	if !isFirstContribution(*(pr.AuthorAssociation)) {
		return nil
	}

	// generate PR comment body
	body := fmt.Sprintf(utils.FirstCommitComment, *(pr.User.Login))
	newComment := &github.IssueComment{
		Body: &body,
	}
	return prp.Client.AddCommentToPR(*(pr.Number), newComment)
}

// isFirstContribution returns true if the author_assiciate field is FIRST_TIME_CONTRIBUTOR.
func isFirstContribution(str string) bool {
	return str == "FIRST_TIME_CONTRIBUTOR"
}
