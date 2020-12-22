FROM scalair/robokops-base:0.5.3
COPY src /home/builder/src
RUN sudo chown -R builder:builder /home/builder
RUN echo 'kubewatch' | sudo tee /name
VOLUME /conf