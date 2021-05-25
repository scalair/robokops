FROM scalair/robokops-base:0.6.2

# jsonnet
RUN go get github.com/google/go-jsonnet/cmd/jsonnet

# jb
RUN go get github.com/jsonnet-bundler/jsonnet-bundler/cmd/jb

# gojsontoyaml
RUN go get github.com/brancz/gojsontoyaml

COPY src /home/builder/src
RUN sudo chown -R builder:builder /home/builder
RUN echo 'monitoring' | sudo tee /name
VOLUME /conf