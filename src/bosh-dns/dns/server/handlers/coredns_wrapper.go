package handlers

import (
	"github.com/miekg/dns"
	"golang.org/x/net/context"
)

type corednsHandlerWrapper struct {
	Next dns.Handler
}

type requestContext struct {
	withMetrics bool
}

func (w corednsHandlerWrapper) ServeDNS(ctx context.Context, writer dns.ResponseWriter, m *dns.Msg) (int, error) {
	reqContext := ctx.Value("indicator").(*reqContext)
	reqContext.withMetrics = false

	w.Next.ServeDNS(writer, m)
	return 0, nil
}

func (w corednsHandlerWrapper) Name() string {
	return "CorednsHandlerWrapper"
}
