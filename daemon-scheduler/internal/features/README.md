# Setup
- Have [etcd](https://github.com/coreos/etcd) running locally.
- AWS Profile appropriately configured (e.g. AWS_PROFILE=default)
- AWS Region appropriately configured (e.g. AWS_REGION=us-west-2)
- Start blox cluster state service:
    ```
    AWS_PROFILE=default AWS_REGION=us-west-2 CSS_LOG_FILE=var/output/logs/css.log  ./out/cluster-state-service --queue event_stream --etcd-endpoint localhost:2379 --bind localhost:3000
    ```
- Start blox daemon scheduler:
    ```
    AWS_PROFILE=default AWS_REGION=us-west-2 DS_LOG_FILE=var/output/logs/ds.log DS_LOG_LEVEL=debug ./out/daemon-scheduler --bind localhost:2000 --etcd-endpoint localhost:2379 --css-endpoint localhost:3000
    ```

# Usage
- In the package root directory (parent of internal)
    ```
    AWS_REGION=... AWS_PROFILE=... gucumber -tags=@e2e
    ```

    ```
    AWS_REGION=... AWS_PROFILE=... gucumber -tags=@integ
    ```

# Customization
  Users can pass in their own custom parameters to the tests through environment variables:
  * Cluster (custom name has to be different from default name: "DSTestCluster")
  * Autoscaling Group (custom name has to be different from default name: "DSClusterASG")
  * EC2 key pair

  ```
  ECS_CLUSTER=<custom-cluster> ECS_CLUSTER_ASG=<custom-autoscaling-group> EC2_KEY_PAIR=<ec2-key-pair>
  ```
