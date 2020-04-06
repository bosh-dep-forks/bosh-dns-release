package handlers

import (
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	"github.com/coredns/coredns/plugin/metrics"
	"github.com/miekg/dns"
	"golang.org/x/net/context"
)

type MetricsDNSHandler struct {
	Metrics *metrics.Metrics
	next    dns.Handler
	logger  boshlog.Logger
	logTag  string
}

func NewMetricsDNSHandler(next dns.Handler, logger boshlog.Logger, listenAddress string) MetricsDNSHandler {
	m := metrics.New(listenAddress)
	m.Next = corednsHandlerWrapper{Next: next}
	return MetricsDNSHandler{
		Metrics: m,
		logTag:  "MetricsDNSHandler",
		next:    next,
		logger:  logger,
	}
}

func (m MetricsDNSHandler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	indicator := &requestContext{
		withMetrics: true,
	}
	reqContext := context.WithValue(context.Background(), "indicator", indicator)
	_, err := m.Metrics.ServeDNS(reqContext, w, r)
	if err != nil {
		m.logger.Error(m.logTag, "Error getting dns metrics:", err.Error())
	}
}

func (m MetricsDNSHandler) Run() error {
	return m.Metrics.OnStartup()
}
