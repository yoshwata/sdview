# sdview

Merge information build pod and the screwdriver information.

## Build
```
$ go build -o sdview cmd/sdview/main.go
```

## Installation

Get the binary from Releases and put it on your bin path.

## Configuration

Put a yaml file `sdview_config.yaml` on your $HOME or the dir `sdview` exists.

example:
```yaml
usertoken: YOUR_SD_USER_TOKEN_HERE
sdapi: https://your.screwdriver.api.co.jp
```


## Usage

```
sdview

Usage:
  sdview [flags]
  sdview [command]

Examples:

        sdview -o="custom-columns=NAME:$.metadata.name,IMAGE:$.spec.containers[0].image" -b="custom-columns=builcClusterName:$.buildClusterName" -j="custom-columns=jobname:$.name" -e="custom-columns=causeMessage:$.causeMessage" -p="custom-columns=REPO:$.scmRepo.name"


Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  version     Print the version number of sdview

Flags:
      --as string                      Username to impersonate for the operation. User could be a regular user or a service account in a namespace.
      --as-group stringArray           Group to impersonate for the operation, this flag can be repeated to specify multiple groups.
      --as-uid string                  UID to impersonate for the operation.
      --cache-dir string               Default cache directory (default "/home/yoshwata/.kube/cache")
      --certificate-authority string   Path to a cert file for the certificate authority
      --client-certificate string      Path to a client certificate file for TLS
      --client-key string              Path to a client key file for TLS
      --cluster string                 The name of the kubeconfig cluster to use
      --context string                 The name of the kubeconfig context to use
  -h, --help                           help for sdview
      --insecure-skip-tls-verify       If true, the server's certificate will not be checked for validity. This will make your HTTPS connections insecure
      --kubeconfig string              Path to the kubeconfig file to use for CLI requests.
  -l, --maxLines int                   Max lines of table.
  -n, --namespace string               If present, the namespace scope for this CLI request
  -o, --output string                  Path of kubernetes pods response
      --request-timeout string         The length of time to wait before giving up on a single server request. Non-zero values should contain a corresponding time unit (e.g. 1s, 2m, 3h). A value of zero means don't timeout requests. (default "0")
  -e, --sdEventPath string             Path of sd's /events response
  -j, --sdJobPath string               Path of sd's /jobs response
  -p, --sdPipelinePath string          Path of sd's /pipelines response
  -b, --sdbuildPath string             Path of sd's /builds response
  -s, --server string                  The address and port of the Kubernetes API server
      --tls-server-name string         Server name to use for server certificate validation. If it is not provided, the hostname used to contact the server is used
      --token string                   Bearer token for authentication to the API server
      --user string                    The name of the kubeconfig user to use

Use "sdview [command] --help" for more information about a command.
```

### example

```
$ ./sdview -o="custom-columns=NAME:$.metadata.name,CPU:$.spec.containers[0].resources.limits.cpu" -b="custom-col
umns=builcClusterName:$.buildClusterName" -j="custom-columns=jobname:$.name" -e="custom-columns=causeMessage:$.causeMessage" -p="custom-columns=REPO:$.scmRepo.name" -l=5
  BUILDID (5)    CAUSEMESSAGE                              BUILCCLUSTERNAME   JOBNAME                REPO                           NAME             CPU
 -------------- ----------------------------------------- ------------------ ---------------------- ------------------------------ ---------------- -----
  84264161       Synchronized by pxxxxxx.git:txxxxxxxxxx   txxx-xxx-sdcd      PR-1309:pull-request   rxxxxx/rexxxxx                 84264161-vfudz   2
  84294563       Manually started by xxxxxxxxxx            txxx-xxx-sdcd      PR-365:test            Wxxxxx/dxxx                    84294563-znd0s   8
  84303575       Synchronized by ghxxxxxx:shxxxxxa         txxx-xxx-sdcd      PR-938:decide-stg      zxx/xxxxsts                    84303575-tdd0n   2
  84327038       Opened by gxxxxxxp:txxxxxxxxi             txxx-xxx-sdcd      PR-379:test            Waxxxx/dxxx                    84327038-fr12c   8
  84335066       Manually started by yxxxxxi               txxx-xxx-sdcd      release-build          yxx/feexxxxx                   84335066-axpv8   2
```
