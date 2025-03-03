package main

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/avast/retry-go/v4"
)

var MAX_RETRY_ATTEMPTS uint = 5

func getResponseBody(url string) ([]byte, error) {
	slog.Debug("Getting response body", "url", url)
	return retry.DoWithData(
		func() ([]byte, error) {
			resp, err := http.Get(url)
			if err != nil {
				return nil, err
			}
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}

			return body, nil
		}, retry.Attempts(MAX_RETRY_ATTEMPTS), retry.OnRetry(func(attempt uint, err error) {
			slog.Debug("Retrying", "url", url, "attempt", attempt, "maxAttempt", MAX_RETRY_ATTEMPTS)
		}),
	)
}
func getExternalIpIpify() (ExternalIPResponse, error) {
	body, err := getResponseBody(IPIFY_URL)

	if err != nil {
		slog.Error(err.Error())
		return ExternalIPResponse{}, err
	}

	var externalIPResponse ExternalIPResponse
	err = json.Unmarshal(body, &externalIPResponse)
	if err != nil {
		slog.Error(err.Error())
		return ExternalIPResponse{}, err
	}
	return externalIPResponse, nil

}

func getExternalIpIpApi() (ExternalIPResponse, error) {
	body, err := getResponseBody(IPAPI_URL)

	if err != nil {
		slog.Error(err.Error())
		return ExternalIPResponse{}, err
	}

	var ipApiPResponse IPAPIResponse
	err = json.Unmarshal(body, &ipApiPResponse)
	if err != nil {
		slog.Error(err.Error())
		return ExternalIPResponse{}, err
	}
	return ExternalIPResponse{IP: ipApiPResponse.Query}, nil

}

func getExternalIp() (ExternalIPResponse, error) {
	ipifyResponse, err := getExternalIpIpify()
	if err != nil {
		ipapiResponse, err := getExternalIpIpApi()
		if err != nil {
			return ExternalIPResponse{}, err
		}
		return ipapiResponse, nil
	}
	return ipifyResponse, nil
}
