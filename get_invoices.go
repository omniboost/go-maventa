package maventa

import (
	"net/http"
	"net/url"
	"time"

	"github.com/omniboost/go-maventa/utils"
)

func (c *Client) NewGetInvoicesRequest() GetInvoicesRequest {
	r := GetInvoicesRequest{
		client:  c,
		method:  http.MethodGet,
		headers: http.Header{},
	}

	r.queryParams = r.NewGetInvoicesQueryParams()
	r.pathParams = r.NewGetInvoicesPathParams()
	r.requestBody = r.NewGetInvoicesRequestBody()
	return r
}

type GetInvoicesRequest struct {
	client      *Client
	queryParams *GetInvoicesQueryParams
	pathParams  *GetInvoicesPathParams
	method      string
	headers     http.Header
	requestBody GetInvoicesRequestBody
}

func (r GetInvoicesRequest) NewGetInvoicesQueryParams() *GetInvoicesQueryParams {
	return &GetInvoicesQueryParams{}
}

type GetInvoicesQueryParams struct {
	Direction string `schema:"direction"`
}

func (p GetInvoicesQueryParams) ToURLValues() (url.Values, error) {
	encoder := utils.NewSchemaEncoder()
	params := url.Values{}

	err := encoder.Encode(p, params)
	if err != nil {
		return params, err
	}

	return params, nil
}

func (r *GetInvoicesRequest) QueryParams() *GetInvoicesQueryParams {
	return r.queryParams
}

func (r GetInvoicesRequest) NewGetInvoicesPathParams() *GetInvoicesPathParams {
	return &GetInvoicesPathParams{}
}

type GetInvoicesPathParams struct {
}

func (p *GetInvoicesPathParams) Params() map[string]string {
	return map[string]string{}
}

func (r *GetInvoicesRequest) PathParams() *GetInvoicesPathParams {
	return r.pathParams
}

func (r *GetInvoicesRequest) SetMethod(method string) {
	r.method = method
}

func (r *GetInvoicesRequest) Method() string {
	return r.method
}

func (r GetInvoicesRequest) NewGetInvoicesRequestBody() GetInvoicesRequestBody {
	return GetInvoicesRequestBody{}
}

type GetInvoicesRequestBody struct{}

func (r *GetInvoicesRequest) RequestBody() *GetInvoicesRequestBody {
	return &r.requestBody
}

func (r *GetInvoicesRequest) SetRequestBody(body GetInvoicesRequestBody) {
	r.requestBody = body
}

func (r *GetInvoicesRequest) NewResponseBody() *GetInvoicesResponseBody {
	return &GetInvoicesResponseBody{}
}

type GetInvoicesResponseBody []struct {
	ID        string `json:"id"`
	Status    string `json:"status"`
	Reference string `json:"reference"`
	Number    string `json:"number"`
	Sender    struct {
		Eia     string `json:"eia"`
		Bid     string `json:"bid"`
		Name    string `json:"name"`
		Country string `json:"country"`
	} `json:"sender"`
	Recipient struct {
		EIA      string `json:"eia"`
		BID      string `json:"bid"`
		Name     string `json:"name"`
		Country  string `json:"country"`
		Operator string `json:"operator"`
	} `json:"recipient"`
	CreatedAt    time.Time     `json:"created_at"`
	Date         string        `json:"date"`
	DateDue      string        `json:"date_due"`
	SourceFormat string        `json:"source_format"`
	Sum          float64       `json:"sum"`
	SumTax       float64       `json:"sum_tax"`
	Currency     string        `json:"currency"`
	Destination  string        `json:"destination"`
	Comment      string        `json:"comment"`
	Files        []interface{} `json:"files"`
}

func (r *GetInvoicesRequest) URL() url.URL {
	return r.client.GetEndpointURL("/v1/invoices", r.PathParams())
}

func (r *GetInvoicesRequest) Do() (GetInvoicesResponseBody, error) {
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
