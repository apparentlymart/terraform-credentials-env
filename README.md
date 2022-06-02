# Terraform Credentials from the Environment

---

**NOTE:** Largely-equivalent functionality to this credentials helper is built in to Terraform v1.2 and later as a standard feature. Therefore this credentials helper is useful only for earlier versions of Terraform that didn't have that built in.

I recommend using Terraform v1.2 or later instead of adopting this credentials helper. I have no plans to develop this credentials helper any further or to produce any new releases of it.

---

`terraform-credentials-env` is a Terraform "credentials helper" plugin that
allows providing credentials for
[Terraform-native services](https://www.terraform.io/docs/internals/remote-service-discovery.html)
(private module registries, Terraform Cloud, etc) via environment variables.

To use it,
[download a release archive](https://github.com/apparentlymart/terraform-credentials-env/releases)
and extract it into the `~/.terraform.d/plugins` directory where Terraform
looks for credentials helper plugins. (The filename of the file inside the
archive is important for Terraform to discover it correctly, so don't rename
it.)

Terraform will take the newest version of the plugin it finds in the plugin
search directory, so if you are switching between versions you may prefer to
remove existing installed versions in order to ensure Terraform selects the
desired version.

Once you've installed the plugin, enable it by adding the following block
to your
[Terraform CLI configuration](https://www.terraform.io/docs/commands/cli-config.html):

```hcl
credentials_helper "env" {}
```

This credentials helper plugin does not take any additional arguments, so the
block must be left empty as shown above.

With this helper installed and enabled, you can set credentials for specific
hostnames in the environment for your shell so that they will be inherited
by `terraform` and then in turn by `terraform-credentials-env`.

The environment variables must be named `TF_TOKEN_` followed by the hostname
the token is for with periods replaced by underscores. For example, to set
a token for `app.terraform.io` (Terraform Cloud) in bash:

```
export TF_TOKEN_app_terraform_io=example_token
```

Terraform will execute the configured credentials helper plugin whenever it
needs to make a request to a Terraform-native service whose credentials aren't
directly configured in the CLI configuration using `credentials` blocks.
`credentials` blocks override credentials helpers though, so if you have any
existing `credentials` block for the hostname you wish to configure you will
need to remove that block first.
