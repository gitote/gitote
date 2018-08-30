package models

import (
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"time"

	gouuid "github.com/satori/go.uuid"
	"gitote.com/yoginth/com"
	log "gopkg.in/clog.v1"

	git "gitote.com/gitote/git-module"

	"gitote.com/gitote/gitote/models/errors"
	"gitote.com/gitote/gitote/pkg/process"
	"gitote.com/gitote/gitote/pkg/setting"
)

// discardLocalRepoBranchChanges discards local commits/changes of
// given branch to make sure it is even to remote branch.
func discardLocalRepoBranchChanges(localPath, branch string) error {
	if !com.IsExist(localPath) {
		return nil
	}
	// No need to check if nothing in the repository.
	if !git.IsBranchExist(localPath, branch) {
		return nil
	}

	refName := "origin/" + branch
	if err := git.ResetHEAD(localPath, true, refName); err != nil {
		return fmt.Errorf("git reset --hard %s: %v", refName, err)
	}
	return nil
}

func (repo *Repository) DiscardLocalRepoBranchChanges(branch string) error {
	return discardLocalRepoBranchChanges(repo.LocalCopyPath(), branch)
}

// checkoutNewBranch checks out to a new branch from the a branch name.
func checkoutNewBranch(repoPath, localPath, oldBranch, newBranch string) error {
	if err := git.Checkout(localPath, git.CheckoutOptions{
		Timeout:   time.Duration(setting.Git.Timeout.Pull) * time.Second,
		Branch:    newBranch,
		OldBranch: oldBranch,
	}); err != nil {
		return fmt.Errorf("git checkout -b %s %s: %v", newBranch, oldBranch, err)
	}
	return nil
}

func (repo *Repository) CheckoutNewBranch(oldBranch, newBranch string) error {
	return checkoutNewBranch(repo.RepoPath(), repo.LocalCopyPath(), oldBranch, newBranch)
}

type UpdateRepoFileOptions struct {
	LastCommitID string
	OldBranch    string
	NewBranch    string
	OldTreeName  string
	NewTreeName  string
	Message      string
	Content      string
	IsNewFile    bool
}

// UpdateRepoFile adds or updates a file in repository.
func (repo *Repository) UpdateRepoFile(doer *User, opts UpdateRepoFileOptions) (err error) {
	repoWorkingPool.CheckIn(com.ToStr(repo.ID))
	defer repoWorkingPool.CheckOut(com.ToStr(repo.ID))

	if err = repo.DiscardLocalRepoBranchChanges(opts.OldBranch); err != nil {
		return fmt.Errorf("DiscardLocalRepoBranchChanges [branch: %s]: %v", opts.OldBranch, err)
	} else if err = repo.UpdateLocalCopyBranch(opts.OldBranch); err != nil {
		return fmt.Errorf("UpdateLocalCopyBranch [branch: %s]: %v", opts.OldBranch, err)
	}

	repoPath := repo.RepoPath()
	localPath := repo.LocalCopyPath()

	if opts.OldBranch != opts.NewBranch {
		// Directly return error if new branch already exists in the server
		if git.IsBranchExist(repoPath, opts.NewBranch) {
			return errors.BranchAlreadyExists{opts.NewBranch}
		}

		// Otherwise, delete branch from local copy in case out of sync
		if git.IsBranchExist(localPath, opts.NewBranch) {
			if err = git.DeleteBranch(localPath, opts.NewBranch, git.DeleteBranchOptions{
				Force: true,
			}); err != nil {
				return fmt.Errorf("DeleteBranch [name: %s]: %v", opts.NewBranch, err)
			}
		}

		if err := repo.CheckoutNewBranch(opts.OldBranch, opts.NewBranch); err != nil {
			return fmt.Errorf("CheckoutNewBranch [old_branch: %s, new_branch: %s]: %v", opts.OldBranch, opts.NewBranch, err)
		}
	}

	oldFilePath := path.Join(localPath, opts.OldTreeName)
	filePath := path.Join(localPath, opts.NewTreeName)
	os.MkdirAll(path.Dir(filePath), os.ModePerm)

	// If it's meant to be a new file, make sure it doesn't exist.
	if opts.IsNewFile {
		if com.IsExist(filePath) {
			return ErrRepoFileAlreadyExist{filePath}
		}
	}

	// Ignore move step if it's a new file under a directory.
	// Otherwise, move the file when name changed.
	if com.IsFile(oldFilePath) && opts.OldTreeName != opts.NewTreeName {
		if err = git.MoveFile(localPath, opts.OldTreeName, opts.NewTreeName); err != nil {
			return fmt.Errorf("git mv %s %s: %v", opts.OldTreeName, opts.NewTreeName, err)
		}
	}

	if err = ioutil.WriteFile(filePath, []byte(opts.Content), 0666); err != nil {
		return fmt.Errorf("WriteFile: %v", err)
	}

	if err = git.AddChanges(localPath, true); err != nil {
		return fmt.Errorf("git add --all: %v", err)
	} else if err = git.CommitChanges(localPath, git.CommitChangesOptions{
		Committer: doer.NewGitSig(),
		Message:   opts.Message,
	}); err != nil {
		return fmt.Errorf("CommitChanges: %v", err)
	} else if err = git.Push(localPath, "origin", opts.NewBranch); err != nil {
		return fmt.Errorf("git push origin %s: %v", opts.NewBranch, err)
	}

	gitRepo, err := git.OpenRepository(repo.RepoPath())
	if err != nil {
		log.Error(2, "OpenRepository: %v", err)
		return nil
	}
	commit, err := gitRepo.GetBranchCommit(opts.NewBranch)
	if err != nil {
		log.Error(2, "GetBranchCommit [branch: %s]: %v", opts.NewBranch, err)
		return nil
	}

	// Simulate push event.
	pushCommits := &PushCommits{
		Len:     1,
		Commits: []*PushCommit{CommitToPushCommit(commit)},
	}
	oldCommitID := opts.LastCommitID
	if opts.NewBranch != opts.OldBranch {
		oldCommitID = git.EMPTY_SHA
	}
	if err := CommitRepoAction(CommitRepoActionOptions{
		PusherName:  doer.Name,
		RepoOwnerID: repo.MustOwner().ID,
		RepoName:    repo.Name,
		RefFullName: git.BRANCH_PREFIX + opts.NewBranch,
		OldCommitID: oldCommitID,
		NewCommitID: commit.ID.String(),
		Commits:     pushCommits,
	}); err != nil {
		log.Error(2, "CommitRepoAction: %v", err)
		return nil
	}

	go AddTestPullRequestTask(doer, repo.ID, opts.NewBranch, true)
	return nil
}

