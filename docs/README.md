# Diver

The main commands are `dtr`, `ucp` and `store` that interact directly with those areas of the EE platform or the Docker store.

```
$ diver
This tool uses the native APIs to "dive" into Docker EE

Usage:
  diver [command]

Available Commands:
  dtr         Docker Trusted Registry
  help        Help about any command
  store       Docker Store
  ucp         Universal Control Plane 
  version     Version and Release information about the diver tool

Flags:
  -h, --help           help for diver
      --logLevel int   Set the logging level [0=panic, 3=warning, 5=debug] (default 4)

```

## *STATUS*

**[UCP](./ucp)**

- Login
- Query/Create/Delete Users
- Query/Create/Delete Organisations
- Query/Create/Delete Teams
- Get/Set Swarm
- Clone and set Roles
- Set Grants from subject, role, object
- Build Images from a local or Remote (github URL) Dockerfile
- Download client bundle
- Inspect Services (Endpoints coming soon)
- Manage Swarm configuration
- Attach/Detach containers from networks

.. More coming

**[DTR](./dtr)**

- Login
- List Replicas and health
- List Repositories
- Create Repositories
- Create/List/Delete WebHooks
- Get/Set DTR Settings

.. More coming

**[STORE](./store)**

- Login
- Retrieve Subscriptions
- Find recent active
- Download licenses
- Create and retrieve trial and prod sub urls

## Debugging Issues

When errors are reported turn up the `--logLevel` to 5, which enables debugging output.
