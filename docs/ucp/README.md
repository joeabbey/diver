### Logging in to UCP

In order to work with the Docker Universal Control Plane the `login` command is used:

```
./diver ucp login --username docker               \
                  --password password             \
                  --url https://docker01.fnnrn.me \
                  --ignorecert
INFO[0000] Succesfully logged into [https://docker01.fnnrn.me] 
```

**Note** the `login` command will create a `~/.ucptoken` file that is used for all further `diver` commands. In the event that login errors start to occur, or logins fail check the permissions of this file or alternatively remove this file.

### Working with multiple UCP sessions

When you log into a new UCP server, then a new session will be created. To view the sessions that are available use the following command:

```
$ ./diver ucp login list
ID   Address                     Active
0    https://192.168.0.140       true
1    https://docker01.fnnrn.me   false
```

To set a different session to active, you can use the setActive command to change the active session.

```
./diver ucp login setActive --id 1
INFO[0000] Set session [1] to active
./diver ucp login list
ID   Address                     Active
0    https://192.168.0.140       false
1    https://docker01.fnnrn.me   true
```


### Checking access

The `diver ucp` command will return all subcommands but also confirm that your login works and returns the username of the user that is logged in:

```
$ diver ucp
Universal Control Plane

Usage:
  diver ucp [flags]
  diver ucp [command]

Available Commands:

{...}

INFO[0000] Current user [docker]                        
$
```