// GetDiffPreview produces and returns diff result of a file which is not yet committed.
func (repo *Repository) GetDiffPreview(branch, treePath, content string) (diff *Diff, err error) {
	repoWorkingPool.CheckIn(com.ToStr(repo.ID))
	defer repoWorkingPool.CheckOut(com.ToStr(repo.ID))

	if err = repo.DiscardLocalRepoBranchChanges(branch); err != nil {
		return nil, fmt.Errorf("DiscardLocalRepoBranchChanges [branch: %s]: %v", branch, err)
	} else if err = repo.UpdateLocalCopyBranch(branch); err != nil {
		return nil, fmt.Errorf("UpdateLocalCopyBranch [branch: %s]: %v", branch, err)
	}

	localPath := repo.LocalCopyPath()
	filePath := path.Join(localPath, treePath)
	os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
	if err = ioutil.WriteFile(filePath, []byte(content), 0666); err != nil {
		return nil, fmt.Errorf("WriteFile: %v", err)
	}

	cmd := exec.Command("git", "diff", treePath)
	cmd.Dir = localPath
	cmd.Stderr = os.Stderr

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("StdoutPipe: %v", err)
	}

	if err = cmd.Start(); err != nil {
		return nil, fmt.Errorf("Start: %v", err)
	}

	pid := process.Add(fmt.Sprintf("GetDiffPreview [repo_path: %s]", repo.RepoPath()), cmd)
	defer process.Remove(pid)

	diff, err = ParsePatch(setting.Git.MaxGitDiffLines, setting.Git.MaxGitDiffLineCharacters, setting.Git.MaxGitDiffFiles, stdout)
	if err != nil {
		return nil, fmt.Errorf("ParsePatch: %v", err)
	}

	if err = cmd.Wait(); err != nil {
		return nil, fmt.Errorf("Wait: %v", err)
	}

	return diff, nil
}

type DeleteRepoFileOptions struct {
	LastCommitID string
	OldBranch    string
	NewBranch    string
	TreePath     string
	Message      string
}

