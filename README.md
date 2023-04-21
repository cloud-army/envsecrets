# envsecrets

Binary to inject secrets as environment variables to another application 

This uses the Workload Identity on GCP https://cloud.google.com/kubernetes-engine/docs/how-to/workload-identity#verify_the_setup
or default credentials

# Settings

By default, the same path of the envsecrets file is used, but you can optionally set the environment variable to change the default path.
- ```ENVSECRETS_CONFIG_PATH:``` set the environment variable to the path where the initial configuration file is located (optional).

For example:
```bash
export ENVSECRETS_CONFIG_PATH:"/opt/envsecrets-config.json"
```

- ```envsecrets-config.json:``` The input file is used to indicate the ```secrets``` to be read from the Secret Manager and set general ```config```

For example:
```
{
    "secrets":[
    {
        "env":"ENV1",
        "name":"projects/GCP_PROJECT_NUMBER/secrets/SECRET_STRING_NAME1/versions/latest"
    },
    {
        "env":"ENV2",
        "name":"projects/GCP_PROJECT_NUMBER/secrets/SECRET_STRING_NAME2/versions/latest"
    },
    {
        "env":"ENV3",
        "name":"projects/GCP_PROJECT_NUMBER/secrets/SECRET_STRING_NAME3/versions/latest"
    },
    {
        "env":"",
        "name":"projects/GCP_PROJECT_NUMBER/secrets/SECRET_JSON_NAME1/versions/latest"
    }
    ],
    "config":
    {
        "convert_to_uppercase_var_names": true
    }    
}
```
Where:


- ```secrets``` is an array of 1 or more with attributes corresponding to env and name for each secret to read from the secret manager.

Note: ```env``` is the name of the environment variable that will be injected and ```name``` is the full path where the secret to be read is located (GCP Secret Manager), if it contains a json secret, it will create an environment variable for each json key

# Usage

if you want to inject secrets to the app named "app1", use the following command line

```bash
./envsecrets ./app1
```

# Compile

In this process, it allows us to generate the envsecrets binary from the golang code, and depending on the distribution where we need to use it, we will execute the corresponding command

- For linux (from Mac)

```bash
go mod download
GOOS=linux GOARCH=amd64 go build envsecrets.go
```

- For Alpine: *

```bash
go mod download
go build -o envsecrets -ldflags "-linkmode external -extldflags -static" -a envsecrets.go
```
- For Linux (Debian and other ditros) *

```bash
go mod download
go build -o envsecrets envsecrets.go
```

*You can use the docker image golang:1.15.0-buster to compile to Alpine and Linux