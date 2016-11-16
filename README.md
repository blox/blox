# Amazon EC2 Container Service Event Stream Handler

The Amazon EC2 Container Service Event Stream Handler is a software developed
to consume events from the Amazon EC2 Container Service Event Stream and
provide a local view of the cluster state.

## Usage

### Setting up
TODO. Add a link to the getting started docs, using CFN

### Running the event stream handler
The following command prints the usage of the Amazon EC2 Container Service
event stream handler:
```bash
$ docker run --rm amazon/amazon-ecs-event-stream-handler --help
amazon-ecs-event-stream-handler handles amazon ecs event stream. It
processes EC2 Container Service events and creates a localized data store, which
provides you a near-real-time view of your cluster state.

Usage:
  amazon-ecs-event-stream-handler [flags]

Flags:
      --etcd-endpoint stringArray   Etcd node addresses
      --queue string                SQS queue name
```

You can also override the logger configuration like the log file and log lever
and AWS CLI parameters like the region and profile when you run the event stream
handler. For example, if running on a local desktop:
```bash
$ docker run -e AWS_REGION=us-east-1 \
    AWS_PROFILE=esh-test \
    CSS_LOG_FILE=/var/output/logs/css.log \
    CSS_LOG_LEVEL=info \
    -v ~/.aws:/.aws \
    -v /tmp/css-logs:/var/output/logs \
    amazon/amazon-ecs-event-stream-handler --queue event_stream
```