func (repo *Repository) DeleteRepoFile(doer *User, opts DeleteRepoFileOptions) (err error) {
	repoWorkingPool.CheckIn(com.ToStr(repo.ID))
	defer repoWorkingPool.CheckOut(com.ToStr(repo.ID))

	if err = repo.DiscardLocalRepoBranchChanges(opts.OldBranch); err != nil {
		return fmt.Errorf("DiscardLocalRepoBranchChanges [branch: %s]: %v", opts.OldBranch, err)
	} else if err = repo.UpdateLocalCopyBranch(opts.OldBranch); err != nil {
		return fmt.Errorf("UpdateLocalCopyBranch [branch: %s]: %v", opts.OldBranch, err)
	}

	if opts.OldBranch != opts.NewBranch {
		if err := repo.CheckoutNewBranch(opts.OldBranch, opts.NewBranch); err != nil {
			return fmt.Errorf("CheckoutNewBranch [old_branch: %s, new_branch: %s]: %v", opts.OldBranch, opts.NewBranch, err)
		}
	}

	localPath := repo.LocalCopyPath()
	if err = os.Remove(path.Join(localPath, opts.TreePath)); err != nil {
		return fmt.Errorf("Remove: %v", err)
	}

	if err = git.AddChanges(localPath, true); err != nil {
		return fmt.Errorf("git add --all: %v", err)
	} else if err = git.CommitChanges(localPath, git.CommitChangesOptions{
		Committer: doer.NewGitSig(),
		Message:   opts.Message,
	}); err != nil {
		return fmt.Errorf("CommitChanges: %v", err)
	} else if err = git.Push(localPath, "origin", opts.NewBranch); err != nil {
		return fmt.Errorf("git push origin %s: %v", opts.NewBranch, err)
	}

	gitRepo, err := git.OpenRepository(repo.RepoPath())
	if err != nil {
		log.Error(2, "OpenRepository: %v", err)
		return nil
	}
	commit, err := gitRepo.GetBranchCommit(opts.NewBranch)
	if err != nil {
		log.Error(2, "GetBranchCommit [branch: %s]: %v", opts.NewBranch, err)
		return nil
	}

	// Simulate push event.
	pushCommits := &PushCommits{
		Len:     1,
		Commits: []*PushCommit{CommitToPushCommit(commit)},
	}
	if err := CommitRepoAction(CommitRepoActionOptions{
		PusherName:  doer.Name,
		RepoOwnerID: repo.MustOwner().ID,
		RepoName:    repo.Name,
		RefFullName: git.BRANCH_PREFIX + opts.NewBranch,
		OldCommitID: opts.LastCommitID,
		NewCommitID: commit.ID.String(),
		Commits:     pushCommits,
	}); err != nil {
		log.Error(2, "CommitRepoAction: %v", err)
		return nil
	}

	go AddTestPullRequestTask(doer, repo.ID, opts.NewBranch, true)
	return nil
}

// Upload represent a uploaded file to a repo to be deleted when moved
type Upload struct {
	ID   int64
	UUID string `xorm:"uuid UNIQUE"`
	Name string
}

// UploadLocalPath returns where uploads is stored in local file system based on given UUID.
func UploadLocalPath(uuid string) string {
	return path.Join(setting.Repository.Upload.TempPath, uuid[0:1], uuid[1:2], uuid)
}

// LocalPath returns where uploads are temporarily stored in local file system.
func (upload *Upload) LocalPath() string {
	return UploadLocalPath(upload.UUID)
}

// NewUpload creates a new upload object.
func NewUpload(name string, buf []byte, file multipart.File) (_ *Upload, err error) {
	upload := &Upload{
		UUID: gouuid.NewV4().String(),
		Name: name,
	}

	localPath := upload.LocalPath()
	if err = os.MkdirAll(path.Dir(localPath), os.ModePerm); err != nil {
		return nil, fmt.Errorf("MkdirAll: %v", err)
	}

	fw, err := os.Create(localPath)
	if err != nil {
		return nil, fmt.Errorf("Create: %v", err)
	}
	defer fw.Close()

	if _, err = fw.Write(buf); err != nil {
		return nil, fmt.Errorf("Write: %v", err)
	} else if _, err = io.Copy(fw, file); err != nil {
		return nil, fmt.Errorf("Copy: %v", err)
	}

	if _, err := x.Insert(upload); err != nil {
		return nil, err
	}

	return upload, nil
}

func GetUploadByUUID(uuid string) (*Upload, error) {
	upload := &Upload{UUID: uuid}
	has, err := x.Get(upload)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, ErrUploadNotExist{0, uuid}
	}
	return upload, nil
}

func GetUploadsByUUIDs(uuids []string) ([]*Upload, error) {
	if len(uuids) == 0 {
		return []*Upload{}, nil
	}

	// Silently drop invalid uuids.
	uploads := make([]*Upload, 0, len(uuids))
	return uploads, x.In("uuid", uuids).Find(&uploads)
}

