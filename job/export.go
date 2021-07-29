package job

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sebastianmontero/dgraph-backup-script/gql"
	"github.com/sebastianmontero/slog-go/slog"
)

type Export struct {
	Name         string
	Admin        *gql.Admin
	Args         *gql.ExportArgs
	slog         *slog.Log
	successCount prometheus.Counter
	failedCount  prometheus.Counter
	lastStatus   prometheus.Gauge
}

func NewExport(name, adminEndpoint string, args *gql.ExportArgs, slog *slog.Log) *Export {
	return &Export{
		Name:  name,
		Admin: gql.NewAdmin(adminEndpoint),
		Args:  args,
		slog:  slog,
		successCount: promauto.NewCounter(prometheus.CounterOpts{
			Subsystem: name,
			Name:      "dgraph_export_success_count",
			Help:      "# of succeded exports",
			ConstLabels: prometheus.Labels{
				"dgraph_instance": name,
			},
		}),
		failedCount: promauto.NewCounter(prometheus.CounterOpts{
			Subsystem: name,
			Name:      "dgraph_export_fail_count",
			Help:      "# of failed exports",
			ConstLabels: prometheus.Labels{
				"dgraph_instance": name,
			},
		}),
		lastStatus: promauto.NewGauge(prometheus.GaugeOpts{
			Subsystem: name,
			Name:      "dgraph_last_export_status",
			Help:      "Last export status",
			ConstLabels: prometheus.Labels{
				"dgraph_instance": name,
			},
		}),
	}
}

func (m *Export) Run() {
	m.slog.Infof("Performing export for: %v", m)
	err := m.Admin.Export(m.Args)
	if err != nil {
		m.lastStatus.Set(0)
		m.failedCount.Inc()
		m.slog.Errorf(err, "export run failed for: %v", m)
	}
	m.lastStatus.Set(1)
	m.successCount.Inc()
	m.slog.Infof("export run succeded for: %v", m)
}

func (m *Export) String() string {
	return fmt.Sprintf(
		`
			Export {
				Admin: %v
				Args: %v		
			}
		`,
		m.Admin,
		m.Args,
	)
}
