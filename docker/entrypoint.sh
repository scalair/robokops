#!/bin/bash
set -e

if [ ! -f /name ]; then
	echo "The container must define its name in /name"
	exit 1
fi

FEATURE_NAME=$(cat /name)
ACTION=$1


# Load custom configuration and authenticate to the cluster
echo "Load custom configuration and authenticate to the cluster" | boxes -d shell -p l4r4
if [ -f /conf/common.conf ]; then
	source /conf/common.conf
fi
if [ -f /conf/${FEATURE_NAME}/${FEATURE_NAME}.conf ]; then
	source /conf/${FEATURE_NAME}/${FEATURE_NAME}.conf
fi
if [ -f /conf/cluster-login.sh ]; then
	/conf/cluster-login.sh
fi

# Customise manifests
if [ -d /conf/${FEATURE_NAME} ]; then
	echo "Updating manifests with custom configurations" | boxes -d shell -p l4r4
	# Since /conf is a docker volume, we have to copy everything before making changes
	cp -r /conf/${FEATURE_NAME} /tmp
	for MANIFEST in $(find /tmp/${FEATURE_NAME}/ -name '*.yaml' | sed "s|/tmp/${FEATURE_NAME}/||"); do
		echo -e "Templating:\t /tmp/${FEATURE_NAME}/${MANIFEST}"
		envsubst < /tmp/${FEATURE_NAME}/${MANIFEST} > /tmp/${FEATURE_NAME}/${MANIFEST}.tmp
		mv /tmp/${FEATURE_NAME}/${MANIFEST}.tmp /tmp/${FEATURE_NAME}/${MANIFEST}
		if [ -f /home/builder/src/${MANIFEST} ]; then
			echo -e "Merging:\t /tmp/${FEATURE_NAME}/${MANIFEST} with /home/builder/src/${MANIFEST}"
			yq m -x -i /home/builder/src/${MANIFEST} /tmp/${FEATURE_NAME}/${MANIFEST}
		else
			echo -e "Copying:\t /tmp/${FEATURE_NAME}/${MANIFEST} into /home/builder/src/${MANIFEST}"
			cp /tmp/${FEATURE_NAME}/${MANIFEST} /home/builder/src/${MANIFEST}
		fi
	done
fi

# Dry run mode
# 	generate all manifests and copy them to /local which is mount to /tmp
if [ "$ACTION" = "dry-run" ]; then
	# Some features need to generate manifests in a different way using build.sh
	if [ -f /home/builder/src/build.sh ]; then
		echo "Run build.sh" | boxes -d shell -p l4r4
		/home/builder/src/build.sh
	fi
	DIR_NAME=$(date +%Y%m%d%H%M%S)
	mkdir -p /local/robokops/${FEATURE_NAME}/${DIR_NAME}
	cp -r /home/builder/src/* /local/robokops/${FEATURE_NAME}/${DIR_NAME}
	echo "You can find the generated manifests in: /tmp/robokops/${FEATURE_NAME}/${DIR_NAME}" | boxes -d shell -p l4r4
	exit 0
fi

# Helm setup (cluster-init will install tiller in the cluster so we don't do anything here)
if [ ${FEATURE_NAME} != "cluster-init" ]; then
	echo "Configuring Helm" | boxes -d shell -p l4r4
	helm init --client-only
	helm repo update
fi

# Apply pre-manifests if present
if [ -d /tmp/${FEATURE_NAME}/pre-manifests ]; then
	if [ "$ACTION" = "delete" ]; then
		echo "Deleting pre-manifests" | boxes -d shell -p l4r4
		kubectl delete -f /tmp/${FEATURE_NAME}/pre-manifests
	else
		echo "Applying pre-manifests" | boxes -d shell -p l4r4
		kubectl apply -f /tmp/${FEATURE_NAME}/pre-manifests
	fi
fi

# Delete post-manifests if present
if [ -d /tmp/${FEATURE_NAME}/post-manifests ]; then
	if [ "$ACTION" = "delete" ]; then
		echo "Deleting post-manifests" | boxes -d shell -p l4r4
		kubectl delete -f /tmp/${FEATURE_NAME}/post-manifests
	fi
fi

# Run the action scripts
if [ -f /conf/${FEATURE_NAME}/pre-$ACTION.sh ]; then
	echo "Run pre-$ACTION.sh" | boxes -d shell -p l4r4
	/conf/${FEATURE_NAME}/pre-$ACTION.sh
fi
echo "Running $ACTION.sh" | boxes -d shell -p l4r4
/home/builder/src/$ACTION.sh
if [ -f /conf/${FEATURE_NAME}/post-$ACTION.sh ]; then
	echo "Run post-$ACTION.sh" | boxes -d shell -p l4r4
	/conf/${FEATURE_NAME}/post-$ACTION.sh
fi

# Apply post-manifests if present
if [ -d /tmp/${FEATURE_NAME}/post-manifests ]; then
	if [ "$ACTION" != "delete" ]; then
		echo "Applying post-manifests" | boxes -d shell -p l4r4
		kubectl apply -f /tmp/${FEATURE_NAME}/post-manifests
	fi
fi
