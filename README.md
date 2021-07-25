This is a lightweight service which exposes an endpoint to scale kubernetes deployment replicas incrementally. 

When a scale is triggered for a deployment, the autoscaler adds/removes a fixed number of replicas as per the deployment's configuration. 
```
target_replicas = current_replicas +/- batch_size
```

## Configuration

A single instance of this application can be used to scale multiple deployments across different namespaces. Every deployment config consists of the following fields:
- `minReplicas`: Minimum number of replicas for the deployment. Autoscaler would never scale below this value. 
- `maxReplicas`: Maximum number of replicas for the deployment. Autoscaler would never scale above this value.
- `scaleUpBatchSize`: Number of replicas to add when `action = "scale-up"`
- `scaleDownBatchSize`: Number of replicas to reduce when `action = "scale-down"`


Please refer to the sample config for reference. 

## Deployment

This is a stateless application and can be deployed as a K8s deployment in the cluster. Please make sure the service account attached to this deployment has sufficient access to trigger scale for the deployments mentioned in the config.yaml

## Integrations

As of now, this is used to trigger scale based on alerts from alertmanager only. Please raise an issue/PR if you wish to add your own integrations.

### AlertManager

This can be used as a webhook reciever in [Alertmanager](https://github.com/prometheus/alertmanager). An alert should contain the following labels for the autoscaler to recognise and scale the deployment. 
- `action` : `scale-up` or `scale-down`. Depends on what the alert is meant for. 
- `namespace` : Namespace in which the target application is deployed. 
- `deployment` : Name of the target deployment in the namespace
