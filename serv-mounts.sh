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
    "accommodation-serv.yml"
    "auth-serv.yml"
    "mail-serv.yml"
    "metrics-command-serv.yml"
    "metrics-query-serv.yml"
    "notification-serv.yml"
    "recommendation-serv.yml"
    "reservation-serv.yml"
    "user-serv.yml"
)

# Apply each YAML file one by one
for file in "${files[@]}"
do
    kubectl apply -f "$file"
done
