# testapp-ThreePilars
I use this to demonstrate the three pilars of observability 

The following repo contains a very basic webserver: 
- [testapp-ping](https://github.com/coffeegoesincodecomesout/testapp-ping)

The following repo contains the same webserver with basic logging:
- [testapp-pingLogging](https://github.com/coffeegoesincodecomesout/testapp-pingLogging)

The following repo contains the same webserver with basic Prometheus instrumentation: 
- [testapp-InstrumentedPing](https://github.com/coffeegoesincodecomesout/testapp-InstrumentedPing)

The following repo contains the same webserver instrumentated for Opentelemetry tracing: 
- [testapp-OTELbasic](https://github.com/coffeegoesincodecomesout/testapp-OTELbasic)

This repo contains that same webserver with metrics, logs and traces enabled - The three pilars of observability.  

### Instructions 

1. Ensure the Red Hat Build of OpenTelemetry operator is installed on your cluster 

2. Create namespace, deployment, service, route, and OTEL collector in sidecar mode

```
oc apply -f manifests/ 
```

3. scale the testapp down and back up, inorder to deploy the OTEL sidecar

```
oc scale --replicas=0 deployment/threepilar-example-deployment
oc scale --replicas=1 deployment/threepilar-example-deployment
```

4. call the endpoint, check the log, check the metrics, view the trace

```
curl -I `oc get route threepilar-example-route -n ns1 | awk 'NR>1 {print $2}'`/ping
oc logs -n ns1 deployment/threepilar-example-deployment -c threepilar-example
oc -n ns1 exec deployment/threepilar-example-deployment -c threepilar-example -- curl -s localhost:8090/metrics
oc logs -n ns1 deployment/threepilar-example-deployment -c otc-container
```
