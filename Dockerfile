ARG ARCH=
FROM ${ARCH}/debian:bullseye AS builder
ENV GOPROXY direct
#ENV ARCH ${ARCH}
RUN echo "deb http://www.deb-multimedia.org bullseye main" >> /etc/apt/sources.list && \
  export DEBIAN_FRONTEND=noninteractive && \
  case $(uname -m) in \
    arm64) export ARCH="arm64";; \
    aarch64) export ARCH="arm64";; \
    x86_64) export ARCH="amd64";; \
    esac && \
  ln -fs /usr/share/zoneinfo/Europe/London /etc/localtime && \
  apt update -oAcquire::AllowInsecureRepositories=true && \
  apt install -y --allow-unauthenticated deb-multimedia-keyring git && \
  apt install -y --allow-unauthenticated ca-certificates libmagickwand-7-dev wget && \
  wget https://golang.org/dl/go1.16.2.linux-${ARCH}.tar.gz && \
  rm -rf /usr/local/go && tar -C /usr/local -xzf go1.16.2.linux-${ARCH}.tar.gz
ADD ./ /src/celexacreams
WORKDIR /src/celexacreams
RUN export PATH=$PATH:/usr/local/go/bin && go build -o /celexacreams .

FROM ${ARCH}/debian:bullseye AS intermediate-amd64
RUN apt update && apt install -y ca-certificates
COPY --from=builder /usr/lib/x86_64-linux-gnu/ /usr/lib/x86_64-linux-gnu/
COPY --from=builder /lib/x86_64-linux-gnu/ /lib/x86_64-linux-gnu/
COPY --from=builder /celexacreams /

FROM ${ARCH}/debian:bullseye AS intermediate-arm64v8
RUN apt update && apt install -y ca-certificates
COPY --from=builder /usr/lib/aarch64-linux-gnu/ /usr/lib/aarch64-linux-gnu/
COPY --from=builder /lib/aarch64-linux-gnu/ /lib/aarch64-linux-gnu/
COPY --from=builder /celexacreams /

FROM intermediate-${ARCH} AS final
CMD ["/celexacreams"]