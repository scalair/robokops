FROM scalair/robokops-base:0.2.0

# Versions
ENV TERRAFORM_VERSION=0.12.12
ENV TERRAGRUNT_VERSION=0.21.1

WORKDIR /tmp

# Python3 require to run aws-clean-resources.py
RUN sudo apk add --no-cache python3 && \
    if [ ! -e /usr/bin/python ]; then ln -sf python3 /usr/bin/python ; fi && \
    python3 -m ensurepip && \
    sudo rm -r /usr/lib/python*/ensurepip && \
    sudo pip3 install --no-cache --upgrade pip setuptools wheel && \
    if [ ! -e /usr/bin/pip ]; then ln -s pip3 /usr/bin/pip ; fi && \
    sudo pip3 install boto3

# Terraform
RUN wget -q https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip && \
	unzip terraform_${TERRAFORM_VERSION}_linux_amd64.zip && \
	sudo mv terraform /usr/local/bin && \
	rm -f terraform_${TERRAFORM_VERSION}_linux_amd64.zip

# Terragrunt
RUN wget -q https://github.com/gruntwork-io/terragrunt/releases/download/v${TERRAGRUNT_VERSION}/terragrunt_linux_amd64 && \
	sudo mv terragrunt_linux_amd64 /usr/local/bin/terragrunt && \
	sudo chmod +x /usr/local/bin/terragrunt

# Keybase
RUN echo "http://dl-5.alpinelinux.org/alpine/edge/testing" | sudo tee -a /etc/apk/repositories && \
    sudo apk --update add keybase-client

# ssh
RUN sudo apk add openssh && \
    mkdir /home/builder/.ssh

WORKDIR /home/builder/src
COPY src /home/builder/src
RUN sudo chown -R builder:builder /home/builder
RUN echo 'terraform' | sudo tee /name
VOLUME /conf
VOLUME /ssh
