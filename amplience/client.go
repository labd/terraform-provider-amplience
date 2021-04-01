package amplience

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Placeholder client configuration struct to pass client info into meta
type ClientConfig struct {
	ID            string
	Secret        string
	HubID         string
	ContentApiUrl string
	AuthUrl       string
}

// AmplienceRequest is a util func to abstract HTTP requests to the Amplience API which will be placeholders for a poc
// before we develop an Amplience Client library to handle the requests
func AmplienceRequest(config *ClientConfig, url string, requestType string, requestBody *bytes.Buffer) (*http.Response, error) {
	var req *http.Request
	var err error
	switch requestType {
	case http.MethodDelete:
		req, err = http.NewRequest(requestType, url, nil)
	case http.MethodGet:
		req, err = http.NewRequest(requestType, url, nil)
	case http.MethodPost:
		req, err = http.NewRequest(requestType, url, requestBody)
	case http.MethodPatch:
		req, err = http.NewRequest(requestType, url, requestBody)
	default:
		return nil, fmt.Errorf("unsupported Amplience RequestType %s", requestType)
	}
	if err != nil {
		return nil, fmt.Errorf("error creating %s request to %s for body %v: %w", requestType, url, requestBody, err)
	}

	token, err := getAmplienceOAuthToken(config)
	if err != nil {
		return nil, fmt.Errorf("could not get Oauth token for request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("content-type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error during %s request: %w", requestType, err)
	}

	return resp, nil
}

func ParseAndUnmarshalAmplienceResponseBody(response *http.Response, data interface{}) error {
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("could not read response body %v: %w", response.Body, err)
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return fmt.Errorf("could not unmarshal %v: %w", response.Body, err)
	}
	if data == nil {
		return fmt.Errorf("struct %v is nil after unmarshalling", data)
	}
	return nil
}

// https://amplience.com/docs/api/dynamic-content/management/index.html#section/Usage/Status-Code-Table
func HandleAmplienceError(response *http.Response) *resource.RetryError {
	if response == nil {
		return resource.NonRetryableError(fmt.Errorf("received nil response, unable to handle error"))
	}
	errResponseBody := ErrorResponse{}
	switch response.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusAccepted, http.StatusNoContent:
		return nil
	case http.StatusInternalServerError:
		data, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("could not read error response body: %w", err))
		}
		err = json.Unmarshal(data, &errResponseBody)
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("could not unmarshal error response body: %w", err))
		}
		return resource.RetryableError(fmt.Errorf("retryable error with code %d received: %s\n "+
			"Amplience Error Response: %s", response.StatusCode, response.Status, StringFormatObject(errResponseBody)))
	default:
		data, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("could not read error response body: %w", err))
		}
		err = json.Unmarshal(data, &errResponseBody)
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("could not unmarshal error response body: %w", err))
		}
		return resource.NonRetryableError(fmt.Errorf("non retryable error with code %d received: %s\n "+
			"Amplience Error Response: %s", response.StatusCode, response.Status, StringFormatObject(errResponseBody)))
	}
}

func getAmplienceOAuthToken(config *ClientConfig) (string, error) {
	data := url.Values{}
	data.Set("client_id", config.ID)
	data.Set("client_secret", config.Secret)
	data.Set("grant_type", "client_credentials")

	req, err := http.NewRequest(http.MethodPost, config.AuthUrl, strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("could not create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error during request %v: %w", req, err)
	}

	tokenStruct := authResponseBody{}
	if resp.StatusCode == http.StatusOK {
		defer resp.Body.Close()
		err = json.NewDecoder(resp.Body).Decode(&tokenStruct)
		if err != nil {
			return "", fmt.Errorf("could not decode %v into %v: %w", resp.Body, tokenStruct, err)
		}
	} else {
		return "", fmt.Errorf("received statuscode %d", resp.StatusCode)
	}
	if tokenStruct.AccessToken == "" {
		return "", fmt.Errorf("did not receive Oauth token")
	}

	return tokenStruct.AccessToken, nil
}

type authResponseBody struct {
	AccessToken      string `json:"access_token"`
	SessionExpiresIn int    `json:"session_expires_in"`
	ExpiresIn        int    `json:"expires_in"`
}

type ErrorResponseItem struct {
	Level        string `json:"level"`
	Code         string `json:"code"`
	Message      string `json:"message"`
	Property     string `json:"property"`
	Entity       string `json:"entity"`
	InvalidValue string `json:"invalidValue"`
}

type ErrorResponse struct {
	Errors []ErrorResponseItem `json:"errors"`
}
