FROM scalair/robokops-base:0.5.2
COPY src /home/builder/src
RUN sudo chown -R builder:builder /home/builder
RUN echo 'dashboard' | sudo tee /name
VOLUME /conf