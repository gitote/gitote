package repo

import (
	api "gitote.com/gitote/go-gitote-client"

	"gitote.com/gitote/gitote/models/errors"
	"gitote.com/gitote/gitote/pkg/context"
	"gitote.com/gitote/gitote/routes/api/v1/convert"
)

func GetBranch(c *context.APIContext) {
	branch, err := c.Repo.Repository.GetBranch(c.Params("*"))
	if err != nil {
		if errors.IsErrBranchNotExist(err) {
			c.Error(404, "GetBranch", err)
		} else {
			c.Error(500, "GetBranch", err)
		}
		return
	}

	commit, err := branch.GetCommit()
	if err != nil {
		c.Error(500, "GetCommit", err)
		return
	}

	c.JSON(200, convert.ToBranch(branch, commit))
}

func ListBranches(c *context.APIContext) {
	branches, err := c.Repo.Repository.GetBranches()
	if err != nil {
		c.Error(500, "GetBranches", err)
		return
	}

	apiBranches := make([]*api.Branch, len(branches))
	for i := range branches {
		commit, err := branches[i].GetCommit()
		if err != nil {
			c.Error(500, "GetCommit", err)
			return
		}
		apiBranches[i] = convert.ToBranch(branches[i], commit)
	}

	c.JSON(200, &apiBranches)
}
