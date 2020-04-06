package handlers_test

import (
	"bosh-dns/dns/server/handlers"
	"bosh-dns/dns/server/handlers/handlersfakes"
	"bosh-dns/dns/server/internal/internalfakes"
	"net"

	"github.com/cloudfoundry/bosh-utils/logger/loggerfakes"
	"github.com/coredns/coredns/plugin/metrics"
	"github.com/miekg/dns"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ bool = Describe("metricsHandler", func() {
	var (
		metricsHandler handlers.MetricsDNSHandler
		fakeWriter     *internalfakes.FakeResponseWriter
		fakeDnsHandler *handlersfakes.FakeDNSHandler
		fakeLogger     *loggerfakes.FakeLogger
		response       *dns.Msg
	)

	BeforeEach(func() {
		fakeDnsHandler = &handlersfakes.FakeDNSHandler{}
		fakeWriter = &internalfakes.FakeResponseWriter{}
		fakeLogger = &loggerfakes.FakeLogger{}
		metricsHandler = handlers.NewMetricsDNSHandler(fakeDnsHandler, fakeLogger, "127.0.0.1:53088")

		response = &dns.Msg{
			Answer: []dns.RR{&dns.A{A: net.ParseIP("99.99.99.99")}},
		}
		response.SetQuestion("my-instance.my-group.my-network.my-deployment.bosh.", dns.TypeANY)
		fakeDnsHandler.ServeDNSStub = func(responseWriter dns.ResponseWriter, r *dns.Msg) {
			response.SetRcode(r, dns.RcodeSuccess)
			responseWriter.WriteMsg(response)
		}
	})

	Describe("ServeDNS", func() {
		It("collects metrics", func() {
			m := &dns.Msg{}
			m.SetQuestion("my-instance.my-group.my-network.my-deployment.bosh.", dns.TypeANY)
			metricsHandler.ServeDNS(fakeWriter, m)

			Expect(findMetric(metricsHandler.Metrics, "coredns_dns_request_count_total")).To(Equal(1.0))
		})
	})
})

func findMetric(m *metrics.Metrics, key string) float64 {
	metricFamilies, _ := m.Reg.Gather()
	for _, mf := range metricFamilies {
		if mf.GetName() == key {
			return *mf.GetMetric()[0].Counter.Value
		}
	}
	return -1.0
}
