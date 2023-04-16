package pindo

import "context"

func (c *Client) Send(ctx context.Context, in *SendRequest) (*SendResponse, error) {
	var endpoint = "/sms"

	out := &SendResponse{}

	if err := c.Do(ctx, "POST", endpoint, in, out, nil); err != nil {
		return nil, err
	}
	return out, nil
}
