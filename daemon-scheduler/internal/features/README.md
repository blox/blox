# Setup
- Have [etcd](https://github.com/coreos/etcd) running locally.
- AWS Credentials available using default providers (environment variables, instance-profile)
- AWS Region appropriately configured (e.g. AWS_REGION=us-west-2)
- Start blox-daemon-scheduler (blox-daemon-scheduler --bind localhost:2000 --etcd-endpoint localhost:2379)
- Create ECS cluster (e.g. name=q-daemon-scheduler-test)
- Create AutoScaling group (e.g. name=ecs-daemon-scheduler-test-asg) which can launch instances into the above ECS cluster. Follow the steps [here](http://docs.aws.amazon.com/AmazonECS/latest/developerguide/launch_container_instance.html)

# Usage
- In the package root directory (parent of internal)
    ```
    ECS_CLUSTER=... ECS_CLUSTER_ASG=... AWS_REGION=... AWS_PROFILE=... gucumber -tags=@e2e
    ```

    ```
    ECS_CLUSTER=... ECS_CLUSTER_ASG=... AWS_REGION=... AWS_PROFILE=... gucumber -tags=@integ
    ```
