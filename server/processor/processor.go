package processor

import (
	"fmt"

	"github.com/pouchcontainer/pouchrobot/server/processor/issueCommentProcessor"
	"github.com/pouchcontainer/pouchrobot/server/processor/issueProcessor"
	"github.com/pouchcontainer/pouchrobot/server/processor/prCommentProcessor"
	"github.com/pouchcontainer/pouchrobot/server/processor/pullRequestProcessor"
	"github.com/pouchcontainer/pouchrobot/server/utils"

	"github.com/pouchcontainer/pouchrobot/server/gh"
	"github.com/sirupsen/logrus"
)

type processor interface {
	// Process processes item automan gets, and then execute operations torwards items on GitHub
	Process(data []byte) error
}

// Processor contains several specific processors
type Processor struct {
	IssueProcessor        *issueProcessor.IssueProcessor
	PullRequestProcessor  *pullRequestProcessor.PullRequestProcessor
	IssueCommentProcessor *issueCommentProcessor.IssueCommentProcessor
	PRCommentProcessor    *prCommentProcessor.PRCommentProcessor
}

// New initializes a brand new processor.
func New(client *gh.Client) *Processor {
	return &Processor{
		IssueProcessor: &issueProcessor.IssueProcessor{
			Client: client,
		},
		PullRequestProcessor: &pullRequestProcessor.PullRequestProcessor{
			Client: client,
		},
		IssueCommentProcessor: &issueCommentProcessor.IssueCommentProcessor{
			Client: client,
		},
		PRCommentProcessor: &prCommentProcessor.PRCommentProcessor{
			Client: client,
		},
	}
}

// HandleEvent processes an event received from github
func (p *Processor) HandleEvent(eventType string, data []byte) error {
	switch eventType {
	case "issues":
		p.IssueProcessor.Process(data)
	case "pull_request":
		p.PullRequestProcessor.Process(data)
	case "issue_comment":
		// since pr is also a kind of issue, we need to first make it clear
		issueType := judgeIssueOrPR(data)
		logrus.Infof("get issueType: %s", issueType)
		if issueType == "issue" {
			p.IssueCommentProcessor.Process(data)
			return nil
		}
		if issueType == "pull_request" {
			p.PRCommentProcessor.Process(data)
			return nil
		}
		return nil
	default:
		return fmt.Errorf("unknown event type %s", eventType)
	}
	return nil
}

func judgeIssueOrPR(data []byte) string {
	issue, err := utils.ExactIssue(data)
	if err != nil {
		return ""
	}

	if issue.PullRequestLinks == nil {
		return "issue"
	}
	return "pull_request"
}
