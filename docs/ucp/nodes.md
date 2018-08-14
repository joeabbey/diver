## Nodes

The `node` subcommand allows a user to interigate and make various configuration changes to all nodes that are part of a Docker EE cluster


### List all nodes in the EE cluster

The `diver ucp nodes list` command will return all nodes across the entire cluster.

```
diver ucp nodes list
Name             ID                          Role      Version      Platform       Swarm   Kubernetes
docker03.local   d5kqmkg5elaq1ygk6nf7qzmer   worker    18.03.1-ce   linux/x86_64   true    false
docker01.local   l8h2ejtpxkuf5o2loygwk8zun   manager   18.03.1-ce   linux/x86_64   true    true
docker02.local   tlbmntgk7plu19w3ob98r2nel   worker    18.03.1-ce   linux/x86_64   true    false
```

### Manage Orchestrators

By default in Docker EE 2+ there is the option of using multiple orchestrators to manage the Docker EE nodes, the `orchestrator` command provides the option to set which orchestrator will manage the nodes.

**NOTE** it appears that you can also set a node to have no orchestrator, effectively rendering the node unusable by UCP.

The `--swarm` and `--kubernetes` flags will enable either or both orchestrators:

```
diver ucp nodes orchestrator --swarm --id tlbmntgk7plu19w3ob98r2nel
INFO[0000] Configured Node [tlbmntgk7plu19w3ob98r2nel] to allow kubernetes=false and swarm=true
```

De-activating both orchestrators requires passing no flags:

```
diver ucp nodes orchestrator --id tlbmntgk7plu19w3ob98r2nel
WARN[0000] This node has no orchestrators defined and wont be scheduled any workload
INFO[0000] Configured Node [tlbmntgk7plu19w3ob98r2nel] to allow kubernetes=false and swarm=false
```

### Manage Availability

Nodes have three usage states:

- active - in use
- paused - wont take additional tasks
- drain - remove all running tasks

Setting the availability state:

```
diver ucp nodes availability --id d5kqmkg5elaq1ygk6nf7qzmer --state active
INFO[0000] Succesfully set node [d5kqmkg5elaq1ygk6nf7qzmer] to state [active]
```

### Manage node role

Docker nodes can either be managers (run UCP etc.) or workers, which is where most workloads will run

Set the node role:

```
diver ucp nodes role --role worker --id tlbmntgk7plu19w3ob98r2nel
INFO[0000] Succesfully set node [tlbmntgk7plu19w3ob98r2nel] to swarm role [worker]
```

### Apply labels to a node

```
diver ucp nodes label --key labelkey --value labelvalue --id tlbmntgk7plu19w3ob98r2nel
INFO[0000] Succesfully updated node [tlbmntgk7plu19w3ob98r2nel] with the label [labelkey=labelvalue]
```

### Investigate a nodes labels 

```
diver ucp nodes get --id tlbmntgk7plu19w3ob98r2nel
Label Key                                Label Value
labelkey                                 labelvalue
com.docker.ucp.collection                system
com.docker.ucp.collection.root           true
com.docker.ucp.collection.system         true
com.docker.ucp.orchestrator.swarm        false
com.docker.ucp.SANs                      192.168.0.141,localhost,proxy.local,docker02.local,tmpp1fniumpwgvv5c1vusexcw,172.17.0.1,127.0.0.1,10.96.0.1
com.docker.ucp.access.label              /System
com.docker.ucp.collection.swarm          true
com.docker.ucp.orchestrator.kubernetes   false
```
