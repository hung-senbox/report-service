package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"report-service/helper"
	gw_response "report-service/internal/gateway/dto/response"
	"report-service/pkg/constants"

	"github.com/hashicorp/consul/api"
)

type ClassroomGateway interface {
	GetStudents4ClassroomReport(ctx context.Context, termID string, classroomID string, teacherID string) ([]*gw_response.Student4ClassroomReport, error)
	GetStudentsByClassroomID(ctx context.Context, classroomID string, termID string) ([]*gw_response.Student4ClassroomReport, error)
}

type classroomGateway struct {
	serviceName string
	consul      *api.Client
}

func NewClassroomGateway(serviceName string, consulClient *api.Client) ClassroomGateway {
	return &classroomGateway{
		serviceName: serviceName,
		consul:      consulClient,
	}
}

func (g *classroomGateway) GetStudents4ClassroomReport(ctx context.Context, termID string, classroomID string, teacherID string) ([]*gw_response.Student4ClassroomReport, error) {
	token, ok := ctx.Value(constants.Token).(string)
	if !ok {
		return nil, fmt.Errorf("token not found in context")
	}

	client, err := NewGatewayClient(g.serviceName, token, g.consul, nil)
	if err != nil {
		return nil, fmt.Errorf("init GatewayClient fail: %w", err)
	}

	headers := helper.GetHeaders(ctx)
	resp, err := client.Call("GET", "/api/v1/gateway/classrooms/teacher-assignments?term_id="+termID+"&classroom_id="+classroomID+"&teacher_id="+teacherID+"", nil, headers)
	if err != nil {
		return nil, fmt.Errorf("call API get teacher by user fail: %w", err)
	}

	// Unmarshal response theo format Gateway
	var gwResp gw_response.APIGateWayResponse[[]*gw_response.Student4ClassroomReport]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return nil, fmt.Errorf("unmarshal response fail: %w", err)
	}

	// Check status_code trả về
	if gwResp.StatusCode != 200 {
		return nil, fmt.Errorf("gateway error: %s", gwResp.Message)
	}

	return gwResp.Data, nil
}

func (g *classroomGateway) GetStudentsByClassroomID(ctx context.Context, classroomID string, termID string) ([]*gw_response.Student4ClassroomReport, error) {
	token, ok := ctx.Value(constants.Token).(string)
	if !ok {
		return nil, fmt.Errorf("token not found in context")
	}

	client, err := NewGatewayClient(g.serviceName, token, g.consul, nil)
	if err != nil {
		return nil, fmt.Errorf("init GatewayClient fail: %w", err)
	}

	headers := helper.GetHeaders(ctx)
	resp, err := client.Call("GET", "/api/v1/gateway/classrooms/students?classroom_id="+classroomID+"&term_id="+termID+"", nil, headers)
	if err != nil {
		return nil, fmt.Errorf("call API get teacher by user fail: %w", err)
	}

	// Unmarshal response theo format Gateway
	var gwResp gw_response.APIGateWayResponse[[]*gw_response.Student4ClassroomReport]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return nil, fmt.Errorf("unmarshal response fail: %w", err)
	}

	// Check status_code trả về
	if gwResp.StatusCode != 200 {
		return nil, fmt.Errorf("gateway error: %s", gwResp.Message)
	}

	return gwResp.Data, nil
}
