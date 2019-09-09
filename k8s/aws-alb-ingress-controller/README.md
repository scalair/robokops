# AWS Alb Ingress Controller
Deploy an AWS ALB Ingress Controller.

## How to use it
This is an example of how to create an ALB ingress:
```
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  annotations:
    kubernetes.io/ingress.class: alb
    alb.ingress.kubernetes.io/scheme: "internet-facing"
    alb.ingress.kubernetes.io/target-type: "instance"
    alb.ingress.kubernetes.io/certificate-arn: arn:aws:acm:eu-west-1:123456789012:certificate/abcdefgh-1234-5678-9012-ijklmnopqrst
    alb.ingress.kubernetes.io/listen-ports: '[{"HTTP":80,"HTTPS": 443}]'
    alb.ingress.kubernetes.io/actions.ssl-redirect: '{"Type": "redirect", "RedirectConfig": { "Protocol": "HTTPS", "Port": "443", "StatusCode": "HTTP_301"}}'
  name: foo
  namespace: foo
spec:
  rules:
  - host: foo.example.com
    http:
      paths:
      - path: /*
        backend:
          serviceName: ssl-redirect
          servicePort: use-annotation
      - path: /*
        backend:
          serviceName: foo
          servicePort: https
```