FROM debian:bullseye AS builder
RUN echo "deb http://www.deb-multimedia.org bullseye main" >> /etc/apt/sources.list && \
  export DEBIAN_FRONTEND=noninteractive && \
  ln -fs /usr/share/zoneinfo/Europe/London /etc/localtime && \
  apt update -oAcquire::AllowInsecureRepositories=true && \
  apt install -y --allow-unauthenticated deb-multimedia-keyring && \
  apt install -y --allow-unauthenticated ca-certificates libmagickwand-7-dev wget && \
  wget https://golang.org/dl/go1.16.2.linux-amd64.tar.gz && \
  rm -rf /usr/local/go && tar -C /usr/local -xzf go1.16.2.linux-amd64.tar.gz
ADD ./ /src/celexacreams
WORKDIR /src/celexacreams
RUN export PATH=$PATH:/usr/local/go/bin && go build -o /celexacreams .

FROM debian:bullseye AS final
RUN apt update && apt install -y ca-certificates
COPY --from=builder /usr/lib/x86_64-linux-gnu/ /usr/lib/x86_64-linux-gnu/
COPY --from=builder /lib/x86_64-linux-gnu/ /lib/x86_64-linux-gnu/
COPY --from=builder /celexacreams /
CMD ["/celexacreams"]