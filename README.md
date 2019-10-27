# WOW Designs Front-End API
This is the api used by the front end of wow designs. 

## Overview
The api is written in Go and uses mongoDB for a databse and GraphQL as an interface. The folder structure is hopefully pretty clear but essentially the queries are in the queries folder, the mutations in the mutations folder and so on. 

## Building the code
To build the API run ```make build``` which should build an executable at /bin/api which you can use to run the api.

## Docker
I will soon add a Dockerfile to deploy this application - WIP

## Testing
Testing isn't really done at this point. Mostly to cut down on development time although skipping it at this stage can be quite dangerous. Rigorous testing will follow afterwards and you should be able to run the tests with ```make tests``` - WIP