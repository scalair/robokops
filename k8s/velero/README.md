# Velero
Setup [Velero](https://velero.io/) to backup cluster manifests and volumes.

## Prerequisites
A user with proper permission must have been created and you will have to configure it.

## Configuration

### AWS
Example to configure Velero in AWS:
```
configuration:
  provider: aws
  backupStorageLocation:
    name: aws
    bucket: my.awesome.bucket
    prefix: velero
    config:
      region: eu-west-1
  volumeSnapshotLocation:
    name: aws
    config:
      region: eu-west-1
  extraEnvVars:
    AWS_CLUSTER_NAME: my.awesome.bucket

credentials:
  secretContents:
    cloud: |
      [default]
      aws_access_key_id=${VELERO_AWS_ACCESS_KEY_ID}
      aws_secret_access_key=${VELERO_AWS_SECRET_ACCESS_KEY}
```