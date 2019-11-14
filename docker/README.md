# Robokops - Docker image base
Docker image used as a base for all features of Robokops

## What it contains
* aws-cli
* aws-iam-authenticator
* kubectl
* helm
* vault
* yq
* envsubst
* boxes
* tini

## Entrypoint
The file `entrypoint.sh` is the entrypoint of all features. It's doing all the common work so each containers can focus on what's specific to them.
What the script is doing:
* Takes all yaml files in `/conf`, templates them and merge them with manifests in `/home/builder/src`
* If `action=dry-run` return generated manifests
* Load configuration and authenticate to the cluster
* If `pre-manifests` folder exists, apply or delete them
* If `pre-$ACTION.sh` exists, run it
* Run `$ACTION.sh`
* If `post-$ACTION.sh` exists, run it
* If `post-manifests` folder exists, apply or delete them

`$ACTION` either match the `action` or the `terraform` value passed to Robokops. So `entrypoint.sh` will call the script named `apply.sh`, `deploy.sh`, `delete.sh`, etc.