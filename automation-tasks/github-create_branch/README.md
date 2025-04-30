## What it does 

This automation helps to create a github branch from a specific branch. This can be used in platform for auto branch creation or default branch setup along with repo setup

## How to execute it?

Pre-requisite:

````sh
export GHUB_TOKEN = "<Github token with appropriate access>"

````

### CLI Access:

````sh
$ go run cli_main.go  --repo-name rajeshdeshpande02/test-repo --base-branch main --new-branch test1

````
### REST API Access:

Start Server:
````sh
go run rest_main.go

````

Send Request:
````sh
curl --location 'https://expert-guacamole-9j5ppg5j7qrc7v5p-8080.app.github.dev/create-branch' \
--header 'Content-Type: application/json' \
--data '{
    "repo-name": "rajeshdeshpande02/test-repo",
    "base-branch": "main",
    "new-branch": "rest-test"
}'
````

