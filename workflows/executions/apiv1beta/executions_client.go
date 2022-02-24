// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by protoc-gen-go_gapic. DO NOT EDIT.

package executions

import (
	"context"
	"fmt"
	"math"
	"net/url"

	gax "github.com/googleapis/gax-go/v2"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"google.golang.org/api/option/internaloption"
	gtransport "google.golang.org/api/transport/grpc"
	executionspb "google.golang.org/genproto/googleapis/cloud/workflows/executions/v1beta"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

var newClientHook clientHook

// CallOptions contains the retry settings for each method of Client.
type CallOptions struct {
	ListExecutions  []gax.CallOption
	CreateExecution []gax.CallOption
	GetExecution    []gax.CallOption
	CancelExecution []gax.CallOption
}

func defaultGRPCClientOptions() []option.ClientOption {
	return []option.ClientOption{
		internaloption.WithDefaultEndpoint("workflowexecutions.googleapis.com:443"),
		internaloption.WithDefaultMTLSEndpoint("workflowexecutions.mtls.googleapis.com:443"),
		internaloption.WithDefaultAudience("https://workflowexecutions.googleapis.com/"),
		internaloption.WithDefaultScopes(DefaultAuthScopes()...),
		internaloption.EnableJwtWithScope(),
		option.WithGRPCDialOption(grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(math.MaxInt32))),
	}
}

func defaultCallOptions() *CallOptions {
	return &CallOptions{
		ListExecutions:  []gax.CallOption{},
		CreateExecution: []gax.CallOption{},
		GetExecution:    []gax.CallOption{},
		CancelExecution: []gax.CallOption{},
	}
}

// internalClient is an interface that defines the methods availaible from Workflow Executions API.
type internalClient interface {
	Close() error
	setGoogleClientInfo(...string)
	Connection() *grpc.ClientConn
	ListExecutions(context.Context, *executionspb.ListExecutionsRequest, ...gax.CallOption) *ExecutionIterator
	CreateExecution(context.Context, *executionspb.CreateExecutionRequest, ...gax.CallOption) (*executionspb.Execution, error)
	GetExecution(context.Context, *executionspb.GetExecutionRequest, ...gax.CallOption) (*executionspb.Execution, error)
	CancelExecution(context.Context, *executionspb.CancelExecutionRequest, ...gax.CallOption) (*executionspb.Execution, error)
}

// Client is a client for interacting with Workflow Executions API.
// Methods, except Close, may be called concurrently. However, fields must not be modified concurrently with method calls.
//
// Executions is used to start and manage running instances of
// Workflows called executions.
type Client struct {
	// The internal transport-dependent client.
	internalClient internalClient

	// The call options for this service.
	CallOptions *CallOptions
}

// Wrapper methods routed to the internal client.

// Close closes the connection to the API service. The user should invoke this when
// the client is no longer required.
func (c *Client) Close() error {
	return c.internalClient.Close()
}

// setGoogleClientInfo sets the name and version of the application in
// the `x-goog-api-client` header passed on each request. Intended for
// use by Google-written clients.
func (c *Client) setGoogleClientInfo(keyval ...string) {
	c.internalClient.setGoogleClientInfo(keyval...)
}

// Connection returns a connection to the API service.
//
// Deprecated.
func (c *Client) Connection() *grpc.ClientConn {
	return c.internalClient.Connection()
}

// ListExecutions returns a list of executions which belong to the workflow with
// the given name. The method returns executions of all workflow
// revisions. Returned executions are ordered by their start time (newest
// first).
func (c *Client) ListExecutions(ctx context.Context, req *executionspb.ListExecutionsRequest, opts ...gax.CallOption) *ExecutionIterator {
	return c.internalClient.ListExecutions(ctx, req, opts...)
}

// CreateExecution creates a new execution using the latest revision of the given workflow.
func (c *Client) CreateExecution(ctx context.Context, req *executionspb.CreateExecutionRequest, opts ...gax.CallOption) (*executionspb.Execution, error) {
	return c.internalClient.CreateExecution(ctx, req, opts...)
}

