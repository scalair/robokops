FROM blakelead/robokops-base:0.2.1
COPY src /home/builder/src
RUN sudo chown -R builder:builder /home/builder
RUN echo 'elastic-stack' | sudo tee /name
VOLUME /conf