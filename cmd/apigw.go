package cmd

import (
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/apigw"
	"github.com/spf13/cobra"
)

var apigwCmd = &cobra.Command{
	Use:     "apigw",
	Aliases: []string{"api-gateway", "apigateway"},
	Short:   "Manage API Gateway services, stages, deployments, and API keys",
}

func init() {
	rootCmd.AddCommand(apigwCmd)
}

func newAPIGWClient() *apigw.Client {
	appKey := getAppKey()
	accessKey := getAccessKey()
	secretKey := getSecretKey()
	return apigw.NewClient(getRegion(), appKey, accessKey, secretKey, nil, debug)
}
