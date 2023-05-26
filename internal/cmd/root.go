package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/spf13/cobra"

	"github.com/rclark/dynamic/internal/setup"
)

var cfg aws.Config

func Execute() {
	var err error
	cfg, err = config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

var cmd = &cobra.Command{
	Use:   "dynamic-creds-setup [TFC organization name]",
	Short: "Bootstrap an AWS account for TFC Dynamic Provider Credentials",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		if len(args) != 1 {
			return errors.New("Must provide 1 argument: the name of a TFC organization")
		}
		client := setup.NewClient(cfg)
		oidcARN, err := client.SetupOIDCProvider(ctx)
		if err != nil {
			return err
		}
		roleARN, err := client.SetupAdminRole(ctx, oidcARN, args[0])
		if err != nil {
			return err
		}
		fmt.Println("Success!")
		w := tabwriter.NewWriter(os.Stdout, 0, 4, 1, ' ', 0)
		fmt.Fprintf(w, "  OIDC Provider ARN:\t%s\n", oidcARN)
		fmt.Fprintf(w, "  Role ARN:\t%s\n", roleARN)
		w.Flush()
		fmt.Println("\nYou may want to record these values in case you ever wish to import the resources into Terraform")
		return nil
	},
}
