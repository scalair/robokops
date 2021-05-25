FROM scalair/robokops-base:0.6.1
COPY src /home/builder/src
RUN sudo chown -R builder:builder /home/builder
RUN echo 'dashboard' | sudo tee /name
VOLUME /conf