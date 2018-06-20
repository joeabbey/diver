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

***Working with Roles***

Once logged in you can list/get and create roles as per the example below:

```
dan $ ./diver ucp auth roles list | grep jenkins
998612c1-b367-42af-9d82-b2a5de9f8851    false   jenkins
dan $ ./diver ucp auth roles get --rolename jenkins > jenkins.role
dan $ ./diver ucp auth roles create --rolename jenkins2 --ruleset ./jenkins.role
INFO[0000] Role [jenkins2] created succesfully
dan $ ./diver ucp auth roles list | grep jenkins
260976b1-76d7-4ef0-84e2-6ae6b896eed1    false   jenkins2
998612c1-b367-42af-9d82-b2a5de9f8851    false   jenkins
```

**Downloading the client bundle**

Download the client bundle to your local machine.

```
./diver ucp client-bundle
INFO[0000] Downloading the UCP Client Bundle            
```

**Watching Containers**

This will present a colour coded output on memory usage of all containers that are running in a swarm cluster.. (using [urchin](http://github.com/thebsdbox/urchin) to hit memory reservations in the demo below)


```
./diver ucp containers --top
```

![](img/container-top.jpg)

**Debugging Issues**

When errors are reported turn up the `--logLevel` to 5, which enables debugging output.
