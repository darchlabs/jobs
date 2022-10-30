# Jobs

## Reviewers

- github.com/cagodoy
- github.com/mtavano
- github.com/nicobarra

## Definition

Jobs is the generic module that contains different implementations for smart contracts off-chain interactions.

### Context

There are multiple options for jobs (or keepers) in order to manage off-chain smart contract interactions nowadays. Some of them are very expensive, or some have a lot of friction when trying to setup and configure it. Also, there is a lack of well-delivered information about the jobs state and the interactions.

### Why Do We Do This?

The purpose of this solution is to provide different implementations for the user to choose which he prefers the most in an easy way. Also, make frictionless the setup and funding process.

In addition, it also provides a dashboard with-well delivered information logs and metrics about the services that the user is consuming.

Finally, inside this module could be integrated a cloud based solution for the user that cannot afford the expensive costs of the keepers today.

### Diagrams

- Architecture diagram:

https://app.diagrams.net/#G1PxvFkkQKAgMXKkp0dnIGnzweYV3uBXMT

### Proposed Solution

In the first MVP that will be built by `Darch Labs`, the solution is implemented in the `Golang` programming language.
The code can be divided in different parts managing different tasks that will be connected between them:

- Interfaces definition for implementations, and integrations for them

- An API for the user (from the fronted) to interact with one of the module actions

- Read and Write in a DB the Jobs providers implemented available, the smart contracts being used by the user

- Write in a DB the state of the jobs providers being used, just like using the `synchronizer V2` for getting the smart contracts interactions and writing it in the DB

- A dashboard that will the DB data in the frontend
