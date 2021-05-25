FROM scalair/robokops-base:0.6.2
COPY src /home/builder/src
RUN sudo chown -R builder:builder /home/builder
RUN echo 'elastic-stack' | sudo tee /name
VOLUME /conf