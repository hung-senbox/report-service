package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"report-service/helper"
	gw_response "report-service/internal/gateway/dto/response"
	"report-service/pkg/constants"

	"github.com/hashicorp/consul/api"
)

type MediaGateway interface {
	GetTopicByID(ctx context.Context, topicID string) (*gw_response.TopicResponse, error)
	GetAllTopicsByOrganization(ctx context.Context, organizationID string) ([]gw_response.TopicResponse, error)
	GetTopicByStudentID(ctx context.Context, studentID string) ([]gw_response.TopicResponse, error)
}

type mediaGateway struct {
	serviceName string
	consul      *api.Client
}

func NewMediaGateway(serviceName string, consulClient *api.Client) MediaGateway {
	return &mediaGateway{
		serviceName: serviceName,
		consul:      consulClient,
	}
}

func (g *mediaGateway) GetTopicByID(ctx context.Context, topicID string) (*gw_response.TopicResponse, error) {
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
	headers := helper.GetHeaders(ctx)
	url := fmt.Sprintf("/api/v2/gateway/topics/%s", topicID)
	resp, err := client.Call("GET", url, nil, headers)
	if err != nil {
		return nil, err
	}

	// parse JSON
	var gwResp gw_response.APIGateWayResponse[*gw_response.TopicResponse]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return nil, fmt.Errorf("unmarshal response fail: %w", err)
	}

	// check status
	if gwResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("call gateway get topic fail: %s", gwResp.Message)
	}

	return gwResp.Data, nil
}

func (g *mediaGateway) GetAllTopicsByOrganization(ctx context.Context, organizationID string) ([]gw_response.TopicResponse, error) {
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
	headers := helper.GetHeaders(ctx)
	url := fmt.Sprintf("/api/v2/gateway/topics/organization/%s", organizationID)
	resp, err := client.Call("GET", url, nil, headers)
	if err != nil {
		return nil, err
	}

	// parse JSON
	var gwResp gw_response.APIGateWayResponse[[]gw_response.TopicResponse]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return nil, fmt.Errorf("unmarshal response fail: %w", err)
	}

	// check status
	if gwResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("call gateway get topics fail: %s", gwResp.Message)
	}

	return gwResp.Data, nil
}

func (g *mediaGateway) GetTopicByStudentID(ctx context.Context, studentID string) ([]gw_response.TopicResponse, error) {
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
	headers := helper.GetHeaders(ctx)
	url := fmt.Sprintf("/api/v2/gateway/topics/student/%s", studentID)
	resp, err := client.Call("GET", url, nil, headers)
	if err != nil {
		return nil, err
	}

	// parse JSON
	var gwResp gw_response.APIGateWayResponse[[]gw_response.TopicResponse]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return nil, fmt.Errorf("unmarshal response fail: %w", err)
	}

	// check status
	if gwResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("call gateway get topics fail: %s", gwResp.Message)
	}

	return gwResp.Data, nil
}
