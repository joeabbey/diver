## Docker Enterprise Edition Servics

### List all services

```
$ diver ucp services list urchin
Name              ID
urchin            p59zsqz20308wr9p8y9sp7ieg
ucp-agent-s390x   tuns8xv029cgr6nmoqijus62u
ucp-agent         tzafamk01pqszot3hlquci1ev
ucp-agent-win     z925y17hbkuky6tf5lffh8k6k
```

### Retrieve design of a service

```
$ diver ucp services architecture urchin
INFO[0000] Inspecting service [urchin]                  
ID:        p59zsqz20308wr9p8y9sp7ieg
Version:   19222
Name:      urchin
Image:     thebsdbox/urchin:1.2@sha256:fbadb7d721cd9faabdead81323a02deb1a05993e3e60c0762eb249bed2d168d3
Cmd: /urchin
Args: -w 8080
Labels:
                      com.docker.ucp.access.label       /
                      com.docker.ucp.collection         swarm
                      com.docker.ucp.collection.root    true
                      com.docker.ucp.collection.swarm   true
Memory Reservation:   0
CPU Reservation:      0
Memory Limits:        0
CPU Limits:           0
Replicas:             1
```

## Deep dive into services

### Retrieve service health

```
$ diver ucp services get health --name urchin
Service   Expected   Running   Shutdown   Failed   Tasks   Status
urchin    1          1         19         3        24
```

### Retrieve information about service tasks

```
$ diver ucp services get tasks --name p59zsqz20308wr9p8y9sp7ieg
Hostame                                ID                                                                 Node       IP Address		 Network   State
/urchin.4.17bu1l1r75qf9d0o1q771o94l    3a0f01b2ebd11da6c5a2815ce3d93832b1dd3e9fd32c8e911d6dff09f5e78032   docker03   10.255.6.211/16  ingress  shutdown                         
/urchin.7.29hl9exork2bsl5azv68z9hbl    7ef1a44d80bbe61db80793068c3ff5889018efefbb266dcd6fcb7bf93a50cb4f   docker01   10.255.6.206/16  ingress  shutdown                         
/urchin.19.b4iu2sc158x86gh7k856gqcf8   cc53607a5cd7c7d0d0ba2f195384c77079f5d046490c4090257fc0e3f87b24ea   docker02   10.255.6.218/16  ingress  shutdown                         
/urchin.6.fy4luzc2lt8d17xk5gzhplcal    d3f5b7427b58e608fc7ca4f437cb16732b5850a83f651c08175ef03e7faeafaf   docker03   10.255.6.213/16  ingress  shutdown                         
/urchin.20.ibe8zp7zlv5wezrhc3v3t31dm   e516aa877b0c95845d1432cd8d3cc1bb5ddfc0232966b4875693a6c0c88158ae   docker03   10.255.6.219/16  ingress  shutdown                         
/urchin.3.khvvsgrfg9ctmwrsedt11249h    4a576e2cf8ce3de83395d7e43617e13bdc33212948c826c20d44511a2a67f033   docker01   10.255.6.205/16  ingress  shutdown                         
/urchin.5.m57efzr4l6lhgyrwwfv67zrg3    3fc0d0f673c1f2f2914e34e1e9e0d89fd711c50287a9ce2c503bfed1369e0bbc   docker01   10.255.6.212/16  ingress  shutdown                         
/urchin.9.n72yt7jo03w36m2vzuhmnbpsl    201e32b4b8e3891ca83e264ddf047950c5c38d692cefd34e3bff29db55cbb7d8   docker01   10.255.6.215/16  ingress  shutdown                         
/urchin.10.oi0z604nt1asl5pemr7hsmg1q   22390b5a407daf5b1f5da78db6279863f108104c6d5b0af8fd042f9b2be2dcb2   docker02   10.255.6.216/16  ingress  shutdown                         
/urchin.8.oqgkghga9noihoug8f3fw0lsp    94a859261ef199a687f9be4489f44fdb38ddbcef4f24f83c5f9e086288eedba9   docker02   10.255.6.214/16  ingress  shutdown                         
/urchin.1.qdu50tqmigqno79ntlq8e6ssd    e0d7302ace05b21b2496615736994cd5912dfdb2d7cb57dcffe2cf98c1c8b1d5   docker01   10.255.6.240/16  ingress  10.0.2.26/24  urchin  complete   
/urchin.1.r7it6s218brkbnz427kruhkun    4a49c61d8ecdcc27cbd9b09fa4d0729c64b7c5c17052c2e5f647bbb00a05ea3c   docker01   10.255.6.241/16  ingress  10.0.2.27/24  urchin  running    
```
