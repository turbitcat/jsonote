FROM golang:latest

WORKDIR /tmp/build
RUN git clone https://github.com/turbitcat/jsonote.git
WORKDIR /tmp/build/jsonote
RUN go build
RUN cp jsonote /root
RUN rm -rf /tmp/build
WORKDIR /root

EXPOSE 8088

ENTRYPOINT ["/root/jsonote"]