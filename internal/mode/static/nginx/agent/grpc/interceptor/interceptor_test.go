package interceptor

import (
	"context"
	"errors"
	"net"
	"testing"

	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	authv1 "k8s.io/api/authentication/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type mockServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (m *mockServerStream) Context() context.Context {
	return m.ctx
}

type mockClient struct {
	client.Client
	createErr     error
	authenticated bool
}

func (m *mockClient) Create(_ context.Context, obj client.Object, _ ...client.CreateOption) error {
	tr, ok := obj.(*authv1.TokenReview)
	if !ok {
		return errors.New("couldn't convert object to TokenReview")
	}
	tr.Status.Authenticated = m.authenticated

	return m.createErr
}

func TestInterceptor(t *testing.T) {
	t.Parallel()

	validMetadata := metadata.New(map[string]string{
		headerUUID: "test-uuid",
		headerAuth: "test-token",
	})
	validPeerData := &peer.Peer{
		Addr: &net.TCPAddr{IP: net.ParseIP("127.0.0.1")},
	}

	tests := []struct {
		md            metadata.MD
		peer          *peer.Peer
		createErr     error
		name          string
		authenticated bool
		expErrCode    codes.Code
	}{
		{
			name:          "valid request",
			md:            validMetadata,
			peer:          validPeerData,
			authenticated: true,
			createErr:     nil,
			expErrCode:    codes.OK,
		},
		{
			name:          "missing metadata",
			peer:          validPeerData,
			authenticated: true,
			createErr:     nil,
			expErrCode:    codes.InvalidArgument,
		},
		{
			name: "missing uuid",
			md: metadata.New(map[string]string{
				headerAuth: "test-token",
			}),
			peer:          validPeerData,
			authenticated: true,
			createErr:     nil,
			expErrCode:    codes.Unauthenticated,
		},
		{
			name: "missing authorization",
			md: metadata.New(map[string]string{
				headerUUID: "test-uuid",
			}),
			peer:          validPeerData,
			authenticated: true,
			createErr:     nil,
			expErrCode:    codes.Unauthenticated,
		},
		{
			name:          "missing peer data",
			md:            validMetadata,
			authenticated: true,
			createErr:     nil,
			expErrCode:    codes.InvalidArgument,
		},
		{
			name:          "tokenreview not created",
			md:            validMetadata,
			peer:          validPeerData,
			authenticated: true,
			createErr:     errors.New("not created"),
			expErrCode:    codes.Internal,
		},
		{
			name:          "tokenreview created and not authenticated",
			md:            validMetadata,
			peer:          validPeerData,
			authenticated: false,
			createErr:     nil,
			expErrCode:    codes.Unauthenticated,
		},
	}

	streamHandler := func(_ any, _ grpc.ServerStream) error {
		return nil
	}

	unaryHandler := func(_ context.Context, _ any) (any, error) {
		return nil, nil //nolint:nilnil // unit test
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			g := NewWithT(t)

			mockK8sClient := &mockClient{
				authenticated: test.authenticated,
				createErr:     test.createErr,
			}
			cs := NewContextSetter(mockK8sClient)

			ctx := context.Background()
			if test.md != nil {
				peerCtx := context.Background()
				if test.peer != nil {
					peerCtx = peer.NewContext(context.Background(), test.peer)
				}
				ctx = metadata.NewIncomingContext(peerCtx, test.md)
			}

			stream := &mockServerStream{ctx: ctx}

			err := cs.Stream()(nil, stream, nil, streamHandler)
			if test.expErrCode != codes.OK {
				g.Expect(status.Code(err)).To(Equal(test.expErrCode))
			} else {
				g.Expect(err).ToNot(HaveOccurred())
			}

			_, err = cs.Unary()(ctx, nil, nil, unaryHandler)
			if test.expErrCode != codes.OK {
				g.Expect(status.Code(err)).To(Equal(test.expErrCode))
			} else {
				g.Expect(err).ToNot(HaveOccurred())
			}
		})
	}
}
