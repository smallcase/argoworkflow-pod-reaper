# argoworkflow-pod-reaper

Deletes completed pods that are owned by ArgoWorkflow.

## Usage:

```
go build -o app 

./app --in-cluster=false --delete-failed-after=1 --delete-successful-after=1 --namespaces="ns1,ns2"
```

Example command output:
```
2021/08/24 20:43:06 Deleting pods that have failed after 1 days and succeeded after 1 days
2021/08/24 20:43:06 Would delete pod static-tester in namespace ns1
2021/08/24 20:43:06 Would delete pod static-tester in namespace ns2
```

| flag                      | default  | type     | usage                                                   |
|---------------------------|----------|----------|---------------------------------------------------------|
| --in-cluster              | true     | bool     | use in cluster config, or ~/.kube/config |
| --delete-failed-after     | 10       | int      | delete failed pods after x days                         |
| --delete-successful-after | 5        | int      | delete successful pods after x days                     |
| --namespaces              | default  | []string | namespaces to watch                                     |




