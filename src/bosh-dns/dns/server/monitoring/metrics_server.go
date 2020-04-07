package monitoring

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	"github.com/coredns/coredns/plugin/metrics"
)

type MetricsServer struct {
	Metrics *metrics.Metrics
	logger  boshlog.Logger
	logTag  string
}

func NewMetricsServer(logger boshlog.Logger, listenAddress string) *MetricsServer {
	m := metrics.New(listenAddress)
	return &MetricsServer{
		Metrics: m,
		logTag:  "MetricsServer",
		logger:  logger,
	}
}

func (m *MetricsServer) startup() error {
	err := m.Metrics.OnStartup()
	if err != nil {
		return bosherr.WrapError(err, "setting up metrics on startup")
	}

	return nil
}

func (m *MetricsServer) shutdown() error {
	err := m.Metrics.OnFinalShutdown()
	if err != nil {
		return bosherr.WrapError(err, "tear down and restart of the metrics listener")
	}

	return nil
}

func (m *MetricsServer) Run(shutdown chan struct{}) {
	if err := m.startup(); err != nil {
		m.logger.Error(m.logTag, "running: %s", err)
	}
	for {
		select {
		case <-shutdown:
			if err := m.shutdown(); err != nil {
				m.logger.Error(m.logTag, "running: %s", err)
			}
			return
		}
	}
}
