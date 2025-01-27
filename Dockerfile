FROM alpine:edge
ARG TARGETPLATFORM
LABEL MAINTAINER="i@nn.ci"
WORKDIR /bin/
COPY bin/${TARGETPLATFORM}/dp ./
RUN chmod +x /bin/dp
CMD ["/bin/dp"]