package maventa

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"text/template"

	"github.com/omniboost/go-maventa/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

const (
	libraryVersion = "0.0.1"
	userAgent      = "go-maventa/" + libraryVersion
	mediaType      = "application/json"
	charset        = "utf-8"
)

var (
	BaseURL = url.URL{
		Scheme: "https",
		Host:   "ax.maventa.com",
		Path:   "",
	}
	ValidatorURL = url.URL{
		Scheme: "https",
		Host:   "validator.maventa.com",
		Path:   "",
	}
)

// NewClient returns a new Exact Globe Client client
func NewClient(httpClient *http.Client, clientID, clientSecret, vendorAPIKey string) *Client {
	client := &Client{}

	client.SetClientID(clientID)
	client.SetClientSecret(clientSecret)
	client.SetVendorAPIKey(vendorAPIKey)
	client.SetBaseURL(BaseURL)
	client.SetValidatorURL(ValidatorURL)
	client.SetDebug(false)
	client.SetUserAgent(userAgent)
	client.SetMediaType(mediaType)
	client.SetCharset(charset)

	if httpClient == nil {
		httpClient = client.DefaultClient()
	}
	client.SetHTTPClient(httpClient)

	return client
}

// Client manages communication with Exact Globe Client
type Client struct {
	// HTTP client used to communicate with the Client.
	http *http.Client

	debug        bool
	baseURL      url.URL
	validatorURL url.URL

	// credentials
	clientID     string
	clientSecret string
	vendorAPIKey string

	// User agent for client
	userAgent string

	mediaType             string
	charset               string
	disallowUnknownFields bool

	// Optional function called after every successful request made to the DO Clients
	onRequestCompleted RequestCompletionCallback
}

// RequestCompletionCallback defines the type of the request callback function
type RequestCompletionCallback func(*http.Request, *http.Response)

func (c *Client) DefaultClient() *http.Client {
	u := c.GetEndpointURL(c.BaseURL(), "/oauth2/token", AccessTokenPathParams{})

	baseURL := c.BaseURL()
	u2 := baseURL.String()
	oauth2.RegisterBrokenAuthHeaderProvider(u2)
	config := &clientcredentials.Config{
		ClientID:     c.ClientID(),
		ClientSecret: c.ClientSecret(),
		Scopes:       c.Scopes(),
		TokenURL:     u.String(),
		EndpointParams: url.Values(map[string][]string{
			"vendor_api_key": []string{c.VendorAPIKey()},
		}),
	}
	return config.Client(context.Background())
}

func (c *Client) Scopes() []string {
	return []string{
		// "eui",     // Recommended to use when integrating to EUI. Alias for eui:open, company:read, company:write, lookup, receivables:assignments, document:send, document:receive, invoice:receive, invoice:send
		"global",                  // Alias for company:read, document:receive, document:send, lookup
		"company",                 // Alias for company:read, company:write
		"lookup",                  // grants access to the lookup operations
		"document:receive",        // grants access to document receive operations
		"document:send",           // grants access to document send operations
		"invoice:receive",         // grants access to invoice receive operations
		"invoice:send",            // grants access to invoice send operations
		"company:read",            // grants read access to company settings, profiles and notifications
		"company:write",           // grants write access to company settings, profiles and notifications
		"validate",                // grants access to the AutoInvoice validator service
		"receivables:assignments", // grants access to assignments in the receivables service
		"analysis",                // grants access to analysis service
		// "operator:documents:receive",       // grants access to fetch received documents
		// "operator:documents:send",          // grants access to send documents
		// "operator:lookup", // grants access to perform actions related to lookups
		// "operator:participants",            // grants access to perform actions on operator participants
		// "operator:notifications",           // grants access to perform actions on operator notifications
		// "operator:validate", // grants access to the AutoInvoice validator service
		// "operator:receivables:assignments", // grants access to assignments the receivables service
		// "operator:receivables:assignments:create", // grants access to create assignments in the receivables servicee
		// "operator:analysis",                       // grants access to analysis service
	}
}

func (c *Client) SetHTTPClient(client *http.Client) {
	c.http = client
}

func (c Client) Debug() bool {
	return c.debug
}

func (c *Client) SetDebug(debug bool) {
	c.debug = debug
}

func (c Client) ClientID() string {
	return c.clientID
}

func (c *Client) SetClientID(clientID string) {
	c.clientID = clientID
}

