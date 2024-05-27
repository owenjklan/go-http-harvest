package main

import "net/http"

// Inspired by this answer from Stack Overflow:
// https://stackoverflow.com/a/75150227

// All functions referenced in the dispatch map are expected to have
// this signature.
type extrasFunc func(response ConnectedResponse, httpResponse *http.Response) ConnectedResponse

var dispatchMap = map[string]extrasFunc{
	"authHeader": addWWWAuthHeader,
}

func addWWWAuthHeader(response ConnectedResponse, httpResponse *http.Response) ConnectedResponse {
	response.Extras["WWW-Authenticate Header"] = httpResponse.Header.Get("WWW-Authenticate")
	return response
}

func ProcessExtras(response ConnectedResponse, httpResponse *http.Response) ConnectedResponse {
	for _, extrasFuncName := range response.extrasStack {
		extrasFunc := dispatchMap[extrasFuncName]

		if extrasFunc != nil {
			response = extrasFunc(response, httpResponse)
		}
	}

	return response
}
