package util

import (
	"context"

	"google.golang.org/grpc/metadata"
)

// AttachAuth stamps key onto ctx as an outgoing "authorization: bearer
// <key>" header for the next RPC made with ctx.
func AttachAuth(ctx context.Context, key string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, HeaderAuthorize, TokenPrefix+key)
}

// RelayAuth copies an incoming "authorization" header (e.g. one attached
// by AttachAuth on the cli -> clid leg over the local Unix socket) onto
// the outgoing context, so it can be forwarded verbatim on clid's
// upstream call to the server. No-op if the incoming request carried no
// such header.
func RelayAuth(ctx context.Context) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx
	}
	values := md.Get(HeaderAuthorize)
	if len(values) == 0 {
		return ctx
	}
	return metadata.AppendToOutgoingContext(ctx, HeaderAuthorize, values[0])
}
