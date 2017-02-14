package roles

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/crowdflux/angel/app/models/acl_type"
	"github.com/crowdflux/angel/app/plog"
	"github.com/crowdflux/angel/app/plog/log_tags"
	"github.com/pkg/errors"
	"net/http"
)

//This will make a http request to Haimdall to validate from ACL
func ValidateRequest(token string, roles []string, url string) (bool, error) {
	body, err := json.Marshal(acl_type.ACL_Request{
		token,
		roles,
	})
	if err != nil {
		plog.Error("ACL", err)
		return false, err
	}
	fmt.Println("creating request")
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		plog.Error("ACL", err)
		return false, err
	}
	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		plog.Error("ACL", err)
		return false, err
	}
	var acl_response acl_type.ACL_Response
	err = json.NewDecoder(res.Body).Decode(&acl_response)
	if err != nil {
		plog.Error("ACL", err)
		return false, err
	}
	fmt.Println("response", acl_response)
	if acl_response.Success {
		if acl_response.Response.IsPermitted {
			return true, nil
		} else {
			return false, nil
		}
	} else {
		plog.Error("ACL", errors.New(acl_response.Error.Message), plog.MP(log_tags.ERROR_CODE, acl_response.Error.ErrorCode))
		return false, errors.New(acl_response.Error.Message)
	}
}
