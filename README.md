# This is Paul.
![alt text](docs/images/small-paul.jpg)

---
## Overview
This repo contains everything that makes Paul, Paul!

## TODO
Here is what we still need to do to make Paul a real boy.

- Add security context settings to force the container to run as a non-root user/group and drop all capabilities.
- Tighten up Pauls Kubernetes permissions. We don't think Paul needs *that* much access.
- Implement the interfaces Paul will use to find and query his game servers.
- Teach Paul how to talk and help Paul understand how we're going to talk to him.
- Understand the scope of the second and third item so they can be broken down into actionable subtasks

## Getting Started
Environtment Variables
```
export K8S_NAMESPACE=temporal
export DISCORD_TOKEN=""
```

## This is how we interact with Paul
![alt text](docs/paul.drawio.svg)

#### Label for game server
`gaming.turnbros.app/role=server`

#### Label for type of game server
`gaming.turnbros.app/type=satisfactory`
`gaming.turnbros.app/type=minecraft`
`gaming.turnbros.app/type=rust`
`gaming.turnbros.app/type=avorion`


### Development with Telepresence
Make sure Telepresence is installed
```
brew install datawire/blackbird/telepresence
```

List the services
```
telepresence list -n paul
```

Get the service you wish to highjack
```
kubectl get service -n paul paul --output yaml
```

Intercept the service traffic
```
telepresence intercept -n paul paul --port 8443:http --env-file ./paul-intercept.env
```