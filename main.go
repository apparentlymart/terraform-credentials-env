// terraform-credentials-env is a Terraform credentials helper that reads
// credentials from the process environment.
//
// Specifically, it expects to find environment variables with the prefix
// TF_TOKEN_ followed by the requested hostname, such as
// TF_TOKEN_app.terraform.io for Terraform Cloud.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	svchost "github.com/hashicorp/terraform-svchost"
)

var GitCommit = ""
var Version = "0.0.0"
var PreRelease = "dev"

func main() {
	creds := collectCredentialsFromEnv()

	args := os.Args[1:]
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: terraform-credentials-env get <hostname>")
		fmt.Fprintln(os.Stderr, "\nThis is a Terraform credentials helper, not intended to be run directly from a shell.")
		os.Exit(1)
	}

	switch args[0] {
	case "get":
		// The credentials helper protocol calls for Terraform to provide the
		// hostname already in the "for comparison" form, so we'll assume that
		// here and let this not match if the caller isn't behaving.
		wantedHost := svchost.Hostname(args[1])
		token, ok := creds[wantedHost]
		if !ok {
			fmt.Fprintf(os.Stderr, "No credentials for %s are defined via environment variables.\n", svchost.ForDisplay(args[1]))
			os.Exit(1)
		}
		result := resultJSON{token}
		resultJSON, err := json.Marshal(result)
		if err != nil {
			// Should never happen
			fmt.Fprintf(os.Stderr, "Failed to serialize result: %s\n", err)
			os.Exit(1)
		}
		os.Stdout.Write(resultJSON)
		os.Stdout.WriteString("\n")
		os.Exit(0)

	default:
		fmt.Fprintf(os.Stderr, "The 'env' credentials helper is not able to %s credentials.\n", args[0])
		os.Exit(1)
	}
}

func collectCredentialsFromEnv() map[svchost.Hostname]string {
	const prefix = "TF_TOKEN_"

	ret := make(map[svchost.Hostname]string)
	for _, ev := range os.Environ() {
		eqIdx := strings.Index(ev, "=")
		if eqIdx < 0 {
			continue
		}
		name := ev[:eqIdx]
		value := ev[eqIdx+1:]
		if !strings.HasPrefix(name, prefix) {
			continue
		}
		rawHost := name[len(prefix):]
		// We accept underscores in place of dots because dots are not valid
		// identifiers in most shells and are therefore hard to set.
		// Underscores are not valid in hostnames, so this is unambiguous for
		// valid hostnames.
		rawHost = strings.ReplaceAll(rawHost, "_", ".")

		// Because environment variables are often set indirectly by OS
		// libraries that might interfere with how they are encoded, we'll
		// be tolerant of them being given either directly as UTF-8 IDNs
		// or in Punycode form, normalizing to Punycode form here because
		// that is what the Terraform credentials helper protocol will
		// use in its requests.
		//
		// Using ForDisplay first here makes this more liberal than Terraform
		// itself would usually be in that it will tolerate pre-punycoded
		// hostnames that Terraform normally rejects in other contexts in order
		// to ensure stored hostnames are human-readable.
		dispHost := svchost.ForDisplay(rawHost)
		hostname, err := svchost.ForComparison(dispHost)
		if err != nil {
			// Ignore invalid hostnames
			continue
		}

		ret[hostname] = value
	}

	return ret
}

type resultJSON struct {
	Token string `json:"token"`
}
