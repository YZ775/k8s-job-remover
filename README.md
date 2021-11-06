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

## Installation(WIP)
### Requirement
```
Go
Docker
Kubebuilder
kubectl
```
<br><br>
Before installation, set kubectl context as desired cluster's one.<br><br><br>

Make Docker image
```bash
make docker-build
```
<br><br>
Then, upload controller image to some container registry, and edit ```image``` field in config/manager/manager.yaml

If you are using kind, the command below is fine.
```bash
kind load docker-image controller:latest
```
<br><br>

Apply CRDs
```bash
make install
make deploy
```
<br><br>

If you want to deploy JobRemover as a sample setting,please use this command.
```bash
kubectl apply -f config/samples/remover_v1_jobremover.yaml
```
