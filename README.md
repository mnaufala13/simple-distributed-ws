# Chat Example

This application fork from https://github.com/gorilla/websocket to implement a simple distributed web chat application.

## Running the example

To running this example, you should prepare kubernetes cluster. You can use minikube or kubernetes docker for mac.

    $ helm install redis bitnami/redis
    $ helm install ingress-nginx ingress-nginx/ingress-nginx -f ingress.yaml
    $ kubectl apply -f deployment.yaml

To use the chat example, open http://your_server:8080/ in your browser.

