package maventa

import (
	"net/http"
	"net/url"

	"github.com/omniboost/go-maventa/utils"
)

func (c *Client) NewValidateRequest() ValidateRequest {
	r := ValidateRequest{
		client:  c,
		method:  http.MethodPost,
		headers: http.Header{},
	}

	r.queryParams = r.NewValidateQueryParams()
	r.pathParams = r.NewValidatePathParams()
	r.formParams = r.NewValidateFormParams()
	r.requestBody = r.NewValidateRequestBody()
	return r
}

type ValidateRequest struct {
	client      *Client
	queryParams *ValidateQueryParams
	pathParams  *ValidatePathParams
	formParams  *ValidateFormParams
	method      string
	headers     http.Header
	requestBody ValidateRequestBody
}

func (r ValidateRequest) NewValidateQueryParams() *ValidateQueryParams {
	return &ValidateQueryParams{}
}

type ValidateQueryParams struct {
}

func (p ValidateQueryParams) ToURLValues() (url.Values, error) {
	encoder := utils.NewSchemaEncoder()
	params := url.Values{}

	err := encoder.Encode(p, params)
	if err != nil {
		return params, err
	}

	return params, nil
}

func (r *ValidateRequest) QueryParams() *ValidateQueryParams {
	return r.queryParams
}

func (r *ValidateRequest) FormParams() *ValidateFormParams {
	return r.formParams
}

func (r ValidateRequest) NewValidatePathParams() *ValidatePathParams {
	return &ValidatePathParams{}
}

type ValidatePathParams struct {
}

func (p *ValidatePathParams) Params() map[string]string {
	return map[string]string{}
}

func (r *ValidateRequest) PathParams() *ValidatePathParams {
	return r.pathParams
}

func (r ValidateRequest) NewValidateFormParams() *ValidateFormParams {
	return &ValidateFormParams{}
}

type ValidateFormParams struct {
	File File
}

func (p ValidateFormParams) Values() url.Values {
	return url.Values{}
}

func (p ValidateFormParams) Files() map[string]File {
	return map[string]File{
		"file": p.File,
	}
}

func (r *ValidateRequest) SetMethod(method string) {
	r.method = method
}

func (r *ValidateRequest) Method() string {
	return r.method
}

func (r ValidateRequest) NewValidateRequestBody() ValidateRequestBody {
	return ValidateRequestBody{}
}

type ValidateRequestBody struct{}

func (r *ValidateRequest) RequestBody() *ValidateRequestBody {
	return &r.requestBody
}

func (r *ValidateRequest) SetRequestBody(body ValidateRequestBody) {
	r.requestBody = body
}

func (r *ValidateRequest) NewResponseBody() *ValidateResponseBody {
	return &ValidateResponseBody{}
}

type ValidateResponseBody struct {
}

func (r *ValidateRequest) URL() url.URL {
	return r.client.GetEndpointURL(r.client.ValidatorURL(), "/v1/validate", r.PathParams())
}

func (r *ValidateRequest) Do() (ValidateResponseBody, error) {
	// Create http request
	req, err := r.client.NewFormRequest(nil, r.Method(), r.URL(), r.FormParams())
	if err != nil {
		return *r.NewResponseBody(), err
	}

	// Process query parameters
	err = utils.AddQueryParamsToRequest(r.QueryParams(), req, false)
	if err != nil {
		return *r.NewResponseBody(), err
	}

	responseBody := r.NewResponseBody()
	_, err = r.client.Do(req, responseBody)
	return *responseBody, err
}
