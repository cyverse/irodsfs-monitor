package service

/*
import (
	"github.com/cyverse/irodsfs/pkg/irodsfs"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

// Reporter reports metrics to the irodsfs monitoring service
type Reporter struct {
	URL      string
	Pusher   *push.Pusher
	FSConfig *irodsfs.Config
}

func NewReporter(fsConfig *irodsfs.Config) Reporter {
	return &Reporter{
		URL:      fsConfig.MonitorURL,
		FSConfig: fsConfig,
	}
}

func (reporter *PrometheusReporter) ReportDataTransferVolume(transferredKB float64) {
}

// ReportDataTransferVolume reports data transfer volume in KB
func (reporter *PrometheusReporter) ReportDataTransferVolume(transferredKB float64) {
	logger := log.WithFields(log.Fields{
		"package":  "report",
		"function": "PrometheusReporter.ReportDataTransferVolume",
	})

	counterOpts := prometheus.CounterOpts{
		Name: "kilobytes_transferred",
		Help: "Kilobytes transferred",
	}
	counter := prometheus.NewCounter(counterOpts)

	counter.Add(transferredKB)
	err := reporter.Pusher.Collector(counter).Grouping("user", reporter.FSConfig.ClientUser).Push()
	if err != nil {
		logger.WithError(err).Error("Unable to push a counter")
		return
	}
}

// ReportDataTransferRequest reports data transfer request count
func (reporter *PrometheusReporter) ReportDataTransferRequest(countRequest int64) {

}
*/
