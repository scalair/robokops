#!/bin/bash

# Dashboard
helm upgrade --install --force --namespace kube-system -f kubernetes-dashboard/values.yaml kubernetes-dashboard stable/kubernetes-dashboard --version 1.8.0 --wait

# Metrics server
if [ "${METRICS_SERVER}" = "true" ]; then
	kubectl apply -f heapster/heapster.yaml
	kubectl apply -f heapster/influxdb.yaml
	kubectl apply -f heapster/heapster-rbac.yaml
fi

# Display token to access the dashboard
echo "To access the admin dashboard run:"
cat << 'EOF'
POD_NAME=$(kubectl get pods -n kube-system -l "app=kubernetes-dashboard,release=kubernetes-dashboard" -o jsonpath="{.items[0].metadata.name}")
TOKEN=$(kubectl -n kube-system get secret $(kubectl -n kube-system get secret | grep admin-token | awk '{print $1}') -o jsonpath="{.data.token}" | base64 --decode)
echo http://127.0.0.1:9090/
echo "Token: ${TOKEN}"
kubectl -n kube-system port-forward $POD_NAME 9090:9090
EOF
