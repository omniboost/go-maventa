package maventa

import (
	"net/http"
	"net/url"

	"github.com/omniboost/go-maventa/utils"
)

func (c *Client) NewPostInvoiceRequest() PostInvoiceRequest {
	r := PostInvoiceRequest{
		client:  c,
		method:  http.MethodPost,
		headers: http.Header{},
	}

	r.queryParams = r.NewPostInvoiceQueryParams()
	r.pathParams = r.NewPostInvoicePathParams()
	r.formParams = r.NewPostInvoiceFormParams()
	r.requestBody = r.NewPostInvoiceRequestBody()
	return r
}

type PostInvoiceRequest struct {
	client      *Client
	queryParams *PostInvoiceQueryParams
	pathParams  *PostInvoicePathParams
	formParams  *PostInvoiceFormParams
	method      string
	headers     http.Header
	requestBody PostInvoiceRequestBody
}

func (r PostInvoiceRequest) NewPostInvoiceQueryParams() *PostInvoiceQueryParams {
	return &PostInvoiceQueryParams{}
}

type PostInvoiceQueryParams struct {
	Direction string `schema:"direction"`
}

func (p PostInvoiceQueryParams) ToURLValues() (url.Values, error) {
	encoder := utils.NewSchemaEncoder()
	params := url.Values{}

	err := encoder.Encode(p, params)
	if err != nil {
		return params, err
	}

	return params, nil
}

func (r *PostInvoiceRequest) QueryParams() *PostInvoiceQueryParams {
	return r.queryParams
}

func (r *PostInvoiceRequest) FormParams() *PostInvoiceFormParams {
	return r.formParams
}

func (r PostInvoiceRequest) NewPostInvoicePathParams() *PostInvoicePathParams {
	return &PostInvoicePathParams{}
}

type PostInvoicePathParams struct {
}

func (p *PostInvoicePathParams) Params() map[string]string {
	return map[string]string{}
}

func (r *PostInvoiceRequest) PathParams() *PostInvoicePathParams {
	return r.pathParams
}

func (r PostInvoiceRequest) NewPostInvoiceFormParams() *PostInvoiceFormParams {
	return &PostInvoiceFormParams{}
}

type PostInvoiceFormParams struct {
	File              File
	Format            string
	RecipientType     string
	RecipientEIA      string
	RecipientEmail    string
	RecipientOperator string
	DisabledRoutes    []string
	SenderComment     string
	RouteOrder        []string
}

func (p PostInvoiceFormParams) Values() url.Values {
	return url.Values{
		"format":             []string{p.Format},
		"recipient_type":     []string{p.RecipientType},
		"recipient_eia":      []string{p.RecipientEIA},
		"recipient_email":    []string{p.RecipientEmail},
		"recipient_operator": []string{p.RecipientOperator},
		"disabled_routes":    p.DisabledRoutes,
		"sender_comment":     []string{p.SenderComment},
		"route_order":        p.RouteOrder,
	}
}

func (p PostInvoiceFormParams) Files() map[string]File {
	return map[string]File{
		"file": p.File,
	}
}

func (r *PostInvoiceRequest) SetMethod(method string) {
	r.method = method
}

func (r *PostInvoiceRequest) Method() string {
	return r.method
}

func (r PostInvoiceRequest) NewPostInvoiceRequestBody() PostInvoiceRequestBody {
	return PostInvoiceRequestBody{}
}

type PostInvoiceRequestBody struct{}

func (r *PostInvoiceRequest) RequestBody() *PostInvoiceRequestBody {
	return &r.requestBody
}

func (r *PostInvoiceRequest) SetRequestBody(body PostInvoiceRequestBody) {
	r.requestBody = body
}

func (r *PostInvoiceRequest) NewResponseBody() *PostInvoiceResponseBody {
	return &PostInvoiceResponseBody{}
}

type PostInvoiceResponseBody struct {
}

func (r *PostInvoiceRequest) URL() url.URL {
	return r.client.GetEndpointURL("/v1/invoices", r.PathParams())
}

func (r *PostInvoiceRequest) Do() (PostInvoiceResponseBody, error) {
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