func (c Client) ClientSecret() string {
	return c.clientSecret
}

func (c *Client) SetClientSecret(clientSecret string) {
	c.clientSecret = clientSecret
}

func (c Client) VendorAPIKey() string {
	return c.vendorAPIKey
}

func (c *Client) SetVendorAPIKey(vendorAPIKey string) {
	c.vendorAPIKey = vendorAPIKey
}

func (c Client) BaseURL() url.URL {
	return c.baseURL
}

func (c *Client) SetBaseURL(baseURL url.URL) {
	c.baseURL = baseURL
}

func (c Client) ValidatorURL() url.URL {
	return c.validatorURL
}

func (c *Client) SetValidatorURL(validatorURL url.URL) {
	c.validatorURL = validatorURL
}

func (c *Client) SetMediaType(mediaType string) {
	c.mediaType = mediaType
}

func (c Client) MediaType() string {
	return mediaType
}

func (c *Client) SetCharset(charset string) {
	c.charset = charset
}

func (c Client) Charset() string {
	return charset
}

func (c *Client) SetUserAgent(userAgent string) {
	c.userAgent = userAgent
}

func (c Client) UserAgent() string {
	return userAgent
}

func (c *Client) SetDisallowUnknownFields(disallowUnknownFields bool) {
	c.disallowUnknownFields = disallowUnknownFields
}

func (c *Client) GetEndpointURL(base url.URL, path string, pathParams PathParams) url.URL {
	clientURL := base
	clientURL.Path = clientURL.Path + path

	tmpl, err := template.New("endpoint_url").Parse(clientURL.Path)
	if err != nil {
		log.Fatal(err)
	}

	buf := new(bytes.Buffer)
	params := pathParams.Params()
	err = tmpl.Execute(buf, params)
	if err != nil {
		log.Fatal(err)
	}

	clientURL.Path = buf.String()
	return clientURL
}

func (c *Client) NewRequest(ctx context.Context, method string, URL url.URL, body interface{}) (*http.Request, error) {
	// convert body struct to json
	buf := new(bytes.Buffer)
	if body != nil {
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	// create new http request
	req, err := http.NewRequest(method, URL.String(), buf)
	if err != nil {
		return nil, err
	}

	values := url.Values{}
	err = utils.AddURLValuesToRequest(values, req, true)
	if err != nil {
		return nil, err
	}

	// optionally pass along context
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	// set other headers
	req.Header.Add("Content-Type", fmt.Sprintf("%s; charset=%s", c.MediaType(), c.Charset()))
	req.Header.Add("Accept", c.MediaType())
	req.Header.Add("User-Agent", c.UserAgent())

	return req, nil
}

func (c *Client) NewFormRequest(ctx context.Context, method string, URL url.URL, form Form) (*http.Request, error) {
	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)

	for k, vv := range form.Values() {
		for _, v := range vv {
			err := w.WriteField(k, v)
			if err != nil {
				return nil, err
			}
		}
	}

	for k, f := range form.Files() {
		part, err := w.CreateFormFile(k, f.Filename)
		if err != nil {
			return nil, err
		}
		_, err = io.Copy(part, f.Reader)
	}

	err := w.Close()
	if err != nil {
		return nil, err
	}

	// create new http request
	req, err := http.NewRequest(method, URL.String(), body)
	if err != nil {
		return nil, err
	}

	values := url.Values{}
	err = utils.AddURLValuesToRequest(values, req, true)
	if err != nil {
		return nil, err
	}

	// optionally pass along context
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	// set other headers
	req.Header.Add("Content-Type", fmt.Sprintf("%s; charset=%s", w.FormDataContentType(), c.Charset()))
	req.Header.Add("Accept", c.MediaType())
	req.Header.Add("User-Agent", c.UserAgent())

	return req, nil
}

