#!/usr/bin/env python3
import argparse
import boto3
import time

'''
This script is a workaround to terraform issues and limitations. On of them is regarding the ASG instance protection:
https://github.com/terraform-providers/terraform-provider-aws/issues/5278

What it's doing:
- Remove protectedFromScaleIn protection of the ASG to prevent new protected instances to be created
- Terminate all instances matching the tag Env:[env_name] (which have been created with the protected flag)
- Delete network interfaces that were attached to instances because some of them never get deleted otherwise
- Delete network interfaces that were attached to the worker nodes security group. This should be done when delete NI of instances, but sometimes it doesn't
'''
ec2 = None
ec2_client = None
asg_client = None

def parse_args():
	description = 'AWS Clean Resources'
	parser = argparse.ArgumentParser()
	parser.add_argument('-e', '--env_name',
		help='Environment name tag',
		required=True,
		action='store')

	args = parser.parse_args()
	return args

def main():
	global ec2
	global ec2_client
	global asg_client

	args      = parse_args()
	env_name  = args.env_name
	region = env_name.split('.')[3]

	ec2        = boto3.resource('ec2', region_name=region)
	ec2_client = boto3.client('ec2', region_name=region)
	asg_client = boto3.client('autoscaling', region_name=region)

	asgs = get_asg(env_name)
	remove_asg_protection(asgs)

	instances = get_instances(env_name)
	instances_network_interface_ids = get_instances_network_interfaces(instances)
	terminate_instances(instances)
	delete_network_interfaces(instances_network_interface_ids)

	worker_sg = get_worker_sg(env_name)
	sg_network_interface_ids = get_sg_network_interfaces(worker_sg)
	delete_network_interfaces(sg_network_interface_ids)

	print("terminate_instances finished successfully")


# ASG
def get_asg(env_name):
    paginator = asg_client.get_paginator('describe_auto_scaling_groups')
    page_iterator = paginator.paginate(
        PaginationConfig={'PageSize': 100}
    )
    asgs = page_iterator.search(
        'AutoScalingGroups[] | [?contains(Tags[?Key==`Env`].Value, `'+env_name+'`)]'.format(
            'Application', 'CCP')
    )
    return asgs

def remove_asg_protection(asgs):
    for asg in asgs:
        try:
            asg_client.update_auto_scaling_group(AutoScalingGroupName=asg['AutoScalingGroupName'], NewInstancesProtectedFromScaleIn=False)
            print("Remove protected from scale in of "+asg['AutoScalingGroupName'])
        except Exception as err:
            print("Failed to suspend processes: %s\n" % err)
            return False

# Instances
def get_instances(env_name):
	try:
		instances = ec2_client.describe_instances(Filters=[{'Name': 'tag:Env', 'Values': [env_name]}])
	except Exception as err:
		print("Failed to describe instances: %s\n" % err)
		return False
	return instances

def terminate_instances(instances):
	for reservation in instances['Reservations']:
		for instance in reservation['Instances']:
			try:
				if instance['State']['Name'] in 'running':
					instance_id = instance['InstanceId']
					print("Terminating %s ... " % instance_id)
					ec2_client.terminate_instances(InstanceIds=[instance_id])
					wait_terminated_instances(instance_id)
			except Exception as err:
				print("Failed to terminate %s" % instance['InstanceId'])
				continue

def wait_terminated_instances(instance_id):
	instances = ec2_client.describe_instances(Filters=[{'Name': 'instance-id', 'Values': [instance_id]}])
	for reservation in instances['Reservations']:
		for instance in reservation['Instances']:
			if instance['State']['Name'] not in 'terminated':
				print("Waiting for instance %s to terminate. Current status is %s" % (instance_id, instance['State']['Name']))
				time.sleep(10)
				wait_terminated_instances(instance_id)
	return True

# Security groups
def get_worker_sg(env_name):
	filter_tag = "tag:kubernetes.io/cluster/" + env_name.replace('.', '-')
	try:
		worker_sgs = ec2_client.describe_security_groups(Filters=[{'Name': 'tag:Env', 'Values': [env_name]}, {'Name': filter_tag, 'Values': ['owned']}])
	except Exception as err:
		print("Failed to describe security groups: %s\n" % err)
		return False
	worker_sg = ""
	if worker_sgs["SecurityGroups"][0]["GroupId"]:
		worker_sg = worker_sgs["SecurityGroups"][0]["GroupId"]
	return worker_sg

def get_sg_network_interfaces(sg_id):
	network_interface_ids = []
	try:
		network_interfaces = ec2_client.describe_network_interfaces(Filters=[{'Name': 'group-id', 'Values': [sg_id]}])
	except Exception as err:
		print("Failed to describe security group network interface: %s\n" % err)
		return False
	for network_interface in network_interfaces["NetworkInterfaces"]:
		network_interface_ids.append(network_interface["NetworkInterfaceId"])
	return network_interface_ids

# Network Interfaces
def get_instances_network_interfaces(instances):
	network_interface_ids = []
	for reservation in instances['Reservations']:
		for instance in reservation['Instances']:
			for network_interface in instance['NetworkInterfaces']:
				network_interface_ids.append(network_interface['NetworkInterfaceId'])
	return network_interface_ids

def delete_network_interfaces(network_interface_ids):
	for network_interface_id in network_interface_ids:
		network_interface = ec2.NetworkInterface(network_interface_id)
		# first try to detach it from the instance (in case the instance hasn't been terminated completely yet)
		try:
			print("Detaching %s ... " % network_interface_id)
			network_interface.detach()
			# give it a bit of time...
			time.sleep(45)
		except Exception as err:
			print("Fail to detach %s. Skipping" % network_interface_id)
		# then try to remove it if it still exists
		try:
			print("Deleting %s ... " % network_interface_id)
			network_interface.delete()
		except Exception as err:
			print("Fail to delete %s. Skipping" % network_interface_id)
			continue


if __name__ == "__main__":
    main()
