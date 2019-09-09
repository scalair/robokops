# Vault
The recommended way to manage secrets with Robokops is [Vault](https://www.vaultproject.io/). Of course you can use your own secret management tool, but you may need to make some changes.
Vault is already installed in all containers, so you can easily retrieve secrets in `.conf` files like that:
```
export AWS_ACCESS_KEY_ID=$(vault read -field AWS_SECRET_ACCESS_KEY secret/.../aws)
export AWS_SECRET_ACCESS_KEY=$(vault read -field AWS_SECRET_ACCESS_KEY secret/.../aws)
```
In addition to that, you have to configure Vault with the URL and the token to use with environment variables. For instance:
```
robokops --config example/ --terraform apply --env VAULT_TOKEN=$VAULT_TOKEN --env VAULT_ADDR=$VAULT_ADDR
``` 
or directly to a container:
```
docker run -e "VAULT_TOKEN=$VAULT_TOKEN" -e "VAULT_ADDR=$VAULT_ADDR" -v "$(pwd)/../example/conf/:/conf" robokops-terraform apply
```