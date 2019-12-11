FROM scalair/robokops-base:0.2.3
COPY src /home/builder/src
RUN sudo chown -R builder:builder /home/builder
RUN echo 'gitlabci' | sudo tee /name
VOLUME /conf