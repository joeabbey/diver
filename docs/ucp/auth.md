### Creating users/organisations

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


### Working with Roles

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

### Working with Grants

**Listing**
To list all current grants you can use the following command:
`./diver ucp auth grants list`

To **resolve** the gran UUID to an actual `name` use the `--resolve` flag when listing grants.

**Creating**

To create a grant use the command `./diver ucp auth grants set` with the following flags:

`--collection` - Can either be a collection path or a Kubernetes namespace.

`--subject` - A user or service account.

`--role` - A role that has been created in UCP

`--type` - The type of grant that will be applied to, can be a `collection` grant, a single `namespace` grant or `all` kube namespaces.

**NOTE**: Unless the accounts are pre-configured UCP accounts then the UUIDs will need to be passed to this command.

#### EXAMPLE - Deploying HELM

**Before Installing Helm**

Create a Kubernetes service account:
`kubectl create serviceaccount --namespace kube-system tiller`

Create a grant for the `tiller` service account:

```
  ./diver ucp auth grants set --role fullcontrol       \
  --subject system:serviceaccount:kube-system:tiller  \
  --collection kubernetesnamespaces                    \
  --type all
```

Install (or init) Helm

`helm init`

Correct the service account

`kubectl patch deploy --namespace kube-system tiller-deploy -p ‘{“spec”:{“template”:{“spec”:{“serviceAccount”:”tiller”}}}}’`

Deploy using Helm! 

e.g MySQL deployment.

`helm install --name mysql stable/mysql`