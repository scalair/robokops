FROM scalair/robokops-base:0.2.0
COPY src /home/builder/src
RUN sudo chown -R builder:builder /home/builder
RUN echo 'elastic-stack' | sudo tee /name
VOLUME /conf