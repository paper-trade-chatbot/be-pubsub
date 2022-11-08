#!/bin/bash

# Bail on error.
# set -eu
# set -o pipefail

RED="\033[0;31m"
GREEN="\033[0;32m"
YELLOW="\033[1;33m"
NC="\033[0m"

#====================================================================
#|               Please put your topic here                         |
#====================================================================
TOPIC_NAME=private.compute.compress.json.0
#====================================================================
#|               Please put your topic here                         |
#====================================================================

RETENTION_MS=600000
PARTITIONS=3


printf "${GREEN}Kafka:${NC}topic generate start... [${TOPIC_NAME}]\n"

TOPICS=$(kubectl exec --tty -i kafka-0 -- /opt/bitnami/kafka/bin/kafka-topics.sh --list --bootstrap-server kafka:9092)

COMPARE_RESULT=$(echo "${TOPICS}" | grep -w "${TOPIC_NAME}" )

if [ -z "$COMPARE_RESULT" ]; 
then
    printf "${GREEN}Kafka:${NC}topic [${TOPIC_NAME}] not created yet. start creating topic...\n"
    printf "${GREEN}Kafka:${NC}"
    kubectl exec --tty -i kafka-0 -- /opt/bitnami/kafka/bin/kafka-topics.sh --create --bootstrap-server kafka:9092 \
        --topic ${TOPIC_NAME} \
        --partitions ${PARTITIONS} \
        --replication-factor 1 \
        --config cleanup.policy=delete \
        --config retention.ms=${RETENTION_MS}

    printf "${GREEN}Kafka:${NC}topic generate successed!\n"
else
    printf "${YELLOW}Kafka:${NC}topic already exist. skip...\n"
fi
