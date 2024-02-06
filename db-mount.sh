#!/bin/bash

# Function to handle interrupt signals
interrupt_handler() {
    echo "Script interrupted. Exiting..."
    exit 1
}

# Trap interrupt signals and call the interrupt_handler function
trap interrupt_handler SIGINT

# Define the Kubernetes YAML files
files=(
    "user-db.yml"
    "reservation-db.yml"
    "redis-db.yml"
    "notification-db.yml"
    "neo4j.yml"
    "nats.yml"
    "namenode-db.yml"
    "datanode-db.yml"
    "esdb.yml"
    "auth-db.yml"
    "accommodation-db.yml"
)

# Apply each YAML file one by one
for file in "${files[@]}"
do
    kubectl apply -f "$file"
done
