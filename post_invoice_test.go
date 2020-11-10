package maventa_test

import (
	"encoding/json"
	"log"
	"os"
	"testing"

	maventa "github.com/omniboost/go-maventa"
)

func TestPostInvoice(t *testing.T) {
	f, err := os.Open("test.xml")
	if err != nil {
		t.Error(err)
	}
	defer f.Close()

	req := client.NewPostInvoiceRequest()
	req.FormParams().Format = "FINVOICE30"
	req.FormParams().RecipientType = "consumer"
	req.FormParams().File = maventa.File{
		Filename: "test.xml",
		Reader:   f,
	}
	resp, err := req.Do()
	if err != nil {
		t.Error(err)
	}

	b, _ := json.MarshalIndent(resp, "", "  ")
	log.Println(string(b))
}
