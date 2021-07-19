package main

import (
	"os"

	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog"
	"github.com/sebastianmontero/dgraph-backup-script/gql"
	"github.com/sebastianmontero/dgraph-backup-script/job"
	"github.com/sebastianmontero/dgraph-backup-script/monitoring"
	"github.com/sebastianmontero/dgraph-backup-script/util"
	"github.com/sebastianmontero/slog-go/slog"
)

var (
	log *slog.Log
)

func main() {
	logConfig := &slog.Config{Pretty: true, Level: zerolog.DebugLevel}
	log = slog.New(logConfig, "dgraph-backup-script")
	if len(os.Args) != 2 {
		log.Panic(nil, "Config file has to be specified as the only cmd argument")
	}
	config, err := util.LoadConfig(os.Args[1])
	if err != nil {
		log.Panicf(err, "Unable to load config file: %v", os.Args[1])
	}

	log.Info(config.String())

	go monitoring.SetupEndpoint(config.PrometheusPort)
	if err != nil {
		log.Panic(err, "Error setting up prometheus endpoint")
	}
	c := cron.New()
	for name, jobConfig := range config.ExportJobs {
		exportJob := job.NewExport(
			name,
			jobConfig.GQLAdminURL,
			&gql.ExportArgs{
				Destination: config.GetExportURL(name),
				AccessKey:   config.ExportAccessKey,
				SecretKey:   config.ExportSecretKey,
			},
			slog.New(logConfig, name),
		)
		log.Infof("Adding export job: %v\n with schedule: %v", exportJob, jobConfig.Schedule)
		_, err := c.AddJob(jobConfig.Schedule, exportJob)
		if err != nil {
			log.Panicf(err, "failed adding export job: %v\n with schedule: %v, error: %v", exportJob, jobConfig.Schedule, err)
		}
	}
	c.Start()

	// Causes the program to wait indefinelty
	select {}
}
