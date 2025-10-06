package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	gw_response "report-service/internal/gateway/dto/response"
	"report-service/pkg/constants"

	"github.com/hashicorp/consul/api"
)

type TermGateway interface {
	GetTermByID(ctx context.Context, termID string) (*gw_response.TermResponse, error)
}

type termGateway struct {
	serviceName string
	consul      *api.Client
}

func NewTermGateway(serviceName string, consulClient *api.Client) TermGateway {
	return &termGateway{
		serviceName: serviceName,
		consul:      consulClient,
	}
}

func (g *termGateway) GetTermByID(ctx context.Context, termID string) (*gw_response.TermResponse, error) {
	// lấy token từ context
	token, ok := ctx.Value(constants.Token).(string)
	if !ok {
		return nil, fmt.Errorf("token not found in context")
	}

	// tạo client
	client, err := NewGatewayClient(g.serviceName, token, g.consul, nil)
	if err != nil {
		return nil, err
	}

	// gọi API với query params
	url := fmt.Sprintf("/api/v1/gateway/terms/%s", termID)
	resp, err := client.Call("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// parse JSON
	var gwResp gw_response.APIGateWayResponse[*gw_response.TermResponse]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return nil, fmt.Errorf("unmarshal response fail: %w", err)
	}

	// check status
	if gwResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("call gateway get term fail: %s", gwResp.Message)
	}

	return gwResp.Data, nil
}
