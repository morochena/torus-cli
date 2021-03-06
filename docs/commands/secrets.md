# Secrets
The Torus CLI is used to set and access secrets and config. Each value is encrypted client-side according to our [cryptography](../internals/crypto.md) architecture and stored in the Torus Registry.

A secret is a single piece of configuration which should be encrypted.

Torus exposes your decrypted secrets to your process through environment variables. This means that anything you can store in an environment variable, you can set in Torus. Support for file storage (such as certificates) is planned, but not currently supported.

### Command Options

  Option | Description
  ---- | ----
  --org, ORG, -o ORG | Executing the command for the specified org.
  --project PROJECT, -p PROJECT | Executing the command for the specified project.
  --environment ENV, -e ENV | Executing the command for the specified environment.
  --service SERVICE, -s SERVICE | Execute the command for the specified environment. (default: default)
  --user USER, -u USER | Execute the command for the specified user identity. (default: *)
  --machine MACHINE, -m MACHINE | Execute the command for the specified machine identity. (default: *)
  --instance INSTANCE, -i INSTANCE | Execute the command for the specified instance identity. (default: *)

The `--environment`, `--service`, `--user`, `--machine`, and `--instance` flags can be specified many times when setting, unsetting, or importing secrets.

## set
###### Added [v0.1.0](https://github.com/manifoldco/torus-cli/blob/master/CHANGELOG.md)

`torus set <name|path> <value>` or `torus set <name|path>=<value>` sets a value
for the specified name (or [path](../concepts/path.md)).

This is how all secrets are stored in Torus.

#### Examples

**Using flags**

```bash
# Setting the port for the production auth service inside myorg's api project.
$ torus set -o myorg -p api -e production -s auth PORT 3000

Credentials retrieved
Keypairs retrieved
Encrypting key retrieved
Credential encrypted
Completed Operation

Credential PORT has been set at /myorg/api/production/auth/*/*/PORT
```

**Using flags inside a linked directory**

You can link a directory to an org and project using the [`torus link`](./project-structure.md) command. Once a directory is linked, you only need to specify the flags you want to override.

```bash
# Setting the port for the auth service in staging
$ torus set -e staging -s auth PORT 8080

Credentials retrieved
Keypairs retrieved
Encrypting key retrieved
Credential encrypted
Completed Operation

Credential PORT has been set at /myorg/api/staging/auth/*/*/PORT
```

**Using multiple flags to share a secret between environments**

You can specify the environment, service, user, machine, and instance flags many times to share a secret.

```bash
# Setting the same port for production and staging for the auth service
$ torus set -e production -e staging -s auth PORT 8080

Credentials retrieved
Keypairs retrieved
Encrypting key retrieved
Credential encrypted
Completed Operation

Credential PORT has been set at /myorg/api/[production|staging]/auth/*/*/PORT
```

**Setting a secret with a `*` value**

You can set a secret to be shared across all environments, or services by specifying a value of `*`. For example, if you set an environment to be `*` then any environment (production, staging, dev, etc) will have access to the value.

```bash
# Setting the same port across all environments for the auth service
$ torus set -e * -s auth PORT 8080

Credentials retrieved
Keypairs retrieved
Encrypting key retrieved
Credential encrypted
Completed Operation

Credential PORT has been set at /myorg/api/[production|staging]/auth/*/*/PORT
```

## unset
###### Added [v0.1.0](https://github.com/manifoldco/torus-cli/blob/master/CHANGELOG.md)

`torus unset <name|path>` unsets the value for the specified name (or [path](../concepts/path.md)).

## import
###### Added [v0.25.0](https://github.com/manifoldco/torus-cli/blob/v0.25.0/CHANGELOG.md)

`torus import <file>` or using stdin redirection (e.g. `torus import -e production <prod.env`) imports the contents of an `.env` file to the specified path.

**Example**

```bash
$ cat prod.env
PORT=4000
DOMAIN=mydomain.co
MYSQL_URL=mysql://user:pass@host.com:4321/mydb

$ torus import -e production test.env

Credentials retrieved
Keypairs retrieved
Encrypting key retrieved
Credential encrypted
Credential encrypted
Completed Operation

Credential port has been set at /myorg/myproject/production/default/*/*/PORT
Credential mysql_url has been set at /myorg/myproject/production/default/*/*/MYSQL_URL
```

## view
###### Added [v0.1.0](https://github.com/manifoldco/torus-cli/blob/master/CHANGELOG.md)

`torus view` displays secrets in the current [context](./project-structure.md#link).

By default items are displayed in environment variable format.

### Command Options

  Option | Description
  ---- | ----
  --verbose, -v | List the sources of the secrets (shortcut for --format verbose)
  --format FORMAT, -f FORMAT | Format used to display data (json, env, verbose) (default: env)

## run
###### Added [v0.1.0](https://github.com/manifoldco/torus-cli/blob/master/CHANGELOG.md)

`torus run [--] <command> [<arguments>...]` executes the supplied command and injects your secrets into its environment.

By prefixing your process execution with `torus run` we are able to fetch, decrypt and inject your secrets into the process environment based on the [context](./project-structure.md#link) of the Torus client.

To ensure that your command’s arguments and options are passed correctly you may need to separate your `torus run` definition from your command definition with `--`, for example:

```
# Inject secrets into the process from the production environment and www service.
$ torus run -e production -s www -- node ./bin/www --app api
```

## ls
###### Added [v0.13.0](https://github.com/manifoldco/torus-cli/blob/master/CHANGELOG.md)

`torus ls [path]` allows a user to explore the objects and values set within their organization.

Each level within the organization can be inspected by changing the segments supplied in the path. Wildcards cannot be used for the organization or project segments of the [path](../concepts/path.md).

Path is required, and does not support context.

### Command Options

  Option | Description
  ---- | ----
  --verbose, -v | Show which type of path is being displayed, shortcut for --format=verbose
  --format FORMAT, -f FORMAT | Format used to display data (simple, verbose) (default: simple)

### Examples

List all secrets in a project:
```
$ torus ls /my-org/landing-page/**
/my-org/landing-page/dev-*/[api|www]/*/*/port
/my-org/landing-page/[dev-jeff|dev-sally]/api/*/*/token
```

List orgs you are a member of:
```
$ torus ls / -v
ORGS
/my-org
/another-org
```

List all projects within an org:
```
$ torus ls /my-org/
/my-org/landing-page
/my-org/other-project
```

List all environments which match the path:
```
$ torus ls /my-org/landing-page/dev*
/my-org/landing-page/dev-jeff
/my-org/landing-page/dev-sally
```

Expanding the path to view secrets using recursive:
```
$ torus ls /my-org/landing-page/dev-jeff -r
/my-org/landing-page/dev-*/[api|www]/*/*/port
/my-org/landing-page/[dev-jeff|dev-sally]/www/*/*/token
```
