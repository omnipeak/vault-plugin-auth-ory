# Ory Kratos + Keto Auth Plugin for HashiCorp Vault

This repository contains code for a [HashiCorp Vault](https://github.com/hashicorp/vault) Auth [Plugin](https://developer.hashicorp.com/vault/docs/plugins) that authenticates with [Ory Kratos](https://github.com/ory/kratos) and [Ory Keto](https://github.com/ory/keto) APIs.

## Setup

The setup guide assumes some familiarity with Vault and Vault's plugin
ecosystem. You must have a Vault server already running, unsealed, and
authenticated.

1. Download and decompress the latest plugin binary from the Releases tab. Alternatively you can compile the plugin from source.

2. Move the compiled plugin into Vault's configured `plugin_directory`:

  ```sh
  $ mv vault-auth-plugin-ory /etc/vault/plugins/vault-auth-plugin-ory
  ```

3. Calculate the SHA256 of the plugin and register it in Vault's plugin catalog.
If you are downloading the pre-compiled binary, it is highly recommended that
you use the published checksums to verify integrity.

  ```sh
  $ export SHA256=$(shasum -a 256 "/etc/vault/plugins/vault-auth-plugin-ory" | cut -d' ' -f1)

  $ vault plugin register \
      -sha256="${SHA256}" \
      -command="vault-auth-plugin-ory" \
      auth vault-plugin-auth-ory
  ```

4. Mount the auth method:

  ```sh
  $ vault auth enable \
      -path="ory" \
      -plugin-name="vault-plugin-auth-ory" plugin
  ```

## Development Setup

1. Build the plugin for your platform (`os/arch`) e.g.:

  ```sh
  $ make darwin/arm64
  ```

  or build for all platforms:

  ```sh
  $ make build
  ```

2. Start a Vault server in dev mode pointing to the plugin directory:

  ```sh
  $ make start
  ```
3. Login to Vault as root:

  ```sh
  $ vault login root
  ```

4. Enable the plugin in Vault:

  ```sh
  $ make enable
  ```
5. Write the configs:

  ```sh
  $ make configs
  ```

6. Authenticate with the plugin:

  ```sh
  $ vault write auth/ory/login \
namespace=files \
object=c5cc3e28-e3c3-45ca-be86-a0a55953bfca \
relation=editor \
kratos_session_cookie=ory_kratos_session=MTY2NzgyMjg2M3xBYVJxa2hmNFlOOFAyZnc3U3VidnZKd1A0VmdyWFgyU3ozbUNvRG4zeC1oNU1DS3Z6dkc1ODllTHdua0s5aFdpcW1ZZ0pveVNBVVM3ZXBIRWdQdlJGWXN0aS1iVU5tenVFbUw1WE1QNDRVcms5eWZZRk52R3dOdTJKLVcxYVlFWFU4ajNFUmc0bnc9PXyq29KzMQjNDdZLeJAuNLUBeU1g1-iD7l31nahltn4mZg==
  ```

7. Add a policy that matches the naming convention `namespace_relation` (e.g. `files_editor`) using the example policy found below, replacing the accessor string with the contents returned by:

  ```sh
  $ make accessor
  ```

8. Login with the token provided after running the `write` command:

  ```sh
  $ vault login [token]
  ```

9. Attempt to read secrets:

  ```sh
  $ vault read secret/files/c5cc3e28-e3c3-45ca-be86-a0a55953bfca/some_secret
  ```

## Authenticating with Ory Kratos and Keto

To authenticate, the user supplies a valid Ory Kratos session cookie, along with the namespace,
object, and relation to check against Keto.

```sh
$ vault write auth/ory/login namespace=[namespace] object=[object] relation=[relation] kratos_session_cookie=[full kratos_session_cookie=[...] string]
```

The response will be a standard auth response with some token metadata:

```text
Key                     Value
--------------------------------
token                   [token]
token_accessor          [accessor]
token_duration          [TTL]
token_renewable         false
token_policies          ["default" "[namespace]_[relation]"]
identity_policies       []
policies                ["default" "[namespace]_[relation]"]
```

## Policy Template

When a token is successfully created, the plugin attach a policy that follows the naming schema of `[namespace]_[relation]`.

You must then create a policy with that name in Vault that utilises the metadata stored in the alias. The following policy template will allow access to a KV secret at the path `secret/data/[namespace]/[object]*`:

```hcl
path "secret/data/{{identity.entity.aliases.auth_vault-plugin-auth-ory_e40b77a0.metadata.object}}*" {
  capabilities = ["create", "read", "update", "delete", "list"]
}

path "secret/metadata/{{identity.entity.aliases.auth_vault-plugin-auth-ory_e40b77a0.metadata.object}}/*" {
  capabilities = ["delete", "list"]
}
```

Being sure to replace `auth_vault-plugin-auth-ory_e40b77a0` with the accessor of the auth plugin found by running `vault auth list`.

Alternatively, you can find the accessor by running the following command:

```sh
$ export MOUNT_ACCESSOR=$(vault auth list -format=json | jq -r '."ory/".accessor')
```

As we already know the namespace at this point, you can also simply use the path:

`secret/data/[known namespace]/{{identity.entity.metadata.object}}*`

## License

This code is licensed under the MPLv2 license.
