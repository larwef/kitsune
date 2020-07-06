# Kitsune
Simple message ingestion and distribution service.

Under development. Currently only inplements a simple in memory repository.

## Goals
Primarily focusing on simplicity in as many aspects as possible:
- Easy to use as a client.
- Easy to deploy.
- Easy to manage.
- Easy to develop and maintain.

Currently the focus is to optimize for deploying as docker container on a cluster. Think of ECS using a spot fleet in AWS. Which
means an instance can go down and up at any time and you typically run multiple instances spread over several zones. The first
consequence of this is that the storage should be decoupled from the application.

