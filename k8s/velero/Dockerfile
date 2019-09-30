FROM scalair/robokops-base:0.1.0
COPY src /home/builder/src
RUN sudo chown -R builder:builder /home/builder
RUN echo 'velero' | sudo tee /name
VOLUME /conf