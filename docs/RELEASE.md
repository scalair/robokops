# Release process

Thos document describes the release process for features and also for the entire project

## Releasing a feature

Releasing a feature means updating the code for that feature in a branch, and merging that branch to the master so that Gihub Action will create the corresponding Docker image and push it to Docker Hub.

In order to enable this process, the branch must be named precisely:

```
<feature_name>/<feature_version>
# example
velero/0.4.1
```
The version follows the semantic versioning convention, and in order to know the next version to apply, read the feature `CHANGELOG.md`.

Important :
- Update the `CHANGELOG.md` of the feature before merging the branch
- Update the `bom.yaml` file for the new feature version to be used

## Releasing robokops

Github Action will create a release when a tag is created.
