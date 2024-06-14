package mt

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/lokks307/djson/v2"
)

const NET_TIMEOUT = 10
const CRM_NET_TIMEOUT = 300

func MakePath(base string, dd ...string) string {
	format := strings.Repeat("/%s", len(dd)+1)
	args := make([]interface{}, 0)
	args = append(args, url.QueryEscape(base))
	for idx := range dd {
		args = append(args, url.QueryEscape(dd[idx]))
	}
	return fmt.Sprintf(format, args...)
}

func MakePathForCareEase(base string, dd ...string) string {
	format := strings.Repeat("/%s", len(dd)+1)
	args := make([]interface{}, 0)
	args = append(args, base)
	for idx := range dd {
		args = append(args, dd[idx])
	}
	return fmt.Sprintf(format, args...)
}

func MakeUrl(base string, subpath ...string) string {
	u, _ := url.Parse(base)

	for i := range subpath {
		u.Path = path.Join(u.Path, subpath[i])
	}

	return u.String()
}

type httpResponse struct {
	StatusCode  int
	ContentMore bool
	Body        []byte
}

func QueryJsonToUrlRaw(method, baseUrl, subPath, reqBody string, query *djson.JSON, header ...*djson.JSON) (*httpResponse, error) {
	reqBodyReader := strings.NewReader(reqBody)

	ctx, cancel := context.WithTimeout(context.Background(), NET_TIMEOUT*time.Second)
	defer cancel()

	path := MakeUrl(baseUrl, subPath)

	req, err := http.NewRequestWithContext(ctx, method, path, reqBodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Connection", "close")

	if len(header) == 1 && header[0] != nil {
		keys := header[0].GetKeys()
		for _, ek := range keys {
			req.Header.Set(ek, header[0].String(ek))
		}
	}

	if query != nil {
		queries := url.Values{}
		for _, k := range query.GetKeys() {
			queries.Add(k, query.String(k))
		}

		req.URL.RawQuery = queries.Encode()
	}

	client := &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	bodyText, err2 := io.ReadAll(resp.Body)
	if err2 != nil {
		return nil, err2
	}

	return &httpResponse{
		StatusCode:  resp.StatusCode,
		ContentMore: strings.EqualFold(resp.Header.Get("Content-More"), "true"),
		Body:        bodyText,
	}, nil
}

func QueryJsonToCrm(method, baseUrl, subPath, reqBody string, query *djson.JSON, header ...*djson.JSON) (*djson.JSON, error) {
	reqBodyReader := strings.NewReader(reqBody)

	ctx, cancel := context.WithTimeout(context.Background(), CRM_NET_TIMEOUT*time.Second)
	defer cancel()

	path := MakeUrl(baseUrl, subPath)

	req, err := http.NewRequestWithContext(ctx, method, path, reqBodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Connection", "close")

	if len(header) == 1 && header[0] != nil {
		keys := header[0].GetKeys()
		for _, ek := range keys {
			req.Header.Set(ek, header[0].String(ek))
		}
	}

	if query != nil {
		queries := url.Values{}
		for _, k := range query.GetKeys() {
			queries.Add(k, query.String(k))
		}

		req.URL.RawQuery = queries.Encode()
	}

	client := &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if err != nil || (resp != nil && (resp.StatusCode >= 300 || resp.StatusCode < 200)) {
		statusCode := -1
		if resp != nil {
			statusCode = resp.StatusCode
		}
		return nil, errors.New("net statusCode:" + strconv.Itoa(statusCode))
	}

	respBody, err2 := io.ReadAll(resp.Body)
	if err2 != nil {
		return nil, err2
	}

	return djson.New().Parse(string(respBody)), nil
}

func QueryJsonToUrl(method, baseUrl, subPath, reqBody string, query *djson.JSON, header ...*djson.JSON) (*djson.JSON, error) {
	reqBodyReader := strings.NewReader(reqBody)

	ctx, cancel := context.WithTimeout(context.Background(), NET_TIMEOUT*time.Second)
	defer cancel()

	path := MakeUrl(baseUrl, subPath)

	req, err := http.NewRequestWithContext(ctx, method, path, reqBodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Connection", "close")

	if len(header) == 1 && header[0] != nil {
		keys := header[0].GetKeys()
		for _, ek := range keys {
			req.Header.Set(ek, header[0].String(ek))
		}
	}

	if query != nil {
		queries := url.Values{}
		for _, k := range query.GetKeys() {
			queries.Add(k, query.String(k))
		}

		req.URL.RawQuery = queries.Encode()
	}

	client := &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if err != nil || (resp != nil && (resp.StatusCode >= 300 || resp.StatusCode < 200)) {
		statusCode := -1
		if resp != nil {
			statusCode = resp.StatusCode
		}
		return nil, errors.New("net statusCode:" + strconv.Itoa(statusCode))
	}

	respBody, err2 := io.ReadAll(resp.Body)
	if err2 != nil {
		return nil, err2
	}

	return djson.New().Parse(string(respBody)), nil
}
