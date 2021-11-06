# Kubernetes Job Remover

## Description
This Custom Resource deletes completed jobs in the specified namespace that elapsed specified time.

## Sample
You can specify TTL as minutes

```yaml
apiVersion: remover.onikle.com/v1
kind: JobRemover
metadata:
  name: jobremover-sample
spec:
  namespace: "default"
  TTL: 30
```

## Installation
WIP
