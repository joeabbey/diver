## DTR

### Logging into DTR

```
./diver dtr login --username docker               \
                  --password password             \
                  --url https://docker02.fnnrn.me \
                  --ignorecert
INFO[0000] Succesfully logged into [https://docker02.fnnrn.me] 
```


### Repositories

To interact with Repositories the `repos` command is available:

**Create**

```
diver dtr repos create --name newrepo --namespace allsafe
INFO[0003] New Repository [newrepo] created for namespace [allsafe] 
```

**List**

```
diver dtr repos list
Repo              ID                                     Name      Namespace   Description
allsafe/newrepo   5790b01c-c74c-44c2-bfdc-192df23121ea   newrepo   allsafe     
ecorp/nginx       ff621b07-a595-4d36-958e-58f7bb161d00   nginx     ecorp       
ecorp/pause       64bdbaf7-69a4-491b-b577-2183c1e5f3f4   pause     ecorp       
ecorp/test        9cca0a1d-b8e9-4a51-bd14-aa6f8233b7ce   test      ecorp    
```

### List Replicas

```
./diver dtr info replicas

Replica         Status
a3a8ab213a8b     OK
ecb7a768afc4     OK
```

### Settings

```
diver dtr settings get
Configuration           Setting
DTR Host                172.16.49.201
DTR Replica ID          a3a8ab213a8b
SSO                     false
Create Repo on Push     true
Log Protocol            internal
Log Host                
Log Level               INFO
Storage Volume          
NFS Host                
NFS Path                
Scanning Enabled        false
Scanning Online Sync    true
Scanning Auto Recheck   false
```
### Settings Set Commands

  ```
  diver dtr settings set createrepo --setting={false/true}
  ```  
  Create a repository on Push

  ```
  diver dtr settings set online --setting={false/true}
  ```
  Scanning online sync
  
  ```
  diver dtr settings set scanning --setting={false/true}   
  ```
  DTR Image scanning