# Ingress Nginx
Setup ingress using nginx (https://github.com/kubernetes/ingress-nginx)
Documentation to deploy it can be found [here](https://kubernetes.github.io/ingress-nginx/deploy)

## AWS
On AWS, this will create an ELB. You can use that ELB as the one entrypoint to all your services. The downside to use ELB is that it only support one ACM certificate, which is not a problem if you can fit all your services under the same DNS subdomain (using a wildcard for the certificate) but won't work if you have different domain on your cluster.

## Example
This is an example to use the ingress nginx with a SSL redirection:
```
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  annotations:
  	kubernetes.io/ingress.class: nginx
    nginx.ingress.kubernetes.io/force-ssl-redirect: "true"
  name: my-app
  namespace: my-app
spec:
  rules:
  - host: app.example.com
    http:
      paths:
      - backend:
          serviceName: my-service
          servicePort: http
```