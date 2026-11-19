package cmd

import "testing"

// describe-workload read --workload-id without defining it ("unknown flag").
func TestNCSDescribeWorkloadHasWorkloadIDFlag(t *testing.T) {
	if ncsGetWorkloadCmd.Flags().Lookup("workload-id") == nil {
		t.Fatal("ncs describe-workload has no --workload-id flag")
	}
}
