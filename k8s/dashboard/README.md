# Dashbord
Deploy [Kubernetes dashboard](https://github.com/kubernetes/dashboard), which is an administration dashboard with some basic metrics.
Once the dashboard has been deployed, you can access it using port-forward at [http://127.0.0.1:9090](http://127.0.0.1:9090/), for example:
```
POD_NAME=$(kubectl get pods -n kube-system -l "app=kubernetes-dashboard,release=kubernetes-dashboard" -o jsonpath="{.items[0].metadata.name}")
TOKEN=$(kubectl -n kube-system get secret $(kubectl -n kube-system get secret | grep admin-token | awk '{print $1}') -o jsonpath="{.data.token}" | base64 --decode)
echo "Token: ${TOKEN}"
kubectl -n kube-system port-forward $POD_NAME 9090:9090
```
The token displayed after running those command is what you have to enter to authenticate to the dashboard.