#!/bin/bash

# All parameter fields are required for the script to execute
#
# parameter ${1} is autodeploy
# parameter ${2} is the openshift-master-url
# parameter ${3} is the user-token
#
# parameters 4 & 5 are only needed if ${1} == "true"
# parameter ${4} is the project/namespace in openshift where the deploymentconfig resides
# parameter ${5} is the name of the deploymentconfig for rollout of the new image
#
# example usage : build-deploy.sh true https://master.balt1.okd.14west.io:8443 EI_L1xrzGnOdsJwoCfJ3PrkDNc9OLxliYE4hGuTKG10 test poc-microservice

# declare some variables
name="kaniko-simple-microservice"

if [ "$#" -lt 3 ]
then
  echo " "
  echo -e "\033[0;93musage  build-deploy.sh <autodeploy> <master-url> <user-token> <namespace> <deploymentconfig>\033[0m"
  echo " "
  exit -1
fi

# list some gocd variables
echo "GOCD job name         ${GO_JOB_NAME}"
echo "GOCD pipeline name    ${GO_PIPELINE_NAME}"
echo "GOCD pipeline counter ${GO_PIPELINE_COUNTER}"
echo "GOCD pipeline label   ${GO_PIPELINE_LABEL}"
echo "GOCD to revision      ${GO_TO_REVISION}"
echo "GOCD from revision    ${GO_FROM_REVISION}"


# first login
oc login ${2} -n ci-cd --token=${3} --insecure-skip-tls-verify


# we can now execute the job
oc create -f kaniko-job.yml

status=""
while [ "${status}" == "" ]
do
  status=$(oc get job/${name} -o=jsonpath='{.status.conditions[*].type}')
done

pod=$(oc get pods | grep "${name}" | awk '{print $1}')
oc logs po/"${pod}"

if [ "${status}" != "Complete" ];
then
  echo "Failed"
 	exit -1
else
  echo "Passed"
  # if we aren't deploying then just exit
  if [ "${1}" == "false" ];
  then
    oc delete job/"${name}"
    exit 0
  fi
fi

# delete the job
oc delete job/"${name}" 

if [ "${3}" == "true" ];
then
  # we assume that the project resides on the same server (master-url)
  # if not then add a new login call here first
  oc project ${4}
  oc rollout ${5}
  exit 0
fi
