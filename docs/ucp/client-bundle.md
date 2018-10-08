## Using client bundles

Download the client bundle to your local machine.

```
./diver ucp client-bundle get
INFO[0000] Downloading the UCP Client Bundle            
```

List all client bundles

```
$ diver ucp client-bundle ls
ID                                                                 Label
14fb6d770d46ac9b9e431ad52b9f7b67daa2dc914be88baf79d96e0b6397e03f   Generated on Wed, 25 Jul 2018 10:25:52 UTC
acefe7997e5e0dda459d862053222c674f6f036cddc91cbbd801555b3bca9689   Generated on Tue, 31 Jul 2018 08:47:55 UTC
e1523a2642355ae5bd3eae97ca5fc54b51b6a16c6604f050455070c3c7392476   Generated on Tue, 12 Jun 2018 17:52:52 UTC
$
```

Rename or remove client bundles

**Remove**

```
$ diver ucp client-bundle rm <ID>

```

**Rename**

```
$ diver ucp client-bundle rename <Bundle_ID> <new_label>

```
