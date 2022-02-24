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

package binaryauthorization

import (
	"context"
	"fmt"
	"math"
	"net/url"
	"time"

	gax "github.com/googleapis/gax-go/v2"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"google.golang.org/api/option/internaloption"
	gtransport "google.golang.org/api/transport/grpc"
	binaryauthorizationpb "google.golang.org/genproto/googleapis/cloud/binaryauthorization/v1beta1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

var newBinauthzManagementServiceV1Beta1ClientHook clientHook

// BinauthzManagementServiceV1Beta1CallOptions contains the retry settings for each method of BinauthzManagementServiceV1Beta1Client.
type BinauthzManagementServiceV1Beta1CallOptions struct {
	GetPolicy      []gax.CallOption
	UpdatePolicy   []gax.CallOption
	CreateAttestor []gax.CallOption
	GetAttestor    []gax.CallOption
	UpdateAttestor []gax.CallOption
	ListAttestors  []gax.CallOption
	DeleteAttestor []gax.CallOption
}

func defaultBinauthzManagementServiceV1Beta1GRPCClientOptions() []option.ClientOption {
	return []option.ClientOption{
		internaloption.WithDefaultEndpoint("binaryauthorization.googleapis.com:443"),
		internaloption.WithDefaultMTLSEndpoint("binaryauthorization.mtls.googleapis.com:443"),
		internaloption.WithDefaultAudience("https://binaryauthorization.googleapis.com/"),
		internaloption.WithDefaultScopes(DefaultAuthScopes()...),
		internaloption.EnableJwtWithScope(),
		option.WithGRPCDialOption(grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(math.MaxInt32))),
	}
}

func defaultBinauthzManagementServiceV1Beta1CallOptions() *BinauthzManagementServiceV1Beta1CallOptions {
	return &BinauthzManagementServiceV1Beta1CallOptions{
		GetPolicy: []gax.CallOption{
			gax.WithRetry(func() gax.Retryer {
				return gax.OnCodes([]codes.Code{
					codes.DeadlineExceeded,
					codes.Unavailable,
				}, gax.Backoff{
					Initial:    100 * time.Millisecond,
					Max:        60000 * time.Millisecond,
					Multiplier: 1.30,
				})
			}),
		},
		UpdatePolicy: []gax.CallOption{
			gax.WithRetry(func() gax.Retryer {
				return gax.OnCodes([]codes.Code{
					codes.DeadlineExceeded,
					codes.Unavailable,
				}, gax.Backoff{
					Initial:    100 * time.Millisecond,
					Max:        60000 * time.Millisecond,
					Multiplier: 1.30,
				})
			}),
		},
		CreateAttestor: []gax.CallOption{},
		GetAttestor: []gax.CallOption{
			gax.WithRetry(func() gax.Retryer {
				return gax.OnCodes([]codes.Code{
					codes.DeadlineExceeded,
					codes.Unavailable,
				}, gax.Backoff{
					Initial:    100 * time.Millisecond,
					Max:        60000 * time.Millisecond,
					Multiplier: 1.30,
				})
			}),
		},
		UpdateAttestor: []gax.CallOption{
			gax.WithRetry(func() gax.Retryer {
				return gax.OnCodes([]codes.Code{
					codes.DeadlineExceeded,
					codes.Unavailable,
				}, gax.Backoff{
					Initial:    100 * time.Millisecond,
					Max:        60000 * time.Millisecond,
					Multiplier: 1.30,
				})
			}),
		},
		ListAttestors: []gax.CallOption{
			gax.WithRetry(func() gax.Retryer {
				return gax.OnCodes([]codes.Code{
					codes.DeadlineExceeded,
					codes.Unavailable,
				}, gax.Backoff{
					Initial:    100 * time.Millisecond,
					Max:        60000 * time.Millisecond,
					Multiplier: 1.30,
				})
			}),
		},
		DeleteAttestor: []gax.CallOption{
			gax.WithRetry(func() gax.Retryer {
				return gax.OnCodes([]codes.Code{
					codes.DeadlineExceeded,
					codes.Unavailable,
				}, gax.Backoff{
					Initial:    100 * time.Millisecond,
					Max:        60000 * time.Millisecond,
					Multiplier: 1.30,
				})
			}),
		},
	}
}

