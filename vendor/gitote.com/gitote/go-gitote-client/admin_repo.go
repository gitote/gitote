package gitote

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func (c *Client) AdminCreateRepo(user string, opt CreateRepoOption) (*Repository, error) {
	body, err := json.Marshal(&opt)
	if err != nil {
		return nil, err
	}
	repo := new(Repository)
	return repo, c.getParsedResponse("POST", fmt.Sprintf("/admin/users/%s/repos", user),
		jsonHeader, bytes.NewReader(body), repo)
}
