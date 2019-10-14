FROM scalair/robokops-base:0.1.1
COPY src /home/builder/src
RUN sudo chown -R builder:builder /home/builder
RUN echo 'cluster-init' | sudo tee /name
VOLUME /conf