// internalBinauthzManagementServiceV1Beta1Client is an interface that defines the methods availaible from Binary Authorization API.
type internalBinauthzManagementServiceV1Beta1Client interface {
	Close() error
	setGoogleClientInfo(...string)
	Connection() *grpc.ClientConn
	GetPolicy(context.Context, *binaryauthorizationpb.GetPolicyRequest, ...gax.CallOption) (*binaryauthorizationpb.Policy, error)
	UpdatePolicy(context.Context, *binaryauthorizationpb.UpdatePolicyRequest, ...gax.CallOption) (*binaryauthorizationpb.Policy, error)
	CreateAttestor(context.Context, *binaryauthorizationpb.CreateAttestorRequest, ...gax.CallOption) (*binaryauthorizationpb.Attestor, error)
	GetAttestor(context.Context, *binaryauthorizationpb.GetAttestorRequest, ...gax.CallOption) (*binaryauthorizationpb.Attestor, error)
	UpdateAttestor(context.Context, *binaryauthorizationpb.UpdateAttestorRequest, ...gax.CallOption) (*binaryauthorizationpb.Attestor, error)
	ListAttestors(context.Context, *binaryauthorizationpb.ListAttestorsRequest, ...gax.CallOption) *AttestorIterator
	DeleteAttestor(context.Context, *binaryauthorizationpb.DeleteAttestorRequest, ...gax.CallOption) error
}

// BinauthzManagementServiceV1Beta1Client is a client for interacting with Binary Authorization API.
// Methods, except Close, may be called concurrently. However, fields must not be modified concurrently with method calls.
//
// Google Cloud Management Service for Binary Authorization admission policies
// and attestation authorities.
//
// This API implements a REST model with the following objects:
//
//   Policy
//
//   Attestor
type BinauthzManagementServiceV1Beta1Client struct {
	// The internal transport-dependent client.
	internalClient internalBinauthzManagementServiceV1Beta1Client

	// The call options for this service.
	CallOptions *BinauthzManagementServiceV1Beta1CallOptions
}

// Wrapper methods routed to the internal client.

// Close closes the connection to the API service. The user should invoke this when
// the client is no longer required.
func (c *BinauthzManagementServiceV1Beta1Client) Close() error {
	return c.internalClient.Close()
}

// setGoogleClientInfo sets the name and version of the application in
// the `x-goog-api-client` header passed on each request. Intended for
// use by Google-written clients.
func (c *BinauthzManagementServiceV1Beta1Client) setGoogleClientInfo(keyval ...string) {
	c.internalClient.setGoogleClientInfo(keyval...)
}

// Connection returns a connection to the API service.
//
// Deprecated.
func (c *BinauthzManagementServiceV1Beta1Client) Connection() *grpc.ClientConn {
	return c.internalClient.Connection()
}

// GetPolicy a policy specifies the attestors that must attest to
// a container image, before the project is allowed to deploy that
// image. There is at most one policy per project. All image admission
// requests are permitted if a project has no policy.
//
// Gets the policy for this project. Returns a default
// policy if the project does not have one.
func (c *BinauthzManagementServiceV1Beta1Client) GetPolicy(ctx context.Context, req *binaryauthorizationpb.GetPolicyRequest, opts ...gax.CallOption) (*binaryauthorizationpb.Policy, error) {
	return c.internalClient.GetPolicy(ctx, req, opts...)
}

// UpdatePolicy creates or updates a project’s policy, and returns a copy of the
// new policy. A policy is always updated as a whole, to avoid race
// conditions with concurrent policy enforcement (or management!)
// requests. Returns NOT_FOUND if the project does not exist, INVALID_ARGUMENT
// if the request is malformed.
func (c *BinauthzManagementServiceV1Beta1Client) UpdatePolicy(ctx context.Context, req *binaryauthorizationpb.UpdatePolicyRequest, opts ...gax.CallOption) (*binaryauthorizationpb.Policy, error) {
	return c.internalClient.UpdatePolicy(ctx, req, opts...)
}

