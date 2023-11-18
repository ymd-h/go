// request package
package request

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/ymd-h/go/request/json"
)

type (
	// Interface for `*http.Client` implementation
	IHttpClient interface {
		Do(*http.Request) (*http.Response, error)
	}

	// Interface for Request Body Encoder
	IBodyEncoder interface {
		Encode(any) (io.Reader, error)
		ContentType() string
	}

	// Interface for Response Body Decoder
	IBodyDecoder interface {
		Decode(io.Reader, any) error
	}

	// Client class
	Client struct {
		client IHttpClient
		encoder IBodyEncoder
		decoder IBodyDecoder
	}

	// Response class
	Response struct {
		StatusCode int
		Body any
	}
)

var (
	// Default Client
	DefaultClient = &Client{
		client: http.DefaultClient,
		encoder: json.Encoder{},
		decoder: json.Decoder{},
	}
)


// Create new Client
//
// # Arguments
// * `client`: `IClient` - `*http.Client`
//
// # Returns
// * `*Client` - Created Client
func NewClient(
	client IHttpClient,
	encoder IBodyEncoder,
	decoder IBodyDecoder,
) *Client {
	return &Client{client: client, encoder: encoder, decoder: decoder}
}


func (c *Client) newHttpReqest(
	ctx context.Context,
	method, url string,
	request any,
) (*http.Request, error) {
	if request == nil {
		return http.NewRequestWithContext(ctx, method, url, nil)
	}

	body, err := c.encoder.Encode(request)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", c.encoder.ContentType())
	return req, nil
}


// FetchWithContext
//
// # Arguments
// * `ctx`: `context.Context` - Request Context
// * `method`: `string` - http method name
// * `url`: `string` - Target URL
// * `request`: `any` - Request Body Struct annotated with JSON or `nil`
// * `response`: `any` - Pointer to Response Body Struct annotated with JSON or `ResponseDispatcher`
//
// # Returns
// * `*Response` - Response
// * `error` - Error
func (c *Client) FetchWithContext(
	ctx context.Context,
	method, url string,
	request, response any,
) (*Response, error) {
	s := fmt.Sprintf("%s at %s", method, url)

	req, err := c.newHttpReqest(ctx, method, url, request)
	if err != nil {
		return nil, fmt.Errorf("Fail to Create New Request for %s: %w", s, err)
	}

	res, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Fail to %s: %v", s, err)
	}
	defer res.Body.Close()

	ret := &Response{StatusCode: res.StatusCode, Body: nil}
	if response == nil {
		return ret, nil
	}

	if rd, ok := response.(*ResponseDispatcher); ok {
		response, err = rd.Dispatch(ret.StatusCode)
		if err != nil {
			return ret, fmt.Errorf(
				"Fail to Dispatch Response (StatusCode: %d) for %s: %w",
				ret.StatusCode, s, err)
		}
	}

	err = c.decoder.Decode(res.Body, response)
	if err != nil {
		return ret, fmt.Errorf(
			"Fail to Decode Response (StatucCode: %d) for %s: %w",
			ret.StatusCode, s, err)
	}

	ret.Body = response
	return ret, nil
}

func (c *Client) Fetch(method, url string, request, response any) (*Response, error) {
	return c.FetchWithContext(context.Background(), method, url, request, response)
}

func (c *Client) GetWithContext(
	ctx context.Context,
	url string,
	response any,
) (*Response, error) {
	return c.FetchWithContext(ctx, http.MethodGet, url, nil, response)
}

func (c *Client) Get(url string, response any) (*Response, error) {
	return c.GetWithContext(context.Background(), url, response)
}

func (c *Client) HeadWithContext(ctx context.Context, url string) (*Response, error) {
	return c.FetchWithContext(ctx, http.MethodHead, url, nil, nil)
}

func (c *Client) Head(url string) (*Response, error) {
	return c.HeadWithContext(context.Background(), url)
}

func (c *Client) PostWithContext(
	ctx context.Context,
	url string,
	request, response any,
) (*Response, error) {
	return c.FetchWithContext(ctx, http.MethodPost, url, request, response)
}

func (c *Client) Post(url string, request, response any) (*Response, error) {
	return c.PostWithContext(context.Background(), url, request, response)
}

func (c *Client) PutWithContext(
	ctx context.Context,
	url string,
	request, response any,
) (*Response, error) {
	return c.FetchWithContext(ctx, http.MethodPut, url, request, response)
}

func (c *Client) Put(url string, request, response any) (*Response, error) {
	return c.PutWithContext(context.Background(), url, request, response)
}

func (c *Client) DeleteWitchContext(
	ctx context.Context,
	url string,
	request, response any,
) (*Response, error) {
	return c.FetchWithContext(ctx, http.MethodDelete, url, request, response)
}

func (c *Client) Delete(url string, request, response any) (*Response, error) {
	return c.DeleteWitchContext(context.Background(), url, request, response)
}

func (c *Client) PatchWithContext(
	ctx context.Context,
	url string,
	request, response any,
) (*Response, error) {
	return c.FetchWithContext(ctx, http.MethodPatch, url, request, response)
}

func (c *Client) Patch(url string, request, response any) (*Response, error) {
	return c.PatchWithContext(context.Background(), url, request, response)
}

func FetchWithContext(
	ctx context.Context,
	method, url string,
	request, response any,
) (*Response, error) {
	return DefaultClient.FetchWithContext(ctx, method, url, request, response)
}

func Fetch(method, url string, request, response any) (*Response, error) {
	return FetchWithContext(context.Background(), method, url, request, response)
}

func GetWithContext(
	ctx context.Context,
	url string,
	response any,
) (*Response, error) {
	return DefaultClient.GetWithContext(ctx, url, response)
}

func Get(url string, response any) (*Response, error) {
	return GetWithContext(context.Background(), url, response)
}

func HeadWithContext(ctx context.Context, url string) (*Response, error) {
	return DefaultClient.HeadWithContext(ctx, url)
}

func Head(url string) (*Response, error) {
	return HeadWithContext(context.Background(), url)
}

func PostWithContext(
	ctx context.Context,
	url string,
	request, response any,
) (*Response, error) {
	return DefaultClient.PostWithContext(ctx, url, request, response)
}

func Post(url string, request, response any) (*Response, error) {
	return PostWithContext(context.Background(), url, request, response)
}

func PutWithContext(
	ctx context.Context,
	url string,
	request, response any,
) (*Response, error) {
	return DefaultClient.PutWithContext(ctx, url, request, response)
}

func Put(url string, request, response any) (*Response, error) {
	return PutWithContext(context.Background(), url, request, response)
}

func DeleteWitchContext(
	ctx context.Context,
	url string,
	request, response any,
) (*Response, error) {
	return DefaultClient.DeleteWitchContext(ctx, url, request, response)
}

func Delete(url string, request, response any) (*Response, error) {
	return DeleteWitchContext(context.Background(), url, request, response)
}

func PatchWithContext(
	ctx context.Context,
	url string,
	request, response any,
) (*Response, error) {
	return DefaultClient.PatchWithContext(ctx, url, request, response)
}

func Patch(url string, request, response any) (*Response, error) {
	return PatchWithContext(context.Background(), url, request, response)
}
