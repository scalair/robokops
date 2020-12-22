FROM scalair/robokops-base:0.5.4
COPY src /home/builder/src
RUN sudo chown -R builder:builder /home/builder
RUN echo 'aws-alb-ingress-controller' | sudo tee /name
VOLUME /conf
