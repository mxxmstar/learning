package http_status_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/mxxmstar/learning/gate_server/gate_config"
	status_def "github.com/mxxmstar/learning/pkg/def/status"
)

type StatusClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewStatusClient(c gate_config.Config, httpClient *http.Client) *StatusClient {
	ss := c.GetStatusServerAddress()
	baseURL := fmt.Sprintf("http://%s", ss)
	return &StatusClient{
		baseURL:    baseURL,
		httpClient: httpClient,
	}
}

// url
const (
	// BaseURL =
	ServiceDiscoveryURL = "/gate/discovery"
)

// 状态码映射
const (
	CodeSuccess             = 200
	CodeBadRequest          = 400
	CodeConflict            = 409
	CodeNotFound            = 404
	CodeInsufficientStorage = 507
	CodeInternalServerError = 500
)

// 业务响应消息映射
var ResponseMessages = map[int]string{
	CodeSuccess:             "Success",
	CodeBadRequest:          "Bad Request",
	CodeConflict:            "Conflict",
	CodeNotFound:            "NotFound",
	CodeInsufficientStorage: "Insufficient Storage",
	CodeInternalServerError: "Internal Server Error",
}

// 服务注册响应消息映射
var ServiceRegisterMessages = map[int]string{
	CodeSuccess:             "Service registered successfully",
	CodeConflict:            "Service ID already exists",
	CodeBadRequest:          "Invalid service name or request parameters",
	CodeInsufficientStorage: "Service registry is full",
	CodeInternalServerError: "Internal server error",
}

// 服务发现响应消息映射
var ServiceDiscoveryMessages = map[int]string{
	CodeSuccess:             "Service discovery successful",
	CodeBadRequest:          "Invalid request parameters",
	CodeNotFound:            "Service not found",
	CodeInternalServerError: "Internal server error",
}

// TODO:封装一个通用的获取实例接口，参数为服务名，策略。返回一个实例

// 获取一个验证服务实例
func GetVerifyServer(c *StatusClient) (*status_def.ServiceInfo, error) {
	req := &status_def.ServiceDiscoveryByTagsRequest{
		ServiceName: "verify_server",
		Strategy:    "Load",
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	response, err := c.httpClient.Post(
		fmt.Sprintf("%s%s/verify", c.baseURL, ServiceDiscoveryURL),
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var res status_def.ServiceDiscoveryByTagsResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}

	return HandleServiceDiscoveryByTagsResponse(&res)
}

// 处理服务发现响应
func HandleServiceDiscoveryByTagsResponse(res *status_def.ServiceDiscoveryByTagsResponse) (*status_def.ServiceInfo, error) {
	switch res.Code {
	case CodeSuccess:
		return res.Services, nil
	default:
		// 使用映射获取标准错误消息
		if message, exists := ServiceDiscoveryMessages[res.Code]; exists {
			return nil, fmt.Errorf("%s: %s", message, res.Message)
		}
		// 如果没有预定义的消息，使用通用错误格式
		return nil, fmt.Errorf("unexpected response code %d: %s", res.Code, res.Message)
	}
}
