# kubectl-lab

A kubectl plugin available like `kubectl get`.

## Installation

```
$ go get github.com/micnncim/kubectl-lab/cmd/kubectl-lab
```

## Usage

```
$ kubectl lab pods 
NAMESPACE   NAME                       AGE     READY   LABELS
default     dev-app-587f4d4cb5-9cwqt   3d10h           app=app,pod-template-hash=587f4d4cb5
default     dev-app-587f4d4cb5-cz6dm   3d10h           app=app,pod-template-hash=587f4d4cb5
default     dev-app-587f4d4cb5-xbzqm   3d10h           app=app,pod-template-hash=587f4d4cb5
```
