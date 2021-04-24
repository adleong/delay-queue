# delay-queue

A prototype for experimenting with RabbitMQ's delayed message exchange.  A Go server accepts requests to a `/send`
endpoint and writes a message to the delayed message queue.  It also consumes messages from the same queue and
prints the messages when it receives them.  This demonstrates the delay between writing a message to the queue
and consuming it from the queue.

The `/send` endpoint requires a `msg` query param for the body of the message and a `delay` query params for
the delay in milliseconds.

## Install RabbitMQ

1. Install the RabbitMQ kubectl plugin: https://www.rabbitmq.com/kubernetes/operator/kubectl-plugin.html
1. Install the RabbitMQ cluster operator: `kubectl rabbitmq install-cluster-operator`
1. Build RabbitMQ docker image which includes the delayed message plugin: `docker build -t rabbitmq:3.8.8-delay .`
  1. (push/load docker image to Kubernetes cluster as necessary)
1. Create RabbitMQ cluster named hello: `kubectl apply -f hello.yml` 

## Build and deploy messages appliction

1. `cd messages`
1. Build messages application `docker build -t ghcr.io/adleong/messages:0.0.1 .`
  1. (push/load docker image to Kubernetes cluster as necessary)
1. Deply messages application: `kubectl apply -f deployment.yml`

## Send messages

1. Create a port-forward: `kubectl port-forward deploy/messages 8888`
1. Send a message `curl 'localhost:8888/send?msg=hi&delay=5000'`
1. Watch the logs `kubectl logs -f deploy/messages`
