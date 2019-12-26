FROM golang:1.12

RUN apt-get update \
 && apt-get install -y openjdk-8-jdk \
 && echo "deb [arch=amd64] http://storage.googleapis.com/bazel-apt stable jdk1.8" | tee /etc/apt/sources.list.d/bazel.list \
 && curl https://bazel.build/bazel-release.pub.gpg | apt-key add - \
 && apt-get update \
 && apt-get install -y \
    bazel \
    patch

RUN go get k8s.io/repo-infra/kazel

RUN echo "deb [signed-by=/usr/share/keyrings/cloud.google.gpg] http://packages.cloud.google.com/apt cloud-sdk main" | tee -a /etc/apt/sources.list.d/google-cloud-sdk.list \
 && curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key --keyring /usr/share/keyrings/cloud.google.gpg add - \
 && apt-get update -y \
 && apt-get install -y google-cloud-sdk

WORKDIR /go/src

ADD . /go/src/github.com/tony-yang/realtor-tracker

ENV GO111MODULE=on

EXPOSE 80
