// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later
//

package madmin

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	jwtgo "github.com/golang-jwt/jwt/v4"
	"github.com/openstor/openstor-go/v7/pkg/credentials"
)

const (
	defaultPrometheusJWTExpiry = 100 * 365 * 24 * time.Hour
	libraryMinioURLPrefix      = "/minio"
	prometheusIssuer           = "prometheus"
)

// MetricsClient implements MinIO metrics operations
type MetricsClient struct {
	/// Credentials for authentication
	creds *credentials.Credentials
	// Indicate whether we are using https or not
	secure bool
	// Parsed endpoint url provided by the user.
	endpointURL *url.URL
	// Needs allocation.
	httpClient *http.Client
}

// metricsRequestData - is container for all the values to make a
// request.
type metricsRequestData struct {
	relativePath string // URL path relative to admin API base endpoint
}

// NewMetricsClientWithOptions - instantiate minio metrics client honoring Prometheus format
func NewMetricsClientWithOptions(endpoint string, opts *Options) (*MetricsClient, error) {
	if opts == nil {
		return nil, ErrInvalidArgument("empty options not allowed")
	}

	endpointURL, err := getEndpointURL(endpoint, opts.Secure)
	if err != nil {
		return nil, err
	}

	clnt, err := privateNewMetricsClient(endpointURL, opts)
	if err != nil {
		return nil, err
	}
	return clnt, nil
}

// getPrometheusToken creates a JWT from MinIO access and secret keys
func getPrometheusToken(accessKey, secretKey string) (string, error) {
	jwt := jwtgo.NewWithClaims(jwtgo.SigningMethodHS512, jwtgo.RegisteredClaims{
		ExpiresAt: jwtgo.NewNumericDate(time.Now().UTC().Add(defaultPrometheusJWTExpiry)),
		Subject:   accessKey,
		Issuer:    prometheusIssuer,
	})

	token, err := jwt.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return token, nil
}

func privateNewMetricsClient(endpointURL *url.URL, opts *Options) (*MetricsClient, error) {
	clnt := new(MetricsClient)
	clnt.creds = opts.Creds
	clnt.secure = opts.Secure
	clnt.endpointURL = endpointURL

	tr := opts.Transport
	if tr == nil {
		tr = DefaultTransport(opts.Secure)
	}

	clnt.httpClient = &http.Client{
		Transport: tr,
	}
	return clnt, nil
}

// executeGetRequest - instantiates a Get method and performs the request
func (client *MetricsClient) executeGetRequest(ctx context.Context, reqData metricsRequestData) (res *http.Response, err error) {
	req, err := client.newGetRequest(ctx, reqData)
	if err != nil {
		return nil, err
	}

	v, err := client.creds.GetWithContext(client.CredContext())
	if err != nil {
		return nil, err
	}

	accessKeyID := v.AccessKeyID
	secretAccessKey := v.SecretAccessKey

	jwtToken, err := getPrometheusToken(accessKeyID, secretAccessKey)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+jwtToken)
	req.Header.Set("X-Amz-Security-Token", v.SessionToken)

	return client.httpClient.Do(req)
}

// newGetRequest - instantiate a new HTTP GET request
func (client *MetricsClient) newGetRequest(ctx context.Context, reqData metricsRequestData) (req *http.Request, err error) {
	targetURL, err := client.makeTargetURL(reqData)
	if err != nil {
		return nil, err
	}

	return http.NewRequestWithContext(ctx, http.MethodGet, targetURL.String(), nil)
}

// makeTargetURL make a new target url.
func (client *MetricsClient) makeTargetURL(r metricsRequestData) (*url.URL, error) {
	if client.endpointURL == nil {
		return nil, fmt.Errorf("enpointURL cannot be nil")
	}

	host := client.endpointURL.Host
	scheme := client.endpointURL.Scheme
	prefix := libraryMinioURLPrefix

	urlStr := scheme + "://" + host + prefix + r.relativePath
	return url.Parse(urlStr)
}

// CredContext returns the context for fetching credentials
func (client *MetricsClient) CredContext() *credentials.CredContext {
	return &credentials.CredContext{
		Client: client.httpClient,
	}
}
