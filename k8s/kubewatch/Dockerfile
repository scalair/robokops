FROM scalair/robokops-base:0.2.0
COPY src /home/builder/src
RUN sudo chown -R builder:builder /home/builder
RUN echo 'kubewatch' | sudo tee /name
VOLUME /conf