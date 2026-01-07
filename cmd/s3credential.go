package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/s3credential"
	"github.com/spf13/cobra"
)

var s3credentialCmd = &cobra.Command{
	Use:     "s3-credential",
	Aliases: []string{"s3-cred", "s3cred"},
	Short:   "Manage S3 API credentials",
}

func init() {
	rootCmd.AddCommand(s3credentialCmd)

	s3credentialCmd.AddCommand(s3credListCmd)
	s3credentialCmd.AddCommand(s3credCreateCmd)
	s3credentialCmd.AddCommand(s3credDeleteCmd)

	s3credListCmd.Flags().String("user-id", "", "API user ID (required)")
	s3credListCmd.MarkFlagRequired("user-id")

	s3credCreateCmd.Flags().String("api-user-id", "", "API user ID (required)")
	s3credCreateCmd.Flags().String("tenant-id", "", "Tenant ID for the credential")
	s3credCreateCmd.MarkFlagRequired("api-user-id")

	s3credDeleteCmd.Flags().String("user-id", "", "API user ID (required)")
	s3credDeleteCmd.MarkFlagRequired("user-id")
}

func newS3CredentialClient() *s3credential.Client {
	return s3credential.NewClient(getRegion(), getIdentityCreds(), nil, debug)
}

var s3credListCmd = &cobra.Command{
	Use:   "list",
	Short: "List S3 credentials for a user",
	Run: func(cmd *cobra.Command, args []string) {
		userID, _ := cmd.Flags().GetString("user-id")

		client := newS3CredentialClient()
		result, err := client.ListCredentials(context.Background(), userID)
		if err != nil {
			exitWithError("Failed to list S3 credentials", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ACCESS_KEY\tSECRET_KEY\tUSER_ID\tTENANT_ID\tCREATED_AT")
		for _, cred := range result.Credentials {
			secret := cred.Secret
			if len(secret) > 8 {
				secret = secret[:8] + "..."
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
				cred.Access, secret, cred.UserID, cred.TenantID, cred.CreatedAt.Format("2006-01-02 15:04:05"))
		}
		w.Flush()
	},
}

var s3credCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new S3 credential",
	Run: func(cmd *cobra.Command, args []string) {
		apiUserID, _ := cmd.Flags().GetString("api-user-id")
		credTenantID, _ := cmd.Flags().GetString("tenant-id")

		if credTenantID == "" {
			credTenantID = getTenantID()
		}

		client := newS3CredentialClient()
		input := &s3credential.CreateCredentialInput{
			TenantID: credTenantID,
		}

		result, err := client.CreateCredential(context.Background(), apiUserID, input)
		if err != nil {
			exitWithError("Failed to create S3 credential", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("S3 credential created successfully\n")
		fmt.Printf("Access Key: %s\n", result.Credential.Access)
		fmt.Printf("Secret Key: %s\n", result.Credential.Secret)
		fmt.Printf("User ID:    %s\n", result.Credential.UserID)
		fmt.Printf("Tenant ID:  %s\n", result.Credential.TenantID)
	},
}

var s3credDeleteCmd = &cobra.Command{
	Use:   "delete [access-key]",
	Short: "Delete an S3 credential",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		accessKey := args[0]
		userID, _ := cmd.Flags().GetString("user-id")

		client := newS3CredentialClient()
		if err := client.DeleteCredential(context.Background(), userID, accessKey); err != nil {
			exitWithError("Failed to delete S3 credential", err)
		}

		fmt.Printf("S3 credential %s deleted\n", accessKey)
	},
}
