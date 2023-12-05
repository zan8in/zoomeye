package retryhttp

import (
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"runtime"
	"time"

	randutil "github.com/zan8in/pins/rand"
	"github.com/zan8in/retryablehttp"
)

var (
	RtryRedirect   *retryablehttp.Client
	RtryNoRedirect *retryablehttp.Client

	RtryNoRedirectHttpClient *http.Client
	RtryRedirectHttpClient   *http.Client
	defaultMaxRedirects      = 10
)

const maxDefaultBody = 2 * 1024 * 1024

type Options struct {
	Timeout int
	Retries int

	Proxy string
}

type ClientResponse struct {
	Body       string
	Header     string
	StatusCode int
	BodySize   int
}

func Init(options *Options) (err error) {
	retryableHttpOptions := retryablehttp.DefaultOptionsSpraying
	maxIdleConns := 0
	maxConnsPerHost := 0
	maxIdleConnsPerHost := -1
	disableKeepAlives := true // 默认 false

	// retryableHttpOptions = retryablehttp.DefaultOptionsSingle
	// disableKeepAlives = false
	// maxIdleConnsPerHost = 500
	// maxConnsPerHost = 500

	maxIdleConns = 1000                        //
	maxIdleConnsPerHost = runtime.NumCPU() * 2 //
	idleConnTimeout := 15 * time.Second        //
	tLSHandshakeTimeout := 5 * time.Second     //

	dialer := &net.Dialer{ //
		Timeout:   time.Duration(options.Timeout) * time.Second,
		KeepAlive: 15 * time.Second,
	}

	retryableHttpOptions.RetryWaitMax = 10 * time.Second
	retryableHttpOptions.RetryMax = options.Retries

	tlsConfig := &tls.Config{
		Renegotiation:      tls.RenegotiateOnceAsClient,
		InsecureSkipVerify: true,
		MinVersion:         tls.VersionTLS10,
	}

	transport := &http.Transport{
		DialContext:         dialer.DialContext,
		MaxIdleConns:        maxIdleConns,
		MaxIdleConnsPerHost: maxIdleConnsPerHost,
		MaxConnsPerHost:     maxConnsPerHost,
		TLSClientConfig:     tlsConfig,
		DisableKeepAlives:   disableKeepAlives,
		TLSHandshakeTimeout: tLSHandshakeTimeout, //
		IdleConnTimeout:     idleConnTimeout,     //
	}

	// follow redirects client
	// clientCookieJar, _ := cookiejar.New(nil)

	httpRedirectClient := http.Client{
		Transport: transport,
		Timeout:   time.Duration(options.Timeout) * time.Second,
		// Jar:       clientCookieJar,
	}

	RtryRedirect = retryablehttp.NewWithHTTPClient(&httpRedirectClient, retryableHttpOptions)
	RtryRedirect.CheckRetry = retryablehttp.HostSprayRetryPolicy()
	RtryRedirectHttpClient = RtryRedirect.HTTPClient

	// whitespace

	// disabled follow redirects client
	// clientNoRedirectCookieJar, _ := cookiejar.New(nil)

	httpNoRedirectClient := http.Client{
		Transport: transport,
		Timeout:   time.Duration(options.Timeout) * time.Second,
		// Jar:           clientNoRedirectCookieJar,
		CheckRedirect: makeCheckRedirectFunc(false, defaultMaxRedirects),
	}

	RtryNoRedirect = retryablehttp.NewWithHTTPClient(&httpNoRedirectClient, retryableHttpOptions)
	RtryNoRedirect.CheckRetry = retryablehttp.HostSprayRetryPolicy()
	RtryNoRedirectHttpClient = RtryNoRedirect.HTTPClient

	return err
}

type checkRedirectFunc func(req *http.Request, via []*http.Request) error

func makeCheckRedirectFunc(followRedirects bool, maxRedirects int) checkRedirectFunc {
	return func(req *http.Request, via []*http.Request) error {
		if !followRedirects {
			return http.ErrUseLastResponse
		}

		if maxRedirects == 0 {
			if len(via) > defaultMaxRedirects {
				return http.ErrUseLastResponse
			}
			return nil
		}

		if len(via) > maxRedirects {
			return http.ErrUseLastResponse
		}
		return nil
	}
}

func GetWithApiKey(target, apiKey string) ([]byte, int, error) {
	var (
		body   []byte
		status int
		err    error
	)
	req, err := retryablehttp.NewRequest(http.MethodGet, target, nil)
	if err != nil {
		return body, status, err
	}
	req.Header.Set("User-Agent", randutil.RandomUA())
	req.Header.Set("API-KEY", apiKey)

	resp, err := RtryRedirect.Do(req)
	if err != nil {
		if resp != nil {
			resp.Body.Close()
		}
		return body, status, err
	}
	defer resp.Body.Close()

	status = resp.StatusCode

	if resp.StatusCode == http.StatusSwitchingProtocols {
		return body, status, err
	}

	reader := io.LimitReader(resp.Body, maxDefaultBody)
	respBody, err := io.ReadAll(reader)
	if err != nil {
		return body, status, err
	}

	return respBody, status, nil
}
