# Diver

This is a tool to interact with the APIs of the Docker Enterprise Edition products enabling an end user to provision, manage and monitor the platform. 

### Building

Once the repo is cloned build with `make build`, or `make docker` to create a container locally that contains the binary.

### Usage

The main commands are `dtr` and `ucp` that interact directly with those areas of the EE platform.

```
./diver -h
This tool uses the native APIs to "dive" into Docker EE

Usage:
  diver [command]

Available Commands:
  dtr         Docker Trusted Registry
  help        Help about any command
  ucp         Universal Control Plane 

Flags:
  -h, --help   help for diver

Use "diver [command] --help" for more information about a command.
```

**Logging in to UCP**

```
./diver ucp --username docker               \
            --password password             \
            --url https://docker01.fnnrn.me \
            --ignorecert
INFO[0000] Succesfully logged into [https://docker01.fnnrn.me] 
```

**Creating users/organisations**

This uses the `auth` command as part of `ucp`.

This will create a new **user** called `bob`, to create an organisation use the `-isorg` flag. The `--action` flag identifies what operation will take place, such as `create`, `delete` and `modify`.

```
./diver ucp auth --active               \
                 --admin                \
                 --fullname "Bob Smith" \
                 --username bob         \
                 --password chess123    \
                 --action create
```

**Watching Containers**

This will present a colour coded output on memory usage of all containers that are running in a swarm cluster.. (using [urchin](http://github.com/thebsdbox/urchin) to hit memory reservations in the demo below)

![](img/container-top.jpg)

**Debugging Issues**

When errors are reported turn up the `--logLevel` to 5, which enables debugging output.
