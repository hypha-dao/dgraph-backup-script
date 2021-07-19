package gql_test

import (
	"log"
	"os"
	"testing"

	"github.com/sebastianmontero/dgraph-backup-script/gql"
	"github.com/sebastianmontero/dgraph-backup-script/util"
)

var admin *gql.Admin
var config *util.Config

func TestMain(m *testing.M) {
	beforeAll()
	// exec test and this returns an exit code to pass to os
	retCode := m.Run()
	afterAll()
	// If exit code is distinct of zero,
	// the test will be failed (red)
	os.Exit(retCode)
}

func beforeAll() {
	var err error
	config, err = util.LoadConfig("./config.yml")
	if err != nil {
		log.Fatal(err, "Failed to load configuration")
	}
	admin = gql.NewAdmin(config.ExportJobs["dho"].GQLAdminURL)
}

func afterAll() {
}
