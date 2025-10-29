// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later
//

package madmin

import (
	"fmt"
	"net/url"
	"testing"

	jwtgo "github.com/golang-jwt/jwt/v4"
)

func TestMakeTargetUrlBuildsURLWithClientAndRelativePath(t *testing.T) {
	clnt := MetricsClient{
		endpointURL: &url.URL{
			Host:   "localhost:9000",
			Scheme: "http",
		},
	}
	requestData := metricsRequestData{
		relativePath: "/some/path",
	}

	targetURL, err := clnt.makeTargetURL(requestData)
	if err != nil {
		t.Errorf("error not expected, got: %v", err)
	}

	expectedURL := "http://localhost:9000/minio/some/path"
	if expectedURL != targetURL.String() {
		t.Errorf("target url: %s  not equal to expected url: %s", targetURL, expectedURL)
	}
}

func TestMakeTargetUrlReturnsErrorIfEndpointURLNotSet(t *testing.T) {
	clnt := MetricsClient{}
	requestData := metricsRequestData{
		relativePath: "/some/path",
	}

	_, err := clnt.makeTargetURL(requestData)
	if err == nil {
		t.Errorf("error expected got nil")
	}
}

func TestMakeTargetUrlReturnsErrorOnURLParse(t *testing.T) {
	clnt := MetricsClient{
		endpointURL: &url.URL{},
	}
	requestData := metricsRequestData{
		relativePath: "/some/path",
	}

	_, err := clnt.makeTargetURL(requestData)
	if err == nil {
		t.Errorf("error expected got nil")
	}
}

func TestGetPrometheusTokenReturnsValidJwtTokenFromAccessAndSecretKey(t *testing.T) {
	accessKey := "myaccessKey"
	secretKey := "mysecretKey"

	jwtToken, err := getPrometheusToken(accessKey, secretKey)
	if err != nil {
		t.Errorf("error not expected, got: %v", err)
	}

	token, err := jwtgo.Parse(jwtToken, func(token *jwtgo.Token) (interface{}, error) {
		// Set same signing method used in our function
		if _, ok := token.Method.(*jwtgo.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		t.Errorf("error not expected, got: %v", err)
	}
	if !token.Valid {
		t.Errorf("invalid token: %s", jwtToken)
	}
}
