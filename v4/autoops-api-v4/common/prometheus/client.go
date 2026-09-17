package prometheus

import (
	"context"
	"fmt"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

// Client Prometheus客户端
type Client struct {
	api     v1.API
	address string
}

// QueryResult 查询结果
type QueryResult struct {
	Metric map[string]string `json:"metric"`
	Value  float64           `json:"value"`
	Time   time.Time         `json:"time"`
}

// NewClient 创建Prometheus客户端
func NewClient(address string) (*Client, error) {
	client, err := api.NewClient(api.Config{
		Address: address,
	})
	if err != nil {
		return nil, fmt.Errorf("创建Prometheus客户端失败: %v", err)
	}

	return &Client{
		api:     v1.NewAPI(client),
		address: address,
	}, nil
}

// Query 执行PromQL查询
func (c *Client) Query(ctx context.Context, query string) ([]QueryResult, error) {
	result, warnings, err := c.api.Query(ctx, query, time.Now())
	if err != nil {
		return nil, fmt.Errorf("查询Prometheus失败: %v", err)
	}

	if len(warnings) > 0 {
		fmt.Printf("Prometheus查询警告: %v\n", warnings)
	}

	return c.parseQueryResult(result)
}

// QueryWithTime 在指定时间点执行PromQL查询
func (c *Client) QueryWithTime(ctx context.Context, query string, ts time.Time) ([]QueryResult, error) {
	result, warnings, err := c.api.Query(ctx, query, ts)
	if err != nil {
		return nil, fmt.Errorf("查询Prometheus失败: %v", err)
	}

	if len(warnings) > 0 {
		fmt.Printf("Prometheus查询警告: %v\n", warnings)
	}

	return c.parseQueryResult(result)
}

// QueryRange 范围查询
func (c *Client) QueryRange(ctx context.Context, query string, start, end time.Time, step time.Duration) ([]QueryResult, error) {
	r := v1.Range{
		Start: start,
		End:   end,
		Step:  step,
	}

	result, warnings, err := c.api.QueryRange(ctx, query, r)
	if err != nil {
		return nil, fmt.Errorf("范围查询Prometheus失败: %v", err)
	}

	if len(warnings) > 0 {
		fmt.Printf("Prometheus查询警告: %v\n", warnings)
	}

	return c.parseQueryResult(result)
}

// parseQueryResult 解析查询结果
func (c *Client) parseQueryResult(value model.Value) ([]QueryResult, error) {
	var results []QueryResult

	switch value.Type() {
	case model.ValVector:
		vector := value.(model.Vector)
		for _, sample := range vector {
			metric := make(map[string]string)
			for k, v := range sample.Metric {
				metric[string(k)] = string(v)
			}

			results = append(results, QueryResult{
				Metric: metric,
				Value:  float64(sample.Value),
				Time:   sample.Timestamp.Time(),
			})
		}

	case model.ValMatrix:
		matrix := value.(model.Matrix)
		for _, stream := range matrix {
			metric := make(map[string]string)
			for k, v := range stream.Metric {
				metric[string(k)] = string(v)
			}

			// 对于范围查询，只返回最新的值
			if len(stream.Values) > 0 {
				lastValue := stream.Values[len(stream.Values)-1]
				results = append(results, QueryResult{
					Metric: metric,
					Value:  float64(lastValue.Value),
					Time:   lastValue.Timestamp.Time(),
				})
			}
		}

	case model.ValScalar:
		scalar := value.(*model.Scalar)
		results = append(results, QueryResult{
			Metric: map[string]string{},
			Value:  float64(scalar.Value),
			Time:   scalar.Timestamp.Time(),
		})

	default:
		return nil, fmt.Errorf("不支持的结果类型: %v", value.Type())
	}

	return results, nil
}

// GetAddress 获取Prometheus地址
func (c *Client) GetAddress() string {
	return c.address
}