// CreateAttestor creates an attestor, and returns a copy of the new
// attestor. Returns NOT_FOUND if the project does not exist,
// INVALID_ARGUMENT if the request is malformed, ALREADY_EXISTS if the
// attestor already exists.
func (c *BinauthzManagementServiceV1Beta1Client) CreateAttestor(ctx context.Context, req *binaryauthorizationpb.CreateAttestorRequest, opts ...gax.CallOption) (*binaryauthorizationpb.Attestor, error) {
	return c.internalClient.CreateAttestor(ctx, req, opts...)
}

// GetAttestor gets an attestor.
// Returns NOT_FOUND if the attestor does not exist.
func (c *BinauthzManagementServiceV1Beta1Client) GetAttestor(ctx context.Context, req *binaryauthorizationpb.GetAttestorRequest, opts ...gax.CallOption) (*binaryauthorizationpb.Attestor, error) {
	return c.internalClient.GetAttestor(ctx, req, opts...)
}

// UpdateAttestor updates an attestor.
// Returns NOT_FOUND if the attestor does not exist.
func (c *BinauthzManagementServiceV1Beta1Client) UpdateAttestor(ctx context.Context, req *binaryauthorizationpb.UpdateAttestorRequest, opts ...gax.CallOption) (*binaryauthorizationpb.Attestor, error) {
	return c.internalClient.UpdateAttestor(ctx, req, opts...)
}

// ListAttestors lists attestors.
// Returns INVALID_ARGUMENT if the project does not exist.
func (c *BinauthzManagementServiceV1Beta1Client) ListAttestors(ctx context.Context, req *binaryauthorizationpb.ListAttestorsRequest, opts ...gax.CallOption) *AttestorIterator {
	return c.internalClient.ListAttestors(ctx, req, opts...)
}

// DeleteAttestor deletes an attestor. Returns NOT_FOUND if the
// attestor does not exist.
func (c *BinauthzManagementServiceV1Beta1Client) DeleteAttestor(ctx context.Context, req *binaryauthorizationpb.DeleteAttestorRequest, opts ...gax.CallOption) error {
	return c.internalClient.DeleteAttestor(ctx, req, opts...)
}

// binauthzManagementServiceV1Beta1GRPCClient is a client for interacting with Binary Authorization API over gRPC transport.
//
// Methods, except Close, may be called concurrently. However, fields must not be modified concurrently with method calls.
type binauthzManagementServiceV1Beta1GRPCClient struct {
	// Connection pool of gRPC connections to the service.
	connPool gtransport.ConnPool

	// flag to opt out of default deadlines via GOOGLE_API_GO_EXPERIMENTAL_DISABLE_DEFAULT_DEADLINE
	disableDeadlines bool

	// Points back to the CallOptions field of the containing BinauthzManagementServiceV1Beta1Client
	CallOptions **BinauthzManagementServiceV1Beta1CallOptions

	// The gRPC API client.
	binauthzManagementServiceV1Beta1Client binaryauthorizationpb.BinauthzManagementServiceV1Beta1Client

	// The x-goog-* metadata to be sent with each request.
	xGoogMetadata metadata.MD
}

// NewBinauthzManagementServiceV1Beta1Client creates a new binauthz management service v1 beta1 client based on gRPC.
// The returned client must be Closed when it is done being used to clean up its underlying connections.
//
// Google Cloud Management Service for Binary Authorization admission policies
// and attestation authorities.
//
// This API implements a REST model with the following objects:
//
//   Policy
//
//   Attestor
func NewBinauthzManagementServiceV1Beta1Client(ctx context.Context, opts ...option.ClientOption) (*BinauthzManagementServiceV1Beta1Client, error) {
	clientOpts := defaultBinauthzManagementServiceV1Beta1GRPCClientOptions()
	if newBinauthzManagementServiceV1Beta1ClientHook != nil {
		hookOpts, err := newBinauthzManagementServiceV1Beta1ClientHook(ctx, clientHookParams{})
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
	client := BinauthzManagementServiceV1Beta1Client{CallOptions: defaultBinauthzManagementServiceV1Beta1CallOptions()}

	c := &binauthzManagementServiceV1Beta1GRPCClient{
		connPool:                               connPool,
		disableDeadlines:                       disableDeadlines,
		binauthzManagementServiceV1Beta1Client: binaryauthorizationpb.NewBinauthzManagementServiceV1Beta1Client(connPool),
		CallOptions:                            &client.CallOptions,
	}
	c.setGoogleClientInfo()

	client.internalClient = c

	return &client, nil
}

// Connection returns a connection to the API service.
//
// Deprecated.
func (c *binauthzManagementServiceV1Beta1GRPCClient) Connection() *grpc.ClientConn {
	return c.connPool.Conn()
}

// setGoogleClientInfo sets the name and version of the application in
// the `x-goog-api-client` header passed on each request. Intended for
// use by Google-written clients.
func (c *binauthzManagementServiceV1Beta1GRPCClient) setGoogleClientInfo(keyval ...string) {
	kv := append([]string{"gl-go", versionGo()}, keyval...)
	kv = append(kv, "gapic", getVersionClient(), "gax", gax.Version, "grpc", grpc.Version)
	c.xGoogMetadata = metadata.Pairs("x-goog-api-client", gax.XGoogHeader(kv...))
}

// Close closes the connection to the API service. The user should invoke this when
// the client is no longer required.
func (c *binauthzManagementServiceV1Beta1GRPCClient) Close() error {
	return c.connPool.Close()
}

func (c *binauthzManagementServiceV1Beta1GRPCClient) GetPolicy(ctx context.Context, req *binaryauthorizationpb.GetPolicyRequest, opts ...gax.CallOption) (*binaryauthorizationpb.Policy, error) {
	if _, ok := ctx.Deadline(); !ok && !c.disableDeadlines {
		cctx, cancel := context.WithTimeout(ctx, 600000*time.Millisecond)
		defer cancel()
		ctx = cctx
	}
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "name", url.QueryEscape(req.GetName())))

	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append((*c.CallOptions).GetPolicy[0:len((*c.CallOptions).GetPolicy):len((*c.CallOptions).GetPolicy)], opts...)
	var resp *binaryauthorizationpb.Policy
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.binauthzManagementServiceV1Beta1Client.GetPolicy(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *binauthzManagementServiceV1Beta1GRPCClient) UpdatePolicy(ctx context.Context, req *binaryauthorizationpb.UpdatePolicyRequest, opts ...gax.CallOption) (*binaryauthorizationpb.Policy, error) {
	if _, ok := ctx.Deadline(); !ok && !c.disableDeadlines {
		cctx, cancel := context.WithTimeout(ctx, 600000*time.Millisecond)
		defer cancel()
		ctx = cctx
	}
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "policy.name", url.QueryEscape(req.GetPolicy().GetName())))

	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append((*c.CallOptions).UpdatePolicy[0:len((*c.CallOptions).UpdatePolicy):len((*c.CallOptions).UpdatePolicy)], opts...)
	var resp *binaryauthorizationpb.Policy
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.binauthzManagementServiceV1Beta1Client.UpdatePolicy(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *binauthzManagementServiceV1Beta1GRPCClient) CreateAttestor(ctx context.Context, req *binaryauthorizationpb.CreateAttestorRequest, opts ...gax.CallOption) (*binaryauthorizationpb.Attestor, error) {
	if _, ok := ctx.Deadline(); !ok && !c.disableDeadlines {
		cctx, cancel := context.WithTimeout(ctx, 600000*time.Millisecond)
		defer cancel()
		ctx = cctx
	}
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "parent", url.QueryEscape(req.GetParent())))

	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append((*c.CallOptions).CreateAttestor[0:len((*c.CallOptions).CreateAttestor):len((*c.CallOptions).CreateAttestor)], opts...)
	var resp *binaryauthorizationpb.Attestor
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.binauthzManagementServiceV1Beta1Client.CreateAttestor(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *binauthzManagementServiceV1Beta1GRPCClient) GetAttestor(ctx context.Context, req *binaryauthorizationpb.GetAttestorRequest, opts ...gax.CallOption) (*binaryauthorizationpb.Attestor, error) {
	if _, ok := ctx.Deadline(); !ok && !c.disableDeadlines {
		cctx, cancel := context.WithTimeout(ctx, 600000*time.Millisecond)
		defer cancel()
		ctx = cctx
	}
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "name", url.QueryEscape(req.GetName())))

	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append((*c.CallOptions).GetAttestor[0:len((*c.CallOptions).GetAttestor):len((*c.CallOptions).GetAttestor)], opts...)
	var resp *binaryauthorizationpb.Attestor
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.binauthzManagementServiceV1Beta1Client.GetAttestor(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *binauthzManagementServiceV1Beta1GRPCClient) UpdateAttestor(ctx context.Context, req *binaryauthorizationpb.UpdateAttestorRequest, opts ...gax.CallOption) (*binaryauthorizationpb.Attestor, error) {
	if _, ok := ctx.Deadline(); !ok && !c.disableDeadlines {
		cctx, cancel := context.WithTimeout(ctx, 600000*time.Millisecond)
		defer cancel()
		ctx = cctx
	}
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "attestor.name", url.QueryEscape(req.GetAttestor().GetName())))

	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append((*c.CallOptions).UpdateAttestor[0:len((*c.CallOptions).UpdateAttestor):len((*c.CallOptions).UpdateAttestor)], opts...)
	var resp *binaryauthorizationpb.Attestor
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.binauthzManagementServiceV1Beta1Client.UpdateAttestor(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *binauthzManagementServiceV1Beta1GRPCClient) ListAttestors(ctx context.Context, req *binaryauthorizationpb.ListAttestorsRequest, opts ...gax.CallOption) *AttestorIterator {
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "parent", url.QueryEscape(req.GetParent())))

	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append((*c.CallOptions).ListAttestors[0:len((*c.CallOptions).ListAttestors):len((*c.CallOptions).ListAttestors)], opts...)
	it := &AttestorIterator{}
	req = proto.Clone(req).(*binaryauthorizationpb.ListAttestorsRequest)
	it.InternalFetch = func(pageSize int, pageToken string) ([]*binaryauthorizationpb.Attestor, string, error) {
		resp := &binaryauthorizationpb.ListAttestorsResponse{}
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
			resp, err = c.binauthzManagementServiceV1Beta1Client.ListAttestors(ctx, req, settings.GRPC...)
			return err
		}, opts...)
		if err != nil {
			return nil, "", err
		}

		it.Response = resp
		return resp.GetAttestors(), resp.GetNextPageToken(), nil
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

func (c *binauthzManagementServiceV1Beta1GRPCClient) DeleteAttestor(ctx context.Context, req *binaryauthorizationpb.DeleteAttestorRequest, opts ...gax.CallOption) error {
	if _, ok := ctx.Deadline(); !ok && !c.disableDeadlines {
		cctx, cancel := context.WithTimeout(ctx, 600000*time.Millisecond)
		defer cancel()
		ctx = cctx
	}
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "name", url.QueryEscape(req.GetName())))

	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append((*c.CallOptions).DeleteAttestor[0:len((*c.CallOptions).DeleteAttestor):len((*c.CallOptions).DeleteAttestor)], opts...)
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		_, err = c.binauthzManagementServiceV1Beta1Client.DeleteAttestor(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	return err
}

// AttestorIterator manages a stream of *binaryauthorizationpb.Attestor.
type AttestorIterator struct {
	items    []*binaryauthorizationpb.Attestor
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
	InternalFetch func(pageSize int, pageToken string) (results []*binaryauthorizationpb.Attestor, nextPageToken string, err error)
}

// PageInfo supports pagination. See the google.golang.org/api/iterator package for details.
func (it *AttestorIterator) PageInfo() *iterator.PageInfo {
	return it.pageInfo
}

// Next returns the next result. Its second return value is iterator.Done if there are no more
// results. Once Next returns Done, all subsequent calls will return Done.
func (it *AttestorIterator) Next() (*binaryauthorizationpb.Attestor, error) {
	var item *binaryauthorizationpb.Attestor
	if err := it.nextFunc(); err != nil {
		return item, err
	}
	item = it.items[0]
	it.items = it.items[1:]
	return item, nil
}

func (it *AttestorIterator) bufLen() int {
	return len(it.items)
}

func (it *AttestorIterator) takeBuf() interface{} {
	b := it.items
	it.items = nil
	return b
}
