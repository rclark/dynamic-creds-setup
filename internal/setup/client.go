package setup

import (
	"bytes"
	"context"
	"crypto/sha1"
	"crypto/tls"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
)

var domain string = "app.terraform.io"
var aud string = "aws.workload.identity"

type Client struct {
	iam *iam.Client
}

func NewClient(cfg aws.Config) Client {
	client := Client{
		iam: iam.NewFromConfig(cfg),
	}

	return client
}

func getFingerprint(ctx context.Context, domain string) (string, error) {
	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:443", domain), &tls.Config{})
	if err != nil {
		return "", err
	}

	cert := conn.ConnectionState().PeerCertificates[0]
	fingerprint := sha1.Sum(cert.Raw)
	var buf bytes.Buffer
	for _, f := range fingerprint {
		// if i > 0 {
		// 	fmt.Fprintf(&buf, ":")
		// }
		fmt.Fprintf(&buf, "%02X", f)
	}

	return buf.String(), nil
}

func (c Client) SetupOIDCProvider(ctx context.Context) (string, error) {
	fingerprint, err := getFingerprint(ctx, domain)
	if err != nil {
		return "", err
	}

	oidc, err := c.iam.CreateOpenIDConnectProvider(ctx, &iam.CreateOpenIDConnectProviderInput{
		ThumbprintList: []string{fingerprint},
		Url:            aws.String(fmt.Sprintf("https://%s", domain)),
		ClientIDList:   []string{aud},
	})
	if err != nil {
		return "", err
	}

	return *oidc.OpenIDConnectProviderArn, nil
}

// PolicyDocument defines a policy document as a Go struct that can be serialized
// to JSON. A big "fuck you" goes out the IAM team for never providing this
// struct in their SDKs.
type PolicyDocument struct {
	Version   string
	Statement []PolicyStatement
}

// PolicyStatement defines a statement in a policy document.
type PolicyStatement struct {
	Effect    string
	Action    []string
	Principal map[string]string            `json:",omitempty"`
	Resource  *string                      `json:",omitempty"`
	Condition map[string]map[string]string `json:",omitempty"`
}

func (pd PolicyDocument) String() string {
	if pd.Version == "" {
		pd.Version = "2012-10-17"
	}

	data, _ := json.Marshal(pd)
	return string(data)
}

func (c Client) SetupAdminRole(ctx context.Context, oidcArn, orgName string) (string, error) {
	// allows any project in any workspace from your org to access this IAM role
	subject := fmt.Sprintf("organization:%s:project:*:workspace:*:run_phase:*", orgName)

	assume := PolicyDocument{
		Statement: []PolicyStatement{{
			Effect:    "Allow",
			Principal: map[string]string{"Federated": oidcArn},
			Action:    []string{"sts:AssumeRoleWithWebIdentity"},
			Condition: map[string]map[string]string{
				"StringEquals": {fmt.Sprintf("%s:aud", domain): aud},
				"StringLike":   {fmt.Sprintf("%s:sub", domain): subject},
			},
		}},
	}

	name := "tfc-power-user"

	role, err := c.iam.CreateRole(ctx, &iam.CreateRoleInput{
		AssumeRolePolicyDocument: aws.String(assume.String()),
		RoleName:                 aws.String(name),
		Description:              aws.String("Provides TFC with power-user permissions"),
	})
	if err != nil {
		return "", err
	}

	_, err = c.iam.AttachRolePolicy(ctx, &iam.AttachRolePolicyInput{
		PolicyArn: aws.String("arn:aws:iam::aws:policy/PowerUserAccess"),
		RoleName:  aws.String(name),
	})
	if err != nil {
		return "", err
	}

	return *role.Role.Arn, nil
}
