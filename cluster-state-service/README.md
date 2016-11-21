# Cluster State Service 

The Cluster State Service is a software developed to consume events 
from the Amazon EC2 Container Service Event Stream and provide a 
local view of the cluster state.

## Usage

### Setting up
TODO. Add a link to the getting started docs, using CFN

### Running the event stream handler
The following command prints the usage of the Cluster State Service:
```bash
$ docker run --rm amazon/blox-cluster-state-service --help
cluster-state-service processes EC2 Container Service events and  creates
a localized data store, which provides you a near-real-time view of your cluster state.

Usage:
  cluster-state-service [flags]

Flags:
      --bind string                 Cluster State Service listen address
      --etcd-endpoint stringArray   Etcd node addresses
      --queue string                SQS queue name
```

You can also override the logger configuration like the log file and log lever
and AWS CLI parameters like the region and profile when you run the event stream
handler. For example, if running on a local desktop:
```bash
$ docker run -e AWS_REGION=us-east-1 \
    AWS_PROFILE=css-test \
    CSS_LOG_FILE=/var/output/logs/css.log \
    CSS_LOG_LEVEL=info \
    -v ~/.aws:/.aws \
    -v /tmp/css-logs:/var/output/logs \
    amazon/blox-cluster-state-service --queue event_stream
```