// GetExecution returns an execution of the given name.
func (c *Client) GetExecution(ctx context.Context, req *executionspb.GetExecutionRequest, opts ...gax.CallOption) (*executionspb.Execution, error) {
	return c.internalClient.GetExecution(ctx, req, opts...)
}

// CancelExecution cancels an execution of the given name.
func (c *Client) CancelExecution(ctx context.Context, req *executionspb.CancelExecutionRequest, opts ...gax.CallOption) (*executionspb.Execution, error) {
	return c.internalClient.CancelExecution(ctx, req, opts...)
}

// gRPCClient is a client for interacting with Workflow Executions API over gRPC transport.
//
// Methods, except Close, may be called concurrently. However, fields must not be modified concurrently with method calls.
type gRPCClient struct {
	// Connection pool of gRPC connections to the service.
	connPool gtransport.ConnPool

	// flag to opt out of default deadlines via GOOGLE_API_GO_EXPERIMENTAL_DISABLE_DEFAULT_DEADLINE
	disableDeadlines bool

	// Points back to the CallOptions field of the containing Client
	CallOptions **CallOptions

	// The gRPC API client.
	client executionspb.ExecutionsClient

	// The x-goog-* metadata to be sent with each request.
	xGoogMetadata metadata.MD
}

// NewClient creates a new executions client based on gRPC.
// The returned client must be Closed when it is done being used to clean up its underlying connections.
//
// Executions is used to start and manage running instances of
// Workflows called executions.
func NewClient(ctx context.Context, opts ...option.ClientOption) (*Client, error) {
	clientOpts := defaultGRPCClientOptions()
	if newClientHook != nil {
		hookOpts, err := newClientHook(ctx, clientHookParams{})
		if err != nil {
			return nil, err
		}
		clientOpts = append(clientOpts, hookOpts...)
	}

	disableDeadlines, err := checkDisableDeadlines()
	if err != nil {
		return nil, err
	}

	connPool, err := gtransport.DialPool(ctx, append(clientOpts, opts...)...)
	if err != nil {
		return nil, err
	}
	client := Client{CallOptions: defaultCallOptions()}

	c := &gRPCClient{
		connPool:         connPool,
		disableDeadlines: disableDeadlines,
		client:           executionspb.NewExecutionsClient(connPool),
		CallOptions:      &client.CallOptions,
	}
	c.setGoogleClientInfo()

	client.internalClient = c

	return &client, nil
}

// Connection returns a connection to the API service.
//
// Deprecated.
func (c *gRPCClient) Connection() *grpc.ClientConn {
	return c.connPool.Conn()
}

// setGoogleClientInfo sets the name and version of the application in
// the `x-goog-api-client` header passed on each request. Intended for
// use by Google-written clients.
func (c *gRPCClient) setGoogleClientInfo(keyval ...string) {
	kv := append([]string{"gl-go", versionGo()}, keyval...)
	kv = append(kv, "gapic", getVersionClient(), "gax", gax.Version, "grpc", grpc.Version)
	c.xGoogMetadata = metadata.Pairs("x-goog-api-client", gax.XGoogHeader(kv...))
}

// Close closes the connection to the API service. The user should invoke this when
// the client is no longer required.
func (c *gRPCClient) Close() error {
	return c.connPool.Close()
}