func DeleteUploads(uploads ...*Upload) (err error) {
	if len(uploads) == 0 {
		return nil
	}

	sess := x.NewSession()
	defer sess.Close()
	if err = sess.Begin(); err != nil {
		return err
	}

	ids := make([]int64, len(uploads))
	for i := 0; i < len(uploads); i++ {
		ids[i] = uploads[i].ID
	}
	if _, err = sess.In("id", ids).Delete(new(Upload)); err != nil {
		return fmt.Errorf("delete uploads: %v", err)
	}

	for _, upload := range uploads {
		localPath := upload.LocalPath()
		if !com.IsFile(localPath) {
			continue
		}

		if err := os.Remove(localPath); err != nil {
			return fmt.Errorf("remove upload: %v", err)
		}
	}

	return sess.Commit()
}

func DeleteUpload(u *Upload) error {
	return DeleteUploads(u)
}

func DeleteUploadByUUID(uuid string) error {
	upload, err := GetUploadByUUID(uuid)
	if err != nil {
		if IsErrUploadNotExist(err) {
			return nil
		}
		return fmt.Errorf("GetUploadByUUID: %v", err)
	}

	if err := DeleteUpload(upload); err != nil {
		return fmt.Errorf("DeleteUpload: %v", err)
	}

	return nil
}

type UploadRepoFileOptions struct {
	LastCommitID string
	OldBranch    string
	NewBranch    string
	TreePath     string
	Message      string
	Files        []string // In UUID format.
}

func (repo *Repository) UploadRepoFiles(doer *User, opts UploadRepoFileOptions) (err error) {
	if len(opts.Files) == 0 {
		return nil
	}

	uploads, err := GetUploadsByUUIDs(opts.Files)
	if err != nil {
		return fmt.Errorf("GetUploadsByUUIDs [uuids: %v]: %v", opts.Files, err)
	}

	repoWorkingPool.CheckIn(com.ToStr(repo.ID))
	defer repoWorkingPool.CheckOut(com.ToStr(repo.ID))

	if err = repo.DiscardLocalRepoBranchChanges(opts.OldBranch); err != nil {
		return fmt.Errorf("DiscardLocalRepoBranchChanges [branch: %s]: %v", opts.OldBranch, err)
	} else if err = repo.UpdateLocalCopyBranch(opts.OldBranch); err != nil {
		return fmt.Errorf("UpdateLocalCopyBranch [branch: %s]: %v", opts.OldBranch, err)
	}

	if opts.OldBranch != opts.NewBranch {
		if err = repo.CheckoutNewBranch(opts.OldBranch, opts.NewBranch); err != nil {
			return fmt.Errorf("CheckoutNewBranch [old_branch: %s, new_branch: %s]: %v", opts.OldBranch, opts.NewBranch, err)
		}
	}

	localPath := repo.LocalCopyPath()
	dirPath := path.Join(localPath, opts.TreePath)
	os.MkdirAll(dirPath, os.ModePerm)

	// Copy uploaded files into repository.
	for _, upload := range uploads {
		tmpPath := upload.LocalPath()
		targetPath := path.Join(dirPath, upload.Name)
		if !com.IsFile(tmpPath) {
			continue
		}

		if err = com.Copy(tmpPath, targetPath); err != nil {
			return fmt.Errorf("Copy: %v", err)
		}
	}

	if err = git.AddChanges(localPath, true); err != nil {
		return fmt.Errorf("git add --all: %v", err)
	} else if err = git.CommitChanges(localPath, git.CommitChangesOptions{
		Committer: doer.NewGitSig(),
		Message:   opts.Message,
	}); err != nil {
		return fmt.Errorf("CommitChanges: %v", err)
	} else if err = git.Push(localPath, "origin", opts.NewBranch); err != nil {
		return fmt.Errorf("git push origin %s: %v", opts.NewBranch, err)
	}

	gitRepo, err := git.OpenRepository(repo.RepoPath())
	if err != nil {
		log.Error(2, "OpenRepository: %v", err)
		return nil
	}
	commit, err := gitRepo.GetBranchCommit(opts.NewBranch)
	if err != nil {
		log.Error(2, "GetBranchCommit [branch: %s]: %v", opts.NewBranch, err)
		return nil
	}

	// Simulate push event.
	pushCommits := &PushCommits{
		Len:     1,
		Commits: []*PushCommit{CommitToPushCommit(commit)},
	}
	if err := CommitRepoAction(CommitRepoActionOptions{
		PusherName:  doer.Name,
		RepoOwnerID: repo.MustOwner().ID,
		RepoName:    repo.Name,
		RefFullName: git.BRANCH_PREFIX + opts.NewBranch,
		OldCommitID: opts.LastCommitID,
		NewCommitID: commit.ID.String(),
		Commits:     pushCommits,
	}); err != nil {
		log.Error(2, "CommitRepoAction: %v", err)
		return nil
	}

	go AddTestPullRequestTask(doer, repo.ID, opts.NewBranch, true)
	return DeleteUploads(uploads...)
}
