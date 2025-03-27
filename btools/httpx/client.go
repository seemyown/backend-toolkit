package httpx

import (
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/seemyown/backend-toolkit/btools/logging"
	"time"
)

type RestClient struct {
	client  *resty.Client
	rootURL string
	headers map[string]string
	logger  logging.Logger
}

func NewRestClient(rootURL string, headers map[string]string, logger *logging.Logger) *RestClient {
	client := resty.New().SetBaseURL(rootURL).SetTimeout(60 * time.Second)
	if headers != nil {
		client.SetHeaders(headers)
	}
	if logger == nil {
		logger = logging.New(logging.Config{FileName: "rest_cleint", Name: "client"})
	}
	return &RestClient{
		client:  client,
		rootURL: rootURL,
		headers: headers,
		logger:  *logger,
	}
}

func (rc *RestClient) requestBuilder(
	ctx context.Context,
	queryParams map[string]string,
	requestBody interface{},
	timeout time.Duration,
) *resty.Request {
	if timeout > 0 {
		rc.client.SetTimeout(timeout)
	}
	req := rc.client.R().SetContext(ctx)
	if queryParams != nil {
		req.SetQueryParams(queryParams)
	}
	if requestBody != nil {
		req.SetBody(requestBody)
	}
	return req
}

func (rc *RestClient) MakeGetRequest(
	ctx context.Context,
	endpoint string,
	queryParams map[string]string,
	timeout time.Duration,
	maxAttempts int,
) (*resty.Response, error) {
	req := rc.requestBuilder(ctx, queryParams, nil, timeout)
	rc.logger.Debug("GET request to %s, params %+v, timeout %d", endpoint, queryParams, timeout)
	for attempt := 0; attempt < maxAttempts; attempt++ {
		rc.logger.Debug("GET request to %s attempt %d / %d", endpoint, attempt, maxAttempts)
		resp, err := req.Get(endpoint)
		if err != nil {
			rc.logger.Error(err, "failed to get response from endpoint %s", endpoint)
			continue
		}
		rc.logger.Debug("response from endpoint %s: %s [%d]", endpoint, resp.String(), len(resp.String()))
		return resp, nil
	}
	return nil, fmt.Errorf("failed to get response from endpoint %s", endpoint)
}

func (rc *RestClient) MakePostRequest(
	ctx context.Context,
	endpoint string,
	queryParams map[string]string,
	requestBody interface{},
	timeout time.Duration,
) (*resty.Response, error) {
	req := rc.requestBuilder(ctx, queryParams, requestBody, timeout)
	return req.Post(endpoint)
}

func (rc *RestClient) MakePutRequest(
	ctx context.Context,
	endpoint string,
	queryParams map[string]string,
	requestBody interface{},
	timeout time.Duration,
) (*resty.Response, error) {
	req := rc.requestBuilder(ctx, queryParams, requestBody, timeout)
	return req.Put(endpoint)
}

func (rc *RestClient) MakePatchRequest(
	ctx context.Context,
	endpoint string,
	queryParams map[string]string,
	requestBody interface{},
	timeout time.Duration,
) (*resty.Response, error) {
	req := rc.requestBuilder(ctx, queryParams, requestBody, timeout)
	return req.Patch(endpoint)
}

func (rc *RestClient) MakeDeleteRequest(
	ctx context.Context,
	endpoint string,
	queryParams map[string]string,
	timeout time.Duration,
) (*resty.Response, error) {
	req := rc.requestBuilder(ctx, queryParams, nil, timeout)
	return req.Delete(endpoint)
}
