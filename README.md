# Operator to delete pod
## Assignment problem statement 
    Write operator to delete the pod, by taking pod name as input in spec CR

## Approach
1. Operator-sdk is the framework used to write the operator.
2. Custom respource with api-version poddelete.example.com/v1alpha1 and kind: PodDelete is created
3. Spec contains namespace: where pod recides and podName: name of the pod
4. Reconcile method is called by api server for any operation to PodDelete object
5. Pod delete logic is implemented in reconcile funtion - pkg/controller/poddelete/poddelete_controller.go
    ```
    func (r *ReconcilePodDelete) Reconcile(request reconcile.Request) (reconcile.Result, error) {
    ```
6. Currently operator watches all the namespace and can delete the pod in any namespace. Can be restricted via roles.
7. Operator is able to delete itself, but this can be ignored.
8. If pod is not found or is already in terminating state delete operation is ignored.

## Steps to deploy
### Apply crd, cluster-role, cluster-role-bindings, service-accounts for operator
```
$ kubectl apply -f deploy/crds/poddelete.example.com_poddeletes_crd.yaml
$ kubectl apply -f deploy/service_account.yaml
$ kubectl apply -f deploy/role.yaml
$ kubectl apply -f deploy/role_binding.yaml
```

### Deploy operator
```
$ kubectl apply -f deploy/operator.yaml
```

### Sample custom resource
```
apiVersion: poddelete.example.com/v1alpha1
kind: PodDelete
metadata:
  name: <name-of-cr>
spec:
  namespace: <namespace where pod resides>
  podName: <pod-name>
```

## Example
#### Create sample nginx pod
```
$ kubectl run nginx --image=nginx
```

#### Delete the pod nginx
```
$ kubectl apply -f <(echo "
apiVersion: poddelete.example.com/v1alpha1
kind: PodDelete
metadata:
  name: delete-nginx-pod
spec:
  namespace: default
  podName: nginx
")
```
