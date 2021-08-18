package looker

import (
	"bytes"
	"encoding/json"
	"fmt"

	apiclient "github.com/looker-open-source/sdk-codegen/go/sdk/v3"
)

var (
	VALID_WORKSPACES = []string{DEV_WORKSPACE, PROD_WORKSPACE}
)

const (
	DEV_WORKSPACE  = "dev"
	PROD_WORKSPACE = "production"
)

func selectAPISession(client *apiclient.LookerSDK, id string) error {
	if id != DEV_WORKSPACE && id != PROD_WORKSPACE {
		return fmt.Errorf("illegal value for workspace: %+v", id)
	}

	updateSessionBody := apiclient.WriteApiSession{
		WorkspaceId: &id,
	}
	_, err := client.UpdateSession(updateSessionBody, nil)

	return err
}

func JSONMarshal(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(t)
	return buffer.Bytes(), err
}
