# Schooli API



## Getting Started

These instructions will give you a copy of the project up and running on
your local machine for development and testing purposes. See deployment
for notes on deploying the project on a live system.

### Prerequisites

Requirements for the software and other tools to build, test and push 
- Docker
- Docker-compose
- Migrate
- SQLC
- Air 

### Installing

A step by step series of examples that tell you how to get a development
environment running

Cd into ./build/schooli

    docker-compose up

Make sure all docker containers are running

    docker ps -a 

Change address in .config.toml to wlan0 

    ifconfig 


## Running the project

Simply run air in terminal to run the project.

    air

## For the database

Use makefile command to make migrations

    make migrateup

To make queries

    sqlc generate
