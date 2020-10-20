package maventa_test

import (
	"encoding/json"
	"log"
	"testing"
)

func TestGetStatusAuthenticated(t *testing.T) {
	req := client.NewGetStatusAuthenticatedRequest()
	resp, err := req.Do()
	if err != nil {
		t.Error(err)
	}

	b, _ := json.MarshalIndent(resp, "", "  ")
	log.Println(string(b))
}
