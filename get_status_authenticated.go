package maventa

import (
	"net/http"
	"net/url"

	"github.com/omniboost/go-maventa/utils"
)

func (c *Client) NewGetStatusAuthenticatedRequest() GetStatusAuthenticatedRequest {
	r := GetStatusAuthenticatedRequest{
		client:  c,
		method:  http.MethodGet,
		headers: http.Header{},
	}

	r.queryParams = r.NewGetStatusAuthenticatedQueryParams()
	r.pathParams = r.NewGetStatusAuthenticatedPathParams()
	r.requestBody = r.NewGetStatusAuthenticatedRequestBody()
	return r
}

type GetStatusAuthenticatedRequest struct {
	client      *Client
	queryParams *GetStatusAuthenticatedQueryParams
	pathParams  *GetStatusAuthenticatedPathParams
	method      string
	headers     http.Header
	requestBody GetStatusAuthenticatedRequestBody
}

func (r GetStatusAuthenticatedRequest) NewGetStatusAuthenticatedQueryParams() *GetStatusAuthenticatedQueryParams {
	return &GetStatusAuthenticatedQueryParams{}
}

type GetStatusAuthenticatedQueryParams struct {
}

func (p GetStatusAuthenticatedQueryParams) ToURLValues() (url.Values, error) {
	encoder := utils.NewSchemaEncoder()
	params := url.Values{}

	err := encoder.Encode(p, params)
	if err != nil {
		return params, err
	}

	return params, nil
}

func (r *GetStatusAuthenticatedRequest) QueryParams() *GetStatusAuthenticatedQueryParams {
	return r.queryParams
}

func (r GetStatusAuthenticatedRequest) NewGetStatusAuthenticatedPathParams() *GetStatusAuthenticatedPathParams {
	return &GetStatusAuthenticatedPathParams{}
}

type GetStatusAuthenticatedPathParams struct {
}

func (p *GetStatusAuthenticatedPathParams) Params() map[string]string {
	return map[string]string{}
}

func (r *GetStatusAuthenticatedRequest) PathParams() *GetStatusAuthenticatedPathParams {
	return r.pathParams
}

func (r *GetStatusAuthenticatedRequest) SetMethod(method string) {
	r.method = method
}

func (r *GetStatusAuthenticatedRequest) Method() string {
	return r.method
}

func (r GetStatusAuthenticatedRequest) NewGetStatusAuthenticatedRequestBody() GetStatusAuthenticatedRequestBody {
	return GetStatusAuthenticatedRequestBody{}
}

type GetStatusAuthenticatedRequestBody struct{}

func (r *GetStatusAuthenticatedRequest) RequestBody() *GetStatusAuthenticatedRequestBody {
	return &r.requestBody
}

func (r *GetStatusAuthenticatedRequest) SetRequestBody(body GetStatusAuthenticatedRequestBody) {
	r.requestBody = body
}

func (r *GetStatusAuthenticatedRequest) NewResponseBody() *GetStatusAuthenticatedResponseBody {
	return &GetStatusAuthenticatedResponseBody{}
}

type GetStatusAuthenticatedResponseBody struct{}

func (r *GetStatusAuthenticatedRequest) URL() url.URL {
	return r.client.GetEndpointURL("/status/authenticated", r.PathParams())
}

func (r *GetStatusAuthenticatedRequest) Do() (GetStatusAuthenticatedResponseBody, error) {
	// Create http request
	req, err := r.client.NewRequest(nil, r.Method(), r.URL(), nil)
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
