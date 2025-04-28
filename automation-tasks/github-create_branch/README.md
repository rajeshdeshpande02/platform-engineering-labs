## What it does 

This automation helps to create a github branch from a specific branch. This can be used in platform for auto branch creation or default branch setup along with repo setup

## How to execute it?

Run following command:

````sh
 $ go run create_branch.go --repo-name=rajeshdeshpande02/test-repo  --base-branch develop --new-branch test

````

Pre-requisite:

````sh
export GHUB_TOKEN = "<Github token with appropriate access>"

````
