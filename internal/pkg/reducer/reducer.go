package reducer

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/SakiiR/ReduceRequest/internal/pkg/parser"
)

func reduceURIParameters(request *http.Request, parser *parser.Parser) http.Request {
	params := request.URL.Query()
	for key, values := range params {
		// Remove the current query key
		params.Del(key)
		// Construct the new parameters
		request.URL.RawQuery = params.Encode()

		status, _ := validateResponse(parser.InitialResponse, request, parser)
		if status == true {
			fmt.Println(fmt.Sprintf("Ok, parameter %s hasn't impact on response", key))
		} else {
			fmt.Println(fmt.Sprintf("Ok, parameter %s has impact on response", key))
			for _, value := range values {
				params.Add(key, value)
			}
			request.URL.RawQuery = params.Encode()
		}
	}

	return *request
}

func reduceHeaders(request *http.Request, parser *parser.Parser) http.Request {

	headers := request.Header
	for key, values := range headers {

		request.Header.Del(key)

		status, _ := validateResponse(parser.InitialResponse, request, parser)
		if status == true {
			fmt.Println(fmt.Sprintf("Ok, header %s hasn't impact on response", key))
		} else {

			fmt.Println(fmt.Sprintf("Ok, header %s has impact on response", key))
			for _, value := range values {
				request.Header.Add(key, value)
			}
		}

	}

	return *request
}

func serializeCookies(cookies []*http.Cookie) string {

	cookiesStr := make([]string, len(cookies))
	for _, cookie := range cookies {
		if cookie.Value != "" {
			cookiesStr = append(cookiesStr, fmt.Sprintf("%s=%s", cookie.Name, cookie.Value))
		}
	}

	str := strings.Join(cookiesStr, ";")
	str = strings.Trim(str, ";")

	return str
}

func reduceCookies(request *http.Request, parser *parser.Parser) http.Request {
	cookies := request.Cookies()

	for _, cookie := range cookies {
		valueSave := cookie.Value
		cookie.Value = ""

		request.Header.Set("Cookie", serializeCookies(cookies))

		status, _ := validateResponse(parser.InitialResponse, request, parser)
		if status == true {
			fmt.Println(fmt.Sprintf("Ok, cookie %s hasn't impact on response", cookie.Name))
		} else {
			fmt.Println(fmt.Sprintf("Ok, cookie %s has an impact on response", cookie.Name))
			cookie.Value = valueSave
			request.Header.Set("Cookie", serializeCookies(cookies))
		}
	}

	return *request
}

// ReduceRequest reduces request
func ReduceRequest(parser *parser.Parser) (*http.Request, error) {
	r := *parser.Request
	r = reduceURIParameters(&r, parser)
	r = reduceCookies(&r, parser)
	r = reduceHeaders(&r, parser)
	// TODO: iterate over form parameters if form
	// TODO: iterate over json parameters if json
	// TODO: iterate over xml parameters if xml
	fmt.Println()
	fmt.Println("```")
	DumpRequestToStdout(&r)
	fmt.Println("```")
	return nil, nil
}

// DumpRequestToStdout dumps the specified request to stdout
func DumpRequestToStdout(request *http.Request) error {

	data, err := httputil.DumpRequest(request, true)
	if err != nil {
		fmt.Println(fmt.Sprintf("Failed to dump request: %s", err))
		return err
	}

	fmt.Println(string(data))
	return nil
}

// validateResponse checks the request lengths to identify valid/invalid request
func validateResponse(initialResponse *http.Response, request *http.Request, parser *parser.Parser) (bool, error) {

	response, err := parser.Do(request)
	if err != nil {
		return false, err
	}

	data1, err := httputil.DumpResponse(initialResponse, true)
	if err != nil {
		fmt.Println(fmt.Sprintf("Failed to dump response 1: %s", err))
		return false, err
	}

	data2, err := httputil.DumpResponse(response, true)
	if err != nil {
		fmt.Println(fmt.Sprintf("Failed to dump response 2: %s", err))
		return false, err
	}

	if len(data1) != len(data2) {
		return false, nil
	}

	return true, nil
}
