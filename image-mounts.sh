#!/bin/bash

# Function to handle interrupt signals
interrupt_handler() {
    echo "Script interrupted. Exiting..."
    exit 1
}

# Trap interrupt signals and call the interrupt_handler function
trap interrupt_handler SIGINT

# Define the Docker images to load into Minikube
images=(
    "user-service:latest"
#    "notifications-service:latest"
    "accommodations-service:latest"
    "reservations-service:latest"
    "metrics-command:latest"
    "metrics-query:latest"
    "auth-service:latest"
    "mail-service:latest"
    "nats:latest"
    "recommendation-service:latest"
)

# Load each Docker image into Minikube one by one
for image in "${images[@]}"
do
    minikube image load "$image"
done
