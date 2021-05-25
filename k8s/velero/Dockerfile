FROM scalair/robokops-base:0.6.5
COPY src /home/builder/src
RUN sudo chown -R builder:builder /home/builder
RUN echo 'velero' | sudo tee /name
VOLUME /conf