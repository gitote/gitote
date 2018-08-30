package repo

import (
	"gitote.com/gitote/git-module"

	"gitote.com/gitote/gitote/models"
	"gitote.com/gitote/gitote/pkg/context"
	"gitote.com/gitote/gitote/routes/repo"
)

func GetRawFile(c *context.APIContext) {
	if !c.Repo.HasAccess() {
		c.Status(404)
		return
	}

	if c.Repo.Repository.IsBare {
		c.Status(404)
		return
	}

	blob, err := c.Repo.Commit.GetBlobByPath(c.Repo.TreePath)
	if err != nil {
		if git.IsErrNotExist(err) {
			c.Status(404)
		} else {
			c.Error(500, "GetBlobByPath", err)
		}
		return
	}
	if err = repo.ServeBlob(c.Context, blob); err != nil {
		c.Error(500, "ServeBlob", err)
	}
}

func GetArchive(c *context.APIContext) {
	repoPath := models.RepoPath(c.Params(":username"), c.Params(":reponame"))
	gitRepo, err := git.OpenRepository(repoPath)
	if err != nil {
		c.Error(500, "OpenRepository", err)
		return
	}
	c.Repo.GitRepo = gitRepo

	repo.Download(c.Context)
}

func GetEditorconfig(c *context.APIContext) {
	ec, err := c.Repo.GetEditorconfig()
	if err != nil {
		if git.IsErrNotExist(err) {
			c.Error(404, "GetEditorconfig", err)
		} else {
			c.Error(500, "GetEditorconfig", err)
		}
		return
	}

	fileName := c.Params("filename")
	def := ec.GetDefinitionForFilename(fileName)
	if def == nil {
		c.Error(404, "GetDefinitionForFilename", err)
		return
	}
	c.JSON(200, def)
}
