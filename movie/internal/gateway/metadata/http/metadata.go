package http

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mamalmaleki/go_movie/metadata/pkg/model"
	"github.com/mamalmaleki/go_movie/movie/internal/gateway"
	"github.com/mamalmaleki/go_movie/pkg/discovery"
	"log"
	"math/rand"
	"net/http"
)

// Gateway defines a movie metadata HTTP gateway
type Gateway struct {
	registry discovery.Registry
	//addr string
}

// New creates a new HTTP gateway for a movie metadata service.
func New(registry discovery.Registry) *Gateway {
	return &Gateway{registry}
}

//func New(addr string) *Gateway {
//	return &Gateway{addr}
//}

// Get gets movie metadata by a movie id.
func (g *Gateway) Get(ctx context.Context, id string) (*model.Metadata, error) {
	address, err := g.registry.ServiceAddresses(ctx, "metadata")
	if err != nil {
		return nil, err
	}
	url := "http://" + address[rand.Intn(len(address))] + "/metadata"
	log.Println("Calling metadata service. Request: GET " + url)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	//req, err := http.NewRequest(http.MethodGet, g.addr+"/metadata", nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	values := req.URL.Query()
	values.Add("id", id)
	req.URL.RawQuery = values.Encode()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, gateway.ErrNotFound
	} else if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("non-2xx response: %v", resp)
	}
	var v *model.Metadata
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}
	return v, nil
}
