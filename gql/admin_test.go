package gql_test

import (
	"testing"

	"github.com/sebastianmontero/dgraph-backup-script/gql"
	"github.com/stretchr/testify/assert"
)

func TestExport(t *testing.T) {
	err := admin.Export(&gql.ExportArgs{
		Destination: config.GetExportURL("test"),
		AccessKey:   config.ExportAccessKey,
		SecretKey:   config.ExportSecretKey,
	})
	assert.NoError(t, err)
}
