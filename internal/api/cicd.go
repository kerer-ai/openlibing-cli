package api

import (
	"fmt"
	"io"
)

// GetPipelineDetail fetches pipeline run detail from openlibing-cicd.
func (c *Client) GetPipelineDetail(projectID string, limit int) ([]byte, error) {
	resp, err := c.Do("GET", "gateway/openlibing-cicd/project/pipeline/pipeline-run/detail", map[string]string{
		"projectId": projectID,
		"pageSize":  fmt.Sprintf("%d", limit),
	}, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("get pipeline detail: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	return io.ReadAll(resp.Body)
}
