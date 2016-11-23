# Updating to a new version
To update the cluster-state-service to a new version, run the `update-version` script 
in the `scripts` directory. 

# Usage
From the `cluster-state-service`, run the `update-version` script with the version 
number as the argument. Example: 
```
$ ./scripts/update-version 0.1.0
```

This should update/create the `versioning/version.go` file with the appropriate 
version number.
