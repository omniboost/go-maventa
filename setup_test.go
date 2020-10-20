package maventa_test

import (
	"log"
	"net/url"
	"os"
	"testing"

	maventa "github.com/omniboost/go-maventa"
)

var (
	client    *maventa.Client
	companyID int
)

func TestMain(m *testing.M) {
	baseURLString := os.Getenv("BASE_URL")
	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	debug := os.Getenv("DEBUG")

	client = maventa.NewClient(nil, clientID, clientSecret)
	if debug != "" {
		client.SetDebug(true)
	}

	if baseURLString != "" {
		baseURL, err := url.Parse(baseURLString)
		if err != nil {
			log.Fatal(err)
		}
		if baseURL != nil {
			client.SetBaseURL(*baseURL)
			client.SetHTTPClient(client.DefaultClient())
		}
	}
	m.Run()
}
