package client

import (
	"bytes"
	"effective_mobile_test/internal/storage/postgres"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"time"
)

var (
	ErrBadReq      = errors.New("bad request")
	ErrTimeout     = errors.New("timeout reached")
	ErrBadConn     = errors.New("can't connect to server")
	ErrServerFatal = errors.New("server fatal error")
	ErrBadResp     = errors.New("bad response")
)

type SearchRequest struct {
	RegNums []string
}

type SearchResponse struct {
	Cars []postgres.Car
}

type SearchClient struct {
	URL string
}

// FindUsers отправляет запрос во внешнюю систему, которая непосредственно ищет пользователей
func (srv *SearchClient) FindUsers(req SearchRequest) (*SearchResponse, error) {
	rawBody, err := json.Marshal(req)
	if err != nil {
		return nil, ErrBadReq
	}

	body := bytes.NewReader(rawBody)

	searcherReq, _ := http.NewRequest("GET", srv.URL+"/info", body)

	client := &http.Client{Timeout: time.Second}

	resp, err := client.Do(searcherReq)
	if err != nil {
		var err net.Error
		if errors.As(err, &err) && err.Timeout() {
			return nil, ErrTimeout
		}
		return nil, ErrBadConn
	}
	defer resp.Body.Close()
	receivedBody, _ := io.ReadAll(resp.Body) //nolint:errcheck

	switch resp.StatusCode {
	case http.StatusInternalServerError:
		return nil, ErrServerFatal
	case http.StatusBadRequest:
		return nil, ErrBadReq
	}

	result := SearchResponse{}
	err = json.Unmarshal(receivedBody, &result)
	if err != nil {
		return nil, ErrBadResp
	}

	return &result, err
}
