FROM scalair/robokops-base:0.6.2
COPY src /home/builder/src
RUN sudo chown -R builder:builder /home/builder
RUN echo 'aws-efs-csi-driver' | sudo tee /name
VOLUME /conf