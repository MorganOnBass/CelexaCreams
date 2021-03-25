FROM ubuntu
RUN apt update -y && apt install -y ca-certificates
ADD bin/celexacreams /
CMD ["/celexacreams"]