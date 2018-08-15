## Docker Universal Control Plane Upgrade

The UCP Apis expose an endpoint that allows remotely upgrading.

### Identify the current version and available versions for upgrading

To view all available releases of Docker UCP, run the upgrade command without the `--version` flag: 

```
$ diver ucp upgrade
Upgrade Docker Universal Control Plane

Usage:
  diver ucp upgrade [flags]

Flags:
  -h, --help             help for upgrade
      --version string   The version of UCP to upgrade to

Global Flags:
      --logLevel int   Set the logging level [0=panic, 3=warning, 5=debug] (default 4)

Available Versions
3.0.3
3.0.4

Current Version
ucp/3.0.2

FATA[0000] No --version specified
```

### Upgrade Docker UCP to a newer version

**NOTE** The upgrade process will stop any existing sessions to Docker UCP whilst teh components are updated.

The upgrade process provides a 10 second wait process to provide the enduser the option to stop the upgrade process.

```
diver ucp upgrade --version 3.0.3
INFO[0000] Upgrading from [ucp/3.0.2] to [3.0.3]
WARN[0000] The Universal Control Plane will be un-available for > 5 minutes whilst the upgrade takes place
INFO[0000] The procedure will begin in 10 seconds, press ctrl+c to cancel
INFO[0022] Upgrade procedure has begun succesfully
```

**NOTE** The only method to watch the upgrade process is to access the UCP node that the upgrade started on and find the `ucp-agent:x.x.x` that corresponds to the upgrade version and use `docker logs --follow`.