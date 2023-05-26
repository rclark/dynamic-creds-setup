# dynamic-creds-setup

A CLI command to bootstrap an AWS account for Terraform Cloud Dynamic Provider Credentials.

1. Ensure that you have valid AWS credentials setup for the AWS account you wish to bootstrap in your shell session.
1. Run the command. Provide the name of your TFC org as a single argument.
1. Follow [the instructions here](https://developer.hashicorp.com/terraform/cloud-docs/workspaces/dynamic-provider-credentials/aws-configuration#configure-terraform-cloud) to setup your TFC workspace.
