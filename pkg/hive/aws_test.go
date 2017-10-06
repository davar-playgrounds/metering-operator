package hive

import (
	"fmt"
	"testing"
	"time"

	"github.com/coreos-inc/kube-chargeback/pkg/aws"
	cb "github.com/coreos-inc/kube-chargeback/pkg/chargeback/v1"
)

func TestCreateAWSUsageTable(t *testing.T) {
	hiveHost, s3Bucket, _ := setupHiveTest(t)

	conn, err := Connect(hiveHost)
	if err != nil {
		t.Fatal("error connecting: ", err)
	}
	defer conn.Close()

	dropQuery := fmt.Sprint("DROP TABLE ", AWSUsageTableName)
	if err = conn.Query(dropQuery); err != nil {
		t.Errorf("Could not delete existing table: %v", err)
	}

	manifests := getAWSManifests(t)
	if err = CreateAWSUsageTable(conn, AWSUsageTableName, s3Bucket, manifests[0]); err != nil {
		t.Error("error perfoming query: ", err)
	}
}

func getAWSManifests(t *testing.T) []aws.Manifest {
	bucket := "coreos-team-chargeback"
	prefix := "coreos-detailed-billing/coreosinc//coreos-detailed-billing-001/team-chargeback-testing"
	begin := time.Date(2017, time.July, 2, 0, 0, 0, 0, time.UTC)
	end := time.Date(2017, time.July, 29, 0, 0, 0, 0, time.UTC)
	rng := cb.Range{begin, end}
	manifests, err := aws.RetrieveManifests(bucket, prefix, rng)
	if err != nil {
		t.Error("unexpected error: ", err)
	}
	return manifests
}