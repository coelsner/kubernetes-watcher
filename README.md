# Kubernetes Watcher

This tool uses the [kubernetes API](https://kubernetes.io/docs/reference/using-api/api-concepts/) 
to watch for changes in the different resource types.

Currently, pods and events are watched, which for example can be used to identify failing pods 
and the reason of the failure.

By default the `default` namespace is tracked, but which can be changed via the `-namespace <other namespace>` commandline flag.

## Example output

```
INFO 2021/09/24 15:07:30 Started watching kubernetes ... Cancel with CTRL+C
INFO 2021/09/24 15:07:30 Events ResourceVersion: 141898164
EVENT 2021/09/24 15:07:31 EVENT elasticsearch.16a01a931315f5a5 - Failed build model due to ingress: test/elasticsearch: none certificate found for host: elasticsearch.test.org
EVENT 2021/09/24 15:07:31 EVENT k8s-test-database-a39b4b3649.16a01a93425fbc76 - Successfully reconciled
EVENT 2021/09/24 15:07:31 EVENT k8s-test-elastics-1cbcddde66.16a01a93297642f2 - Successfully reconciled
INFO 2021/09/24 15:07:31 Pods ResourceVersion: 141898164
EVENT 2021/09/24 15:07:31 POD 'database-backup-7dc5cb6887-vrkcj' has failed: The node was low on resource: ephemeral-storage.
EVENT 2021/09/24 15:07:31 POD 'database-backup-c5466879f-jnz7c' has failed: The node was low on resource: ephemeral-storage.
EVENT 2021/09/24 15:07:31 POD 'database-backup-restart-1632448800-n8sww' was successful
EVENT 2021/09/24 15:07:31 POD 'database-backup-restart-1632276000-x6scc' was successful
EVENT 2021/09/24 15:07:31 POD 'database-backup-568df6965c-9jptw' has failed: The node was low on resource: ephemeral-storage.
EVENT 2021/09/24 15:07:31 POD 'database-backup-54c8dd8469-q8jmj' has failed: The node was low on resource: ephemeral-storage.
EVENT 2021/09/24 15:07:31 POD 'database-backup-restart-1632362400-5kc8p' was successful
EVENT 2021/09/24 15:07:31 POD 'database-backup-7dc5cb6887-ffpsk' has failed: The node was low on resource: ephemeral-storage.
INFO 2021/09/24 15:07:41 Stopped watching kubernetes
```