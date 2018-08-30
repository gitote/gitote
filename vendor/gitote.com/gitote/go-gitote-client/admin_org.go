package gitote

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func (c *Client) AdminCreateOrg(user string, opt CreateOrgOption) (*Organization, error) {
	body, err := json.Marshal(&opt)
	if err != nil {
		return nil, err
	}
	org := new(Organization)
	return org, c.getParsedResponse("POST", fmt.Sprintf("/admin/users/%s/orgs", user),
		jsonHeader, bytes.NewReader(body), org)
}