func (c *gRPCClient) ListExecutions(ctx context.Context, req *executionspb.ListExecutionsRequest, opts ...gax.CallOption) *ExecutionIterator {
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "parent", url.QueryEscape(req.GetParent())))

	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append((*c.CallOptions).ListExecutions[0:len((*c.CallOptions).ListExecutions):len((*c.CallOptions).ListExecutions)], opts...)
	it := &ExecutionIterator{}
	req = proto.Clone(req).(*executionspb.ListExecutionsRequest)
	it.InternalFetch = func(pageSize int, pageToken string) ([]*executionspb.Execution, string, error) {
		resp := &executionspb.ListExecutionsResponse{}
		if pageToken != "" {
			req.PageToken = pageToken
		}
		if pageSize > math.MaxInt32 {
			req.PageSize = math.MaxInt32
		} else if pageSize != 0 {
			req.PageSize = int32(pageSize)
		}
		err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
			var err error
			resp, err = c.client.ListExecutions(ctx, req, settings.GRPC...)
			return err
		}, opts...)
		if err != nil {
			return nil, "", err
		}

		it.Response = resp
		return resp.GetExecutions(), resp.GetNextPageToken(), nil
	}
	fetch := func(pageSize int, pageToken string) (string, error) {
		items, nextPageToken, err := it.InternalFetch(pageSize, pageToken)
		if err != nil {
			return "", err
		}
		it.items = append(it.items, items...)
		return nextPageToken, nil
	}

	it.pageInfo, it.nextFunc = iterator.NewPageInfo(fetch, it.bufLen, it.takeBuf)
	it.pageInfo.MaxSize = int(req.GetPageSize())
	it.pageInfo.Token = req.GetPageToken()

	return it
}

func (c *gRPCClient) CreateExecution(ctx context.Context, req *executionspb.CreateExecutionRequest, opts ...gax.CallOption) (*executionspb.Execution, error) {
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "parent", url.QueryEscape(req.GetParent())))

	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append((*c.CallOptions).CreateExecution[0:len((*c.CallOptions).CreateExecution):len((*c.CallOptions).CreateExecution)], opts...)
	var resp *executionspb.Execution
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.client.CreateExecution(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *gRPCClient) GetExecution(ctx context.Context, req *executionspb.GetExecutionRequest, opts ...gax.CallOption) (*executionspb.Execution, error) {
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "name", url.QueryEscape(req.GetName())))

	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append((*c.CallOptions).GetExecution[0:len((*c.CallOptions).GetExecution):len((*c.CallOptions).GetExecution)], opts...)
	var resp *executionspb.Execution
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.client.GetExecution(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *gRPCClient) CancelExecution(ctx context.Context, req *executionspb.CancelExecutionRequest, opts ...gax.CallOption) (*executionspb.Execution, error) {
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "name", url.QueryEscape(req.GetName())))

	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append((*c.CallOptions).CancelExecution[0:len((*c.CallOptions).CancelExecution):len((*c.CallOptions).CancelExecution)], opts...)
	var resp *executionspb.Execution
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.client.CancelExecution(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// ExecutionIterator manages a stream of *executionspb.Execution.
type ExecutionIterator struct {
	items    []*executionspb.Execution
	pageInfo *iterator.PageInfo
	nextFunc func() error

	// Response is the raw response for the current page.
	// It must be cast to the RPC response type.
	// Calling Next() or InternalFetch() updates this value.
	Response interface{}

	// InternalFetch is for use by the Google Cloud Libraries only.
	// It is not part of the stable interface of this package.
	//
	// InternalFetch returns results from a single call to the underlying RPC.
	// The number of results is no greater than pageSize.
	// If there are no more results, nextPageToken is empty and err is nil.
	InternalFetch func(pageSize int, pageToken string) (results []*executionspb.Execution, nextPageToken string, err error)
}

// PageInfo supports pagination. See the google.golang.org/api/iterator package for details.
func (it *ExecutionIterator) PageInfo() *iterator.PageInfo {
	return it.pageInfo
}

// Next returns the next result. Its second return value is iterator.Done if there are no more
// results. Once Next returns Done, all subsequent calls will return Done.
func (it *ExecutionIterator) Next() (*executionspb.Execution, error) {
	var item *executionspb.Execution
	if err := it.nextFunc(); err != nil {
		return item, err
	}
	item = it.items[0]
	it.items = it.items[1:]
	return item, nil
}

func (it *ExecutionIterator) bufLen() int {
	return len(it.items)
}

func (it *ExecutionIterator) takeBuf() interface{} {
	b := it.items
	it.items = nil
	return b
}
