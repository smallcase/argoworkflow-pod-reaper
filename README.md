# argoworkflow-pod-reaper

Deletes completed pods that are owned by ArgoWorkflow.

## Usage:

```
go build -o app 

./app --in-cluster=false --delete-failed-after=1 --delete-successful-after=1 --namespaces="ns1,ns2"
```

**Example command output:**
```
2021/08/24 20:43:06 Deleting pods that failed 1 days ago
2021/08/24 20:43:06 Deleting pods that succeeded 1 days ago
2021/08/24 20:43:06 Would delete pod static-tester in namespace ns1
2021/08/24 20:43:06 Would delete pod static-tester in namespace ns2
```

**Container image build:**
```
chmod +x build.sh
./build.sh
```

**Kubernetes job example:**
```yaml
apiVersion: batch/v1
kind: Job
metadata:
  generateName: example-job
spec:
  template:
    spec:
      containers:
      - name: example
        image: $MY_IMAGE:$TAG
        command: ["reaper"]
      restartPolicy: Never
      serviceAccountName: example
  backoffLimit: 0
```

**ArgoWorkflow example:**
```yaml
apiVersion: argoproj.io/v1alpha1
kind: Workflow                  
metadata:
  generateName: hello-world-    
spec:
  serviceAccountName: example
  entrypoint: example          
  templates:
  - name: example             
    container:
      image: $MY_IMAGE:$TAG
      command: [reaper]
```

**Required RBAC:**
```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: example
  namespace: default
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: example
rules:
- apiGroups:
  - ""
  resources:
  - pods
  verbs:
  - get
  - list
  - watch
  - delete
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: example
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: example
subjects:
- kind: ServiceAccount
  name: example
  namespace: default
  ```

---

| flag                      | default  | type     | usage                                                   |
|---------------------------|----------|----------|---------------------------------------------------------|
| --in-cluster              | true     | bool     | use in cluster config, or ~/.kube/config 				|
| --delete-failed-after     | 10       | int      | delete failed pods after x days                         |
| --delete-successful-after | 5        | int      | delete successful pods after x days                     |
| --namespaces              | default  | []string | namespaces to watch                                     |




