## What it does 

This automation helps to create a github repo from a specific repo template. This can be used in platform for auto repo creation from template

## How to execute it?

Pre-requisite:

````sh
export GHUB_TOKEN = "<Github token with appropriate access>"

````

### CLI Access:

````sh
$ go run cli_main.go  --template rajeshdeshpande02/app-template-java --repo_name test-repo1

````
### REST API Access:

Start Server:
````sh
go run rest_main.go

````

Send Request:
````sh
curl --location 'https://expert-guacamole-9j5ppg5j7qrc7v5p-8080.app.github.dev/repo-with-template' \
--header 'Content-Type: application/json' \
--data '{
    "template": "rajeshdeshpande02/app-template-java",
    "repo_name": "java-repo1"
}'
````

