# Terraform
Robokops container packed with Terragrunt and Terraform that can be used to easily create and destroy clusters. This can make sense if you want to automatically create a cluster in the morning and destroy it in the evening.
If you already have a cluster, or if you are not planning to make changes very often to it then it probably won't make sense to use that feature.

## Entrypoint
You can run Terraform with the same action that you would normaly do:
* plan
* apply
* destroy

The only difference is that it will run Terragrunt under the hood with `-all` at the end (`plan-all`, `apply-all` or `destroy-all`) which will run the action against all submodules.

## Secret
When creating a cluster with Terraform, you may want to create users in your cloud provider to access it. If so, you don't want the credentials for those users to be stored unencrypted in Terraform state, so you need a way to encrypt/decrypt them somehow.
The solution we found was to use [Keybase](https://keybase.io/). There is no obligation to use Keybase with Robokops, but we strongly recommand to encrypt sensitive data in state file.

In order to use Keybase, you need to create an account (it's free) and get your credentials.
Within the container, the script `/home/builder/src/decrypt.sh` can be used to decrypt a secret. 

### Example
In the file `example/terraform/post-apply.sh`:
```
# Keybase login
export KEYBASE_USERNAME=$(vault read -field KEYBASE_USERNAME secret/.../keybase)
export KEYBASE_PAPERKEY=$(vault read -field KEYBASE_PAPERKEY secret/.../keybase)
export KEYBASE_PASSPHRASE=$(vault read -field KEYBASE_PASSPHRASE secret/.../keybase)
keybase oneshot

# IAM EKS admin
cd /home/builder/src/eks
ACCESS_KEY_ID=$(terragrunt output iam_access_key_id)
terragrunt output iam_access_key_encrypted_secret | base64 -d > iam_access_key_encrypted_secret.txt
SECRET_ACCESS_KEY=$(/home/builder/src/decrypt.sh ${KEYBASE_PASSPHRASE} iam_access_key_encrypted_secret.txt)
rm iam_access_key_encrypted_secret.txt

vault write secret/.../admin AWS_ACCESS_KEY_ID=${ACCESS_KEY_ID} AWS_SECRET_ACCESS_KEY=${SECRET_ACCESS_KEY##*$'\n'}
```
