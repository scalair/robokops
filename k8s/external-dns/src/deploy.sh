#!/bin/bash
helm upgrade --install --force --namespace kube-system -f external-dns/values.yaml external-dns stable/external-dns --version 2.4.2 --wait