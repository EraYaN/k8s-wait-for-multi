# k8s-wait-for-multi

This is an implementation of k8s-wait-for that allows you to wait for multiple items in one process, so you only need to start one init container.
This uses informers (wrapped in a `sigs.k8s.io/controller-runtime/pkg/cache.Cache`) to get the status updates for all the items that this application is waiting for.

It accepts arguments in the following formats:

- `namespace,service,service-name`
- `namespace,job,job-name`
- `namespace,pod,pod-name`
- `service,service-name` using the namespace from the `--namespace`, `-n` flag or `default`
- `job,job-name` using the namespace from the `--namespace`, `-n` flag or `default`
- `pod,pod-name` using the namespace from the `--namespace`, `-n` flag or `default`
- `pod-name` using the namespace from the `--namespace`, `-n` flag or `default` and the kind `pod` 

For pods it waits until the pod is Ready (`k8s.io/kubectl/pkg/util/podutils.IsPodReady`) and Available (`k8s.io/kubectl/pkg/util/podutils.IsPodAvailable`).

For jobs it wait intil the `Completed` condition is true.

For services it will wait until all pods that match the service selector are Ready and Available (like above).