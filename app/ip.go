package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func getExternalIpIpify() (ExternalIPResponse, error) {
	resp, err := http.Get(IPIFY_URL)
	if err != nil {
		fmt.Println(err)
		return ExternalIPResponse{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return ExternalIPResponse{}, err
	}

	var externalIPResponse ExternalIPResponse
	err = json.Unmarshal(body, &externalIPResponse)
	if err != nil {
		fmt.Println(err)
		return ExternalIPResponse{}, err
	}
	return externalIPResponse, nil

}

func getExternalIpIpApi() (ExternalIPResponse, error) {
	resp, err := http.Get(IPAPI_URL)
	if err != nil {
		fmt.Println(err)
		return ExternalIPResponse{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return ExternalIPResponse{}, err
	}

	var ipApiPResponse IPAPIResponse
	err = json.Unmarshal(body, &ipApiPResponse)
	if err != nil {
		fmt.Println(err)
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
