FROM ubuntu
RUN export DEBIAN_FRONTEND=noninteractive && \
  ln -fs /usr/share/zoneinfo/Europe/London /etc/localtime && \
  apt update -y && \
  apt install -y ca-certificates libgraphicsmagick1-dev
ADD bin/celexacreams /
CMD ["/celexacreams"]