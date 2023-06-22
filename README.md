# k8s-wait-for-multi

[![GitHub release (release name instead of tag name)](https://img.shields.io/github/v/release/erayan/k8s-wait-for-multi?logo=github)](https://github.com/EraYaN/k8s-wait-for-multi/releases) 
[![GitHub Workflow Status (with event)](https://img.shields.io/github/actions/workflow/status/erayan/k8s-wait-for-multi/release.yml?logo=githubactions)](https://github.com/EraYaN/k8s-wait-for-multi/actions/workflows/release.yml) [![Docker Image Version (latest semver)](https://img.shields.io/docker/v/erayan/k8s-wait-for-multi?sort=semver&logo=docker&label=docker)
](https://hub.docker.com/r/erayan/k8s-wait-for-multi)

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

For jobs it wait until the `Completed` condition is true.

For services it will wait until all pods that match the service selector are Ready and Available (like above).

## Docker

Hosted on Docker Hub: https://hub.docker.com/r/erayan/k8s-wait-for-multi

There are `latest`, `nonroot`, `<tag>` and `<tag>-nonroot` labels available, in `amd64` (`v2`), `arm` (`v7`) and `arm64`.