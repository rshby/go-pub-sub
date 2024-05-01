package httpclient

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type LoggingRoundTripper struct {
	Next http.RoundTripper
}

func NewLogginRoundTripper() http.RoundTripper {
	return &LoggingRoundTripper{
		Next: http.DefaultTransport,
	}
}

// RoundTrip implements rountripper interface
func (l *LoggingRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	logrus.Info("masuk sini request httpclient loggin roundtripper")
	req := request.Body
	fmt.Printf("req : %v\n", req)
	//requestBody, _ := io.ReadAll(req)
	//logrus.Infof("request body : ", string(requestBody))

	//request.Body = io.NopCloser(bytes.NewReader(requestBody))
	//_ = request.Body.Close()
	res, err := l.Next.RoundTrip(request)
	if err != nil {
		logrus.Error(err)
		return res, err
	}

	// get response body
	responseBodyByte, err := io.ReadAll(res.Body)
	_ = res.Body.Close()
	logrus.Info(string(responseBodyByte))

	// reassign
	res.Body = io.NopCloser(bytes.NewReader(responseBodyByte))
	_ = res.Body.Close()

	return res, nil
}
