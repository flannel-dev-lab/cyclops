// requester package is used to perform basic http requests
package requester

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// Get will perform a GET request on the url with the headers and query params
func Get(url string, headers map[string]string, queryParams map[string][]string) (response *http.Response, err error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	for headerKey, headerValue := range headers {
		request.Header.Set(headerKey, headerValue)
	}

	if queryParams != nil {
		queryString := request.URL.Query()
		for queryKey, queryValues := range queryParams {
			for _, value := range queryValues {
				queryString.Add(queryKey, value)
			}
		}
		request.URL.RawQuery = queryString.Encode()
	}

	client := &http.Client{}
	response, err = client.Do(request)
	return response, err
}

// Post will perform a POST request on the url with the headers, query params and request body
func Post(url string, headers map[string]string, queryParams map[string][]string, requestBody interface{}) (response *http.Response, err error) {
	data, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	for headerKey, headerValue := range headers {
		request.Header.Set(headerKey, headerValue)
	}

	if queryParams != nil {
		queryString := request.URL.Query()
		for queryKey, queryValues := range queryParams {
			for _, value := range queryValues {
				queryString.Add(queryKey, value)
			}
		}
		request.URL.RawQuery = queryString.Encode()
	}

	client := &http.Client{}
	response, err = client.Do(request)
	return response, err
}

// Delete will perform a Delete request on the url with the headers, query params and request body
func Delete(url string, headers map[string]string, queryParams map[string][]string, requestBody interface{}) (response *http.Response, err error) {
	data, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("DELETE", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	for headerKey, headerValue := range headers {
		request.Header.Set(headerKey, headerValue)
	}

	if queryParams != nil {
		queryString := request.URL.Query()
		for queryKey, queryValues := range queryParams {
			for _, value := range queryValues {
				queryString.Add(queryKey, value)
			}
		}
		request.URL.RawQuery = queryString.Encode()
	}

	client := &http.Client{}
	response, err = client.Do(request)
	return response, err
}

// Put will perform a PUT request on the url with the headers, query params and request body
func Put(url string, headers map[string]string, queryParams map[string][]string, requestBody interface{}) (response *http.Response, err error) {
	data, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("PUT", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	for headerKey, headerValue := range headers {
		request.Header.Set(headerKey, headerValue)
	}

	if queryParams != nil {
		queryString := request.URL.Query()
		for queryKey, queryValues := range queryParams {
			for _, value := range queryValues {
				queryString.Add(queryKey, value)
			}
		}
		request.URL.RawQuery = queryString.Encode()
	}

	client := &http.Client{}
	response, err = client.Do(request)
	return response, err
}

// Patch will perform a Patch request on the url with the headers, query params and request body
func Patch(url string, headers map[string]string, queryParams map[string][]string, requestBody interface{}) (response *http.Response, err error) {
	data, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("PATCH", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	for headerKey, headerValue := range headers {
		request.Header.Set(headerKey, headerValue)
	}

	if queryParams != nil {
		queryString := request.URL.Query()
		for queryKey, queryValues := range queryParams {
			for _, value := range queryValues {
				queryString.Add(queryKey, value)
			}
		}
		request.URL.RawQuery = queryString.Encode()
	}

	client := &http.Client{}
	response, err = client.Do(request)
	return response, err
}