// Do sends an Client request and returns the Client response. The Client response is json decoded and stored in the value
// pointed to by v, or returned as an error if an Client error has occurred. If v implements the io.Writer interface,
// the raw response will be written to v, without attempting to decode it.
func (c *Client) Do(req *http.Request, responseBody interface{}) (*http.Response, error) {
	if c.debug == true {
		dump, _ := httputil.DumpRequestOut(req, true)
		log.Println(string(dump))
	}

	httpResp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}

	if c.onRequestCompleted != nil {
		c.onRequestCompleted(req, httpResp)
	}

	// close body io.Reader
	defer func() {
		if rerr := httpResp.Body.Close(); err == nil {
			err = rerr
		}
	}()

	if c.debug == true {
		dump, _ := httputil.DumpResponse(httpResp, true)
		log.Println(string(dump))
	}

	// check if the response isn't an error
	err = CheckResponse(httpResp)
	if err != nil {
		return httpResp, err
	}

	// check the provided interface parameter
	if httpResp == nil {
		return httpResp, nil
	}

	if responseBody == nil {
		return httpResp, nil
	}

	if httpResp.ContentLength == 0 {
		return httpResp, nil
	}

	// interface implements io.Writer: write Body to it
	// if w, ok := response.Envelope.(io.Writer); ok {
	// 	_, err := io.Copy(w, httpResp.Body)
	// 	return httpResp, err
	// }

	// try to decode body into interface parameter
	if responseBody == nil {
		return httpResp, nil
	}

	apiError := APIError{}
	err = c.Unmarshal(httpResp.Body, &responseBody, &apiError)
	if err != nil {
		return httpResp, err
	}

	if apiError.Error() != "" {
		return httpResp, &ErrorResponse{Response: httpResp, err: apiError}
	}

	return httpResp, nil
}

func (c *Client) Unmarshal(r io.Reader, vv ...interface{}) error {
	if len(vv) == 0 {
		return nil
	}

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	errs := []error{}
	for _, v := range vv {
		r := bytes.NewReader(b)
		dec := json.NewDecoder(r)
		if c.disallowUnknownFields {
			dec.DisallowUnknownFields()
		}

		err := dec.Decode(v)
		if err != nil {
			errs = append(errs, err)
		}

	}

	if len(errs) == len(vv) {
		// Everything errored
		msgs := make([]string, len(errs))
		for i, e := range errs {
			msgs[i] = fmt.Sprint(e)
		}
		return errors.New(strings.Join(msgs, ", "))
	}

	return nil
}

// CheckResponse checks the Client response for errors, and returns them if
// present. A response is considered an error if it has a status code outside
// the 200 range. Client error responses are expected to have either no response
// body, or a json response body that maps to ErrorResponse. Any other response
// body will be silently ignored.
func CheckResponse(r *http.Response) error {
	errorResponse := &ErrorResponse{Response: r}

	// Don't check content-lenght: a created response, for example, has no body
	// if r.Header.Get("Content-Length") == "0" {
	// 	errorResponse.Errors.Message = r.Status
	// 	return errorResponse
	// }

	if c := r.StatusCode; c >= 200 && c <= 299 {
		return nil
	}

	// read data and copy it back
	data, err := ioutil.ReadAll(r.Body)
	r.Body = ioutil.NopCloser(bytes.NewReader(data))
	if err != nil {
		return errorResponse
	}

	err = checkContentType(r)
	if err != nil {
		errorResponse.err = err
		return errorResponse
	}

	if len(data) == 0 {
		errorResponse.err = errors.New("response body is empty")
		return errorResponse
	}

	return nil
}

type ErrorResponse struct {
	// HTTP response that caused this error
	Response *http.Response `json:"-"`

	err error `json:"errors"`
}

func (r ErrorResponse) Error() string {
	if r.Error == nil {
		return ""
	}
	return r.err.Error()
}

// {"code":"invoice_create_api_error","message":"Error while creating invoice","details":["ERROR: INVOICE DATE NOT FOUND"]}
type APIError struct {
	Code    string   `json:"code"`
	Message string   `json:"message"`
	Details []string `json:"details"`
}

func (e APIError) Error() string {
	if e.Code == "" {
		return ""
	}

	return fmt.Sprintf("%s: %s (%s)", e.Code, e.Message, strings.Join(e.Details, ", "))
}

func checkContentType(response *http.Response) error {
	header := response.Header.Get("Content-Type")
	contentType := strings.Split(header, ";")[0]
	if contentType != mediaType && contentType != "application/problem+json" {
		return fmt.Errorf("Expected Content-Type \"%s\", got \"%s\"", mediaType, contentType)
	}

	return nil
}

type PathParams interface {
	Params() map[string]string
}

type AccessTokenPathParams struct{}

func (pp AccessTokenPathParams) Params() map[string]string {
	return map[string]string{}
}
