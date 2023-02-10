# K8s filtering.

## Motivation

Tetragon is configured via `TracingPolicies`. Broadly speaking,
`TracingPolicies` define _what_ events Tetragon should react to, and _how_. The
_what_ can be, for example, specific system calls with specific argument
values. The _how_, can be generating an event (observability), or killing the
process (enforcement).

Here, we are concerned with applying tracing policies only on a subset of pods
running on the system baesd on their namespace, and, in future work, their
labels.

To this end, a new type of policy is introduced: `TracingPolicyNamespaced` that
is exactly the same as the existing `TracingPolicy`, but is _only_ applied to
pods of the namespace that the policy is defined.

As is the case with `TracingPolicies`, `TracingPoliciesNamespaced` are
implemented in the kernel with BPF. This is important for both observability
and enforcement use-cases. For observability, copying only the relevant events
from kernel- to user-space reduces overhead. Second, for the enforcement
use-case, performing the enforcement action in the kernel avoids the
race-condition of doing it in user-space. For example, let us consider the case
where we want to block an application from performing a system call. Performing
the filtering in-kernel means that the application will never finish executing
the system call, which is not possible if enforcement happens in user-space.

To ensure that namespaced tracing policies are always correctly applied,
tetragon needs to perform actions before containers start executing. Tetragon
supports this via [OCI runtime
hooks](https://github.com/opencontainers/runtime-spec/blob/main/config.md#posix-platform-hooks).
If such hooks are not added, Tetragon will apply policies in a best-effort
manner using information from the k8s API server.


## Demo

For this demo, we use containerd and configure appropriate run-time hooks using minikube.

First, let us stars minikube, build and load images, and install tetragon and OCI hooks:
```
$ minikube start --container-runtime=containerd
$ ./contrib/rthooks/minikube-containerd-install-hook.sh
$ make image image-operator
$ minikube image load --daemon=true cilium/tetragon:latest cilium/tetragon-operator:latest
$ minikube ssh -- sudo mount bpffs -t bpf /sys/fs/bpf
$ helm install --namespace kube-system \
	--set tetragonOperator.image.override=cilium/tetragon-operator:latest \
	--set tetragon.image.override=cilium/tetragon:latest  \
	--set tetragon.enablePolicyFilter="true" \
	--set tetragon.grpc.address="unix:///var/run/cilium/tetragon/tetragon.sock" \
	tetragon ./install/kubernetes
```


The new CRD will be installed by the operator container:
```
$ kubectl -n kube-system logs -c tetragon-operator tetragon-xxxx
level=info msg="Tetragon Operator: " subsys=tetragon-operator
level=info msg="CRD (CustomResourceDefinition) is installed and up-to-date" name=TracingPolicy/v1alpha1 subsys=k8s
level=info msg="Creating CRD (CustomResourceDefinition)..." name=TracingPolicyNamespaced/v1alpha1 subsys=k8s
level=info msg="CRD (CustomResourceDefinition) is installed and up-to-date" name=TracingPolicyNamespaced/v1alpha1 subsys=k8s
level=info msg="Initialization complete" subsys=tetragon-operator
```

And the agent should report that the policyfilter (the low level mechanism that implements this) is
enabled:
```
...
level=info msg="Enabling policy filtering"
...
```

For illustration purposes, we will use the lseek system call with an invalid argument. Specifically
a file descriptor (the first argument) of -1.

Let us create a new namespace and start a pod:

```
$ kubectl run test --image=python -it --rm --restart=Never  -- /bin/bash
If you don't see a command prompt, try pressing enter.
root@test:/#
```

There is no policy installed so attempting to do the lseek operation will just return an error.
```
root@test:/# python
Python 3.11.2 (main, Feb  9 2023, 00:38:19) [GCC 10.2.1 20210110] on linux
Type "help", "copyright", "credits" or "license" for more information.
>>> import os
>>> os.lseek(-1,0,0)
Traceback (most recent call last):
  File "<stdin>", line 1, in <module>
OSError: [Errno 9] Bad file descriptor
>>> 
```

In another terminal, we install the policy:
```
$ cat << EOF | kubectl apply -n default -f -
apiVersion: cilium.io/v1alpha1
kind: TracingPolicyNamespaced
metadata:
  name: "lseek-namespaced"
spec:
  kprobes:
  - call: "__x64_sys_lseek"
    syscall: true
    args:
    - index: 0
      type: "int"
    selectors:
    - matchArgs:
      - index: 0
        operator: "Equal"
        values:
        - "-1"
      matchActions:
      - action: Sigkill
EOF
tracingpolicynamespaced.cilium.io/lseek-namespaced created
```

Then, attempting the lseek operation on the previous terminal, will result in the process getting
killed:
```
>>> os.lseek(-1, 0, 0)
Killed
```

The same is true for a newly started container:

```
kubectl run test --image=python -it --rm --restart=Never  -- /bin/bash
If you don't see a command prompt, try pressing enter.
root@test:/# python
Python 3.11.2 (main, Feb  9 2023, 00:38:19) [GCC 10.2.1 20210110] on linux
Type "help", "copyright", "credits" or "license" for more information.
>>> import os
>>> os.lseek(-1, 0, 0)
Killed
```

Dong the same on another namespace, however, will not kill the process:
```
$ kubectl create namespace test
$ kubectl -n test run test --image=python -it --rm --restart=Never  -- /bin/bash
If you don't see a command prompt, try pressing enter.
root@test:/# python
Python 3.11.2 (main, Feb  9 2023, 00:38:19) [GCC 10.2.1 20210110] on linux
Type "help", "copyright", "credits" or "license" for more information.
>>> import os
>>> os.lseek(-1, 0, 0)
Traceback (most recent call last):
  File "<stdin>", line 1, in <module>
OSError: [Errno 9] Bad file descriptor
```
