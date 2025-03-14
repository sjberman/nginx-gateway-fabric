package interceptor

import (
	"context"
	"fmt"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	authv1 "k8s.io/api/authentication/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	grpcContext "github.com/nginx/nginx-gateway-fabric/internal/mode/static/nginx/agent/grpc/context"
)

const (
	headerUUID = "uuid"
	headerAuth = "authorization"
)

// streamHandler is a struct that implements StreamHandler, allowing the interceptor to replace the context.
type streamHandler struct {
	grpc.ServerStream
	ctx context.Context
}

func (sh *streamHandler) Context() context.Context {
	return sh.ctx
}

type ContextSetter struct {
	k8sClient client.Client
}

func NewContextSetter(k8sClient client.Client) ContextSetter {
	return ContextSetter{k8sClient: k8sClient}
}

func (c ContextSetter) Stream() grpc.StreamServerInterceptor {
	return func(
		srv any,
		ss grpc.ServerStream,
		_ *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		ctx, err := c.validateConnection(ss.Context())
		if err != nil {
			return err
		}
		return handler(srv, &streamHandler{
			ServerStream: ss,
			ctx:          ctx,
		})
	}
}

func (c ContextSetter) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		_ *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		if ctx, err = c.validateConnection(ctx); err != nil {
			return nil, err
		}
		return handler(ctx, req)
	}
}

// validateConnection checks that the connection is valid and returns a new
// context containing information used by the gRPC command/file services.
func (c ContextSetter) validateConnection(ctx context.Context) (context.Context, error) {
	gi, err := getGrpcInfo(ctx)
	if err != nil {
		return nil, err
	}

	return c.validateToken(ctx, gi)
}

func getGrpcInfo(ctx context.Context) (*grpcContext.GrpcInfo, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "no metadata")
	}

	id := md.Get(headerUUID)
	if len(id) == 0 {
		return nil, status.Error(codes.Unauthenticated, "no identity")
	}

	auths := md.Get(headerAuth)
	if len(auths) == 0 {
		return nil, status.Error(codes.Unauthenticated, "no authorization")
	}

	p, ok := peer.FromContext(ctx)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "no peer data")
	}

	addr, ok := p.Addr.(*net.TCPAddr)
	if !ok {
		panic(fmt.Sprintf("address %q was not of type net.TCPAddr", p.Addr.String()))
	}

	return &grpcContext.GrpcInfo{
		SystemID:  id[0],
		Token:     auths[0],
		IPAddress: addr.IP.String(),
	}, nil
}

func (c ContextSetter) validateToken(ctx context.Context, gi *grpcContext.GrpcInfo) (context.Context, error) {
	tokenReview := &authv1.TokenReview{
		Spec: authv1.TokenReviewSpec{
			Token: gi.Token,
		},
	}

	createCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if err := c.k8sClient.Create(createCtx, tokenReview); err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("error creating TokenReview: %v", err))
	}

	if !tokenReview.Status.Authenticated {
		return nil, status.Error(codes.Unauthenticated, "invalid authorization")
	}

	return grpcContext.NewGrpcContext(ctx, *gi), nil
}
