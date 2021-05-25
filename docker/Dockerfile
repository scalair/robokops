FROM golang:alpine3.10

# https://pypi.org/project/awscli/
ENV AWSCLI_VERSION=1.17.0
# https://docs.aws.amazon.com/eks/latest/userguide/install-aws-iam-authenticator.html
ENV AWS_IAM_AUTHENTICATOR=1.13.7
# https://github.com/GoogleCloudPlatform/cloud-sdk-docker
ENV GCLOUD_SDK_VERSION=281.0.0
# https://aur.archlinux.org/packages/kubectl-bin/
ENV KUBECTL_VERSION=1.15.0
# https://github.com/kubernetes/helm/releases
ENV HELM_VERSION=2.17.0
# https://releases.hashicorp.com/vault/
ENV VAULT_VERSION=0.11.6
# https://github.com/krallin/tini
ENV TINI_VERSION v0.18.0

# sudo
RUN apk --no-cache add sudo

# Create user builder
RUN addgroup -S builder && adduser --disabled-password -S builder -G builder
RUN echo "builder ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers
USER builder

# Set timezone to UTC by default
RUN sudo ln -sf /usr/share/zoneinfo/Etc/UTC /etc/localtime

# Install base and dev packages
RUN sudo apk update && \
	sudo apk add bash curl ca-certificates unzip git expect tini apache2-utils bison flex gettext

# aws-cli
RUN sudo apk -Uuv add groff less python2 python2-dev py-pip build-base && \
	sudo pip install awscli==$AWSCLI_VERSION && \
	sudo apk --purge -v del py-pip && \
	sudo rm -rf /var/cache/apk/*

# aws-iam-authenticator
RUN sudo wget -q https://amazon-eks.s3-us-west-2.amazonaws.com/${AWS_IAM_AUTHENTICATOR}/2019-06-11/bin/linux/amd64/aws-iam-authenticator -O /usr/local/bin/aws-iam-authenticator && \
	sudo chmod +x /usr/local/bin/aws-iam-authenticator

# gcloud
ENV PATH /google-cloud-sdk/bin:$PATH
RUN curl -O https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-sdk-${GCLOUD_SDK_VERSION}-linux-x86_64.tar.gz && \
	tar xzf google-cloud-sdk-${GCLOUD_SDK_VERSION}-linux-x86_64.tar.gz && \
	sudo mv google-cloud-sdk / && \
	rm google-cloud-sdk-${GCLOUD_SDK_VERSION}-linux-x86_64.tar.gz && \
	gcloud config set core/disable_usage_reporting true && \
	gcloud config set component_manager/disable_update_check true && \
	gcloud config set metrics/environment github_docker_image

# kubectl
RUN sudo wget -q https://storage.googleapis.com/kubernetes-release/release/v${KUBECTL_VERSION}/bin/linux/amd64/kubectl -O /usr/local/bin/kubectl && \
	sudo chmod +x /usr/local/bin/kubectl

# helm
RUN sudo wget -q https://storage.googleapis.com/kubernetes-helm/helm-v${HELM_VERSION}-linux-amd64.tar.gz && \
	sudo tar -xf helm-v${HELM_VERSION}-linux-amd64.tar.gz && \
	sudo mv linux-amd64/helm /usr/local/bin/helm && \
	sudo rm -rf helm-v${HELM_VERSION}-linux-amd64.tar.gz linux-amd64/ && \
	sudo chmod +x /usr/local/bin/helm

# vault
RUN sudo wget -q https://releases.hashicorp.com/vault/${VAULT_VERSION}/vault_${VAULT_VERSION}_linux_amd64.zip && \
	sudo unzip -q vault_${VAULT_VERSION}_linux_amd64.zip && \
	sudo mv vault /usr/bin/vault && \
	sudo rm -f vault_${VAULT_VERSION}_linux_amd64.zip

# yq
RUN sudo curl -L https://github.com/mikefarah/yq/releases/download/2.4.1/yq_linux_amd64 -o /usr/bin/yq && \
	sudo chmod +x /usr/bin/yq

# boxes (https://boxes.thomasjensen.com/)
RUN sudo wget -q https://github.com/ascii-boxes/boxes/archive/v1.3.zip -O boxes-1.3.zip && \
	sudo unzip -q boxes-1.3.zip && \
	sudo make -C boxes-1.3 -s 2>/dev/null 1>/dev/null || true && \
	sudo cp boxes-1.3/src/boxes /usr/local/bin && \
	sudo cp boxes-1.3/boxes-config /usr/share/boxes && \
	sudo rm -rf boxes-1.3 boxes-1.3.zip

# cleanup
RUN sudo apk del unzip bison flex

WORKDIR /home/builder/src
COPY entrypoint.sh /

ENTRYPOINT ["/sbin/tini", "-g", "--", "/entrypoint.sh"]
