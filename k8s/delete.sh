#!/bin/bash
set -e

echo "🗑️  Deleting Daybook from Kubernetes"
echo ""

read -p "Are you sure you want to delete everything? (y/N) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Cancelled"
    exit 0
fi

echo "Deleting all resources..."
kubectl delete namespace daybook

echo ""
echo "✅ All resources deleted"
