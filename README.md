# session1-containers

Now we have a simple application we want to look into how we can run it as a container. There is no one way of doing this, and can differ depending on a few factors like language, runtime requirements, planned runtime environment, etc. So these next steps are meant to be just an example way that we can run a simple application like this one.

## Dockerfiles
But before we start look at our application, lets first discuss what is a container compared to a standard application, how do we build them, and what is the structure of the definition files used to do so (known as a `Dockerfile`).

Looking at the simple name server application, we have a very simple webserver which holds its "state" using a file. It is then quite easy to take our application and related configurations and run them somewhere on a server to be used in production. So why would I want to take this application and move it into a container, what do I benifit from by doing so?

To build a container creates an immutable object which contains all the runtime requirements (this can be anything from binaries for the runtime, static objects, etc), and wrappering it in a layer which can the be executed in lightwieght environment. This image can then be "deployed" and hosted making it accessible to others with a container runtime to be able to use.  This alone brings multiple advantages:

- Portability: This container image can run on any platform, so can be run anywhere very simply. This mean that we could have multiple instances of the application based on the same image running for different reasons. This could be a test/pre-prod/prod type setup, or even a distribution method for providing out an application to multiple consumers.
- Efficient runtimes: I dont have to setup a server for an application to run inside, and maintain that server. Instead I can some sort of cloud provider to manage my containers, making them very computational effiecient.
- Agility: Containers typically are extremely fast to bring up and down, making them very suitable for scaleable HA environments.
- Versioning: Rather than managing versions of an application via an install, container images are tagged (normally with thier version) making it very simple to upgrade, or just change version, with no setup cost.
- Security: Containers are very isolated in terms of access to the host machine running them, and other containers running.

Both Podman and Docker both provide a very simple to use set of commands for building container images. But to build a container we need to define the contents of the container. We do this using a `dockerfile`. Its important to know that container images are typical built on top of other container images, which allow it to inherit anything from the parent image. For example, we can use one of the openJDK images to provide a base image for my Java application, meaning that I dont have to be worried about how to install Java (and which version) into my application image. I would typically have something like this:

```
FROM openjdk:8

COPY /path/to/my/jar /path/for/jar/in/container

ENTRYPOINT ["java", "-jar", "/path/for/jar/in/container"] 
```
I have not had to be concerned with getting java 8 installed. But lets disect the rest of this file. The docker commands are always in upper case, and looking at the 3 here:
- `FROM`: This command states which image we are going to be building from as our starting point. In this case we sat `openjdk:8`, where the `openjdk` is the image name and the `8` is the version. (<Image_Name>:<Version>) 
- `COPY`: This command allows for copy any files from the host machine building the container (aka your computer in our case) and place them inside the container.
-`ENTRYPOINT`: This is typically the executable you want to run to start the container/application up. The `ENTRYPOINT` can be used in conjunction with the `CMD`, which typically provides defaults for your executable in the `ENTRYPOINT`. 

A full reference to all the available commands can be found here: https://docs.docker.com/engine/reference/builder/

## Building our App
The first decision to make, is how to we compile our program? Depending on how we choose to compile the application might change how we choose to containerise our program. So here is two examples:

### Approach one - expect artifacts
This first approach makes the assumption that another step in out automation has already gone away and built our application and provided the build artifacts in the standard(or expected) location. This is the simpliest approach to building container images, as it is pretty much a "lift and shift" in terms of the artifacts needed at runtime. This makes our Dockerfile extremely simple:
```
FROM scratch

COPY bin/server-linux-amd64 /server
COPY assets /assets
COPY names.json /names.json

ENTRYPOINT["./server"]
```

In this file we start from an image called `scratch` which is the smallest simpliest container we can use. Golang compiles into executable binaries so we need very little to run them. This means the resulting container once built is only 6.48MB.

### Approach 2 - compile and build

If we look in the `dockerfileC` file:
```
############################
# STEP 1 build executable binary
############################
FROM golang AS builder
WORKDIR $GOPATH/src/github.com/containers-kubernetes-education/session1-containers
COPY . /go/src/github.com/containers-kubernetes-education/session1-containers

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
ENV GO111MODULE=auto
ENV GOPATH=/go

WORKDIR $GOPATH/src/github.com/containers-kubernetes-education/session1-containers
COPY . /go/src/github.com/containers-kubernetes-education/session1-containers

RUN go build -o /go/bin/run /go/src/github.com/containers-kubernetes-education/session1-containers/cmd/main.go

############################
# STEP 2 build a small image
############################
FROM scratch
COPY --from=builder /go/bin/run /go/bin/run
COPY assets /assets
COPY names.json /names.json
ENTRYPOINT ["/go/bin/run"]
```
A lot more is going on. In this dockerfile I actually make use of two containers for my server. The image `golang` for building and compiling my code, and the `scratch` for my runtime image. I do this as the golang image is quite large (948MB) as it containers all the required build tools for comiling go on any platform. But at runtime I dont need any of this, so i copy the compiled binary i need out to a much much smaller container.

This is a very useful trick as it allows me to compile my code against the latest version of the compiler (you may notice that i didnt include a version on my `FROM golang:<version>`, so it defaults to latest), without having to ever change anything in my automation or build processes. This also means that you can compile golang on your computer without every installing golang. This is something I use quite a lot to avoid installing lots of tools which quickly get outdated onto my PC. I use provided containers (assuming they exist) to compile my code for me.

### Building the image
To now actually build the image we can run the command:
```
podman build -t myimage:v1 .
```
which uses the current directory as context, looking for a dockerfile by default. We can specify a dockerfile however like this:
```
podman build -t myimage:v1 -f dockerfileC .
```
For running the application please move to the branch `step2-running-the-application`