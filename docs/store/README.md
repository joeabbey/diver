
## Docker Store

### Interacting with Docker Store

Logging into the Docker Store through the following command:

`./diver store --username <user> --password <password>`


To retrieve the Docker Store User ID for this user use the following command:

`./diver store user`

The `ID` field is used for retrieving subscriptions and licenses, other users can be examined by using the `--user <username>` flag.

Retrieve subscriptions for this user with the following command:

`./diver store subscriptions ls --id <ID>`

To retrieve the first active subscription use the `--firstactive` flag.

To retrieve a subscription license use the following command:

`./diver store licenses get --subscription <SUBSCRIPTION>`

This will print the raw output so it is advisable to pipe this to a file with the following addition to the command:

`> ./subscription_ID.lic`
