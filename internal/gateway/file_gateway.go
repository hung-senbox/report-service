package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"report-service/helper"
	"report-service/internal/gateway/dto/request"
	dto "report-service/internal/gateway/dto/response"
	"report-service/pkg/constants"

	"github.com/hashicorp/consul/api"
)

type FileGateway interface {
	GetImageUrl(ctx context.Context, req request.GetFileUrlRequest) (*string, error)
}

type fileGateway struct {
	serviceName string
	consul      *api.Client
}

func NewFileGateway(serviceName string, consul *api.Client) FileGateway {
	return &fileGateway{
		serviceName: serviceName,
		consul:      consul,
	}
}

func (g *fileGateway) GetImageUrl(ctx context.Context, req request.GetFileUrlRequest) (*string, error) {
	token, ok := ctx.Value(constants.Token).(string)
	if !ok {
		return nil, fmt.Errorf("token not found in context")
	}

	client, err := NewGatewayClient(g.serviceName, token, g.consul, nil)
	if err != nil {
		return nil, err
	}

	headers := helper.GetHeaders(ctx)
	resp, err := client.Call("POST", "/v1/gateway/images/get-url", req, headers)
	if err != nil {
		return nil, err
	}

	var gwResp dto.APIGateWayResponse[string]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return nil, fmt.Errorf("unmarshal response fail: %w", err)
	}

	if gwResp.StatusCode != 200 {
		return nil, fmt.Errorf("call gateway get image fail: %s", gwResp.Message)
	}

	return &gwResp.Data, nil
}
