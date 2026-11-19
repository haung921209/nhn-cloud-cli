package cmd

import "testing"

// MariaDB create-db-snapshot, enable-multi-az and create-read-replica take the
// fields the v4.0 API requires, like their MySQL counterparts.
func TestMariaDBCommandsHaveRequiredRequestFlags(t *testing.T) {
	if f := createMariaDBSnapshotCmd.Flags().Lookup("backup-method-type"); f == nil || f.DefValue != "FULL" {
		t.Errorf("create-db-snapshot --backup-method-type: %v", f)
	}
	if enableMultiAZMariaDBCmd.Flags().Lookup("db-instance-candidate-name") == nil {
		t.Error("enable-multi-az has no --db-instance-candidate-name")
	}
	if createMariaDBReadReplicaCmd.Flags().Lookup("availability-zone") == nil {
		t.Error("create-read-replica has no --availability-zone")
	}
}
