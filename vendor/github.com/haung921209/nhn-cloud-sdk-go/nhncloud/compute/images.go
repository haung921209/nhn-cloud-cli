package compute

import (
	"context"
	"fmt"
	"net/url"
	"sort"
)

func (c *Client) ListImages(ctx context.Context) (*ListImagesOutput, error) {
	return c.ListImagesWithFilter(ctx, nil)
}

// ListImagesWithFilter calls /images/detail with arbitrary query-string
// parameters appended. Useful for NHN-specific filters like the NKS-only
// view: pass {"nhncloud_allow_nks_cpu_flavor": "true", "visibility": "public"}
// (per docs/api-specs/container/nks.md "베이스 이미지 UUID"). When `params` is
// nil or empty, the request is identical to plain ListImages.
func (c *Client) ListImagesWithFilter(ctx context.Context, params map[string]string) (*ListImagesOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	endpoint := "/images/detail"
	if len(params) > 0 {
		// Sort keys so the query string is deterministic — easier to log,
		// cache-friendly, and predictable in tests.
		keys := make([]string, 0, len(params))
		for k := range params {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		v := url.Values{}
		for _, k := range keys {
			v.Set(k, params[k])
		}
		endpoint += "?" + v.Encode()
	}

	var out ListImagesOutput
	if err := c.httpClient.GET(ctx, endpoint, &out); err != nil {
		return nil, fmt.Errorf("list images: %w", err)
	}
	return &out, nil
}
