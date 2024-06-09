# About containers
Compiled from various sources.

## History
In the beginning, there was a program. Let's call the program run.sh, and what we’d do is we’d copy it to a remote server, and we would run it. However, running arbitrary code on remote computers is insecure and hard to manage and scale. So we invented virtual private servers and user permissions. And things were good.

But little run.sh had dependencies. It needed certain libraries to exist on the host. And it never worked quite the same remotely and locally. So we invented AMIs (Amazon Machine Images) and VMDKs (VMware images) and Vagrantfiles and so on, and things were good.

Well, they were kind-of good. The bundles were big and it was hard to ship them around effectively because they weren’t very standardised. And so, we invented caching.

Caching is what makes Docker images so much more effective than vmdks or vagrantfiles. It lets us ship the deltas over some common base images rather than moving whole images around. It means we can afford to ship the entire environment from one place to another. It’s why when you `docker run whatever` it starts close to immediately even though whatever described the entirety of an operating system image.

## Definition
Containers are lightweight, portable units that encapsulate an application and its dependencies (libraries, configuration files, binaries, etc.).
They ensure that the software runs consistently regardless of the environment, be it a developer’s local machine, a testing environment, or a production server.

## Benefits
- Portability: Containers can run on any system that has the container runtime, ensuring consistent environments across development, testing, and production.
- Isolation: Containers isolate applications, ensuring that they do not interfere with each other or the host system.
- Scalability: Containers can be easily scaled up or down to handle varying loads.
- Efficiency: Containers use fewer resources compared to VMs since they share the host OS kernel.

## Key Components and Terminologies

### Image
A container image is a static snapshot of a container's filesystem and the configuration metadata needed to run an application within a container. Container images are read-only and do not change once they are created. This immutability ensures consistency and repeatability. Images are built in layers. Each layer represents a set of filesystem changes (such as adding, modifying, or deleting files) and is stacked on top of the previous one. This layering system helps with efficiency and reusability. For example, common base layers (like an OS layer) can be shared between multiple images, reducing storage usage. The starting point for building an image. It could be a minimal operating system, like Alpine Linux, or a more feature-rich environment. A JSON file that contains metadata about the image, including details about the layers, environment variables, default command to run, and other configuration parameters.

### Namespaces
Namespaces provide the isolation needed to run multiple containers on one machine while giving each what appears like it’s own environment. There are - at the time of writing - six namespaces. Each can be independently requested and amounts to giving a process (and its children) a view of a subset of the resources of the machine.

#### PID
The pid namespace gives a process and its children their own view of a subset of the processes in the system. Think of it as a mapping table. When a process in a pid namespace asks the kernel for a list of processes, the kernel looks in the mapping table. If the process exists in the table the mapped ID is used instead of the real ID. If it doesn’t exist in the mapping table, the kernel pretends it doesn’t exist at all. The pid namespace makes the first process created within it pid 1 (by mapping whatever its host ID is to 1), giving the appearance of an isolated process tree in the container.

#### MNT
In a way, this one is the most important. The mount namespace gives the process’s contained within it their own mount table. This means they can mount and unmount directories without affecting other namespaces (including the host namespace). More importantly, in combination with the pivot_root syscall - as we’ll see - it allows a process to have its own filesystem. This is how we can have a process think it’s running on ubuntu, or busybox, or alpine — by swapping out the filesystem the container sees.

#### NET
The network namespace gives the processes that use it their own network stack. In general only the main network namespace (the one that the processes that start when you start your computer use) will actually have any real physical network cards attached. But we can create virtual ethernet pairs — linked ethernet cards where one end can be placed in one network namespace and one in another creating a virtual link between the network namespaces. Kind of like having multiple ip stacks talking to each other on one host. With a bit of routing magic this allows each container to talk to the real world while isolating each to its own network stack.

#### UTS
The UTS namespace gives its processes their own view of the system’s hostname and domain name. After entering a UTS namespace, setting the hostname or the domain name will not affect other processes.

#### IPC
The IPC Namespace isolates various inter-process communication mechanisms such as message queues.

#### USER
The user namespace is likely the most powerful from a security perspective. The user namespace maps the uids a process sees to a different set of uids (and gids) on the host. This is extremely useful. Using a user namespace we can map the container's root user ID (i.e. 0) to an arbitrary (and unprivileged) uid on the host. This means we can let a container think it has root access - we can even actually give it root-like permissions on container-specific resources - without actually giving it any privileges in the root namespace. The container is free to run processes as uid 0 - which normally would be synonymous with having root permissions - but the kernel is actually mapping that uid under the covers to an unprivileged real uid. Most container systems don't map any uid in the container to uid 0 in the calling namespace: in other words there simply isn't a uid in the container that has real root permissions.

### Cgroups
Fundamentally cgroups collect a set of process or task ids together and apply limits to them. Where namespaces isolate a process, cgroups enforce fair resource sharing between processes.

Cgroups are exposed by the kernel as a special file system you can mount. You add a process or thread to a cgroup by simply adding process ids to a tasks file, and then read and configure various values by essentially editing files in that directory.

