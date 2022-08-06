# sdview

Merge information build pod and the screwdriver information.

## Build
```
$ go build -o sdview cmd/sdview/main.go
```

## Installation

TBD

## Usage

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
