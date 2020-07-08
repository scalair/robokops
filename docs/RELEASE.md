# Release process

Thos document describes the release process for features and also for the entire project

## Releasing the base Robokops image

When working on the base image (`docker` folder), the branch must be named:
```
base/<version>
# example
base/0.5.1
```

When this branch is merged, Github Action will modify all the features images to use the new base image. All features will then be automatically released.

## Releasing a feature

A feature, in Robokops vocabulary, is a component that is deployed by Robokops, like velero, terraform or elastic-stack.
Features are found in `k8s` directory, except for the terraform feature that can be found in `terraform` directory.

Releasing a feature means updating the code for that feature in a branch, and merging that branch to the master so that Gihub Action will create the corresponding Docker image and push it to Docker Hub.

In order to enable this process, the branch must be named precisely:

```
<feature_name>/<feature_version>
# example
velero/0.4.1
```
The version follows the semantic versioning convention, and in order to know the next version to apply, read the feature `CHANGELOG.md`.

Important :
- Update the `CHANGELOG.md` found in the directory of the feature before merging the branch
- Update the `bom.yaml` file for the new feature version to be used

## Releasing robokops

When working on Robokops (but not on features), the branch must be named:

```
release/<release_version>
# example
release/0.8.1
```

Dont forget to update the `CHANGELOG.md`.

Github Action will create a release when you create the tag, after merging the corresponding branch.
