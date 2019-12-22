
```
docker service create --replicas 1 --name helloworld alpine ping docker.com
```

Suggestions:
 * validate cluster: kops validate cluster
 * list nodes: kubectl get nodes --show-labels
 * ssh to the master: ssh -i ~/.ssh/id_rsa admin@api.jim4-kops.k8s.local
 * the admin user is specific to Debian. If not using Debian please use the appropriate user based on your OS.
 * read about installing addons at: https://github.com/kubernetes/kops/blob/master/docs/addons.md.

Update kops to 500 pods per node:

```
kops edit cluster

spec:
  kubelet:
    maxPods: 500
```

