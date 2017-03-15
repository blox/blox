# Changelog

## 0.3.0

### New features
* Added transactions to the daemon-scheduler for all database writes. [#169](https://github.com/blox/blox/pull/169)
* Automated the cluster-state-service e2e tests setup. [#178](https://github.com/blox/blox/pull/178)
* Automated the daemon-scheduler e2e and integration tests setup. [#186](https://github.com/blox/blox/pull/186)

## 0.2.1

### Notable bug fixes
* Fixed the cluster-state-service information from getting stale after network interruption. [#179](https://github.com/blox/blox/pull/179)
* Fixed the cluster-state-service reconciler to not immediately delete stopped tasks. [#167](https://github.com/blox/blox/pull/167)
* Fixed the cluster-state-service reconciler to replace outdated tasks and container instances. [#171](https://github.com/blox/blox/pull/171)

## 0.2.0

### New features
* Added continuous integration with Travis CI. [#150](https://github.com/blox/blox/pull/150)
* Added Kinesis as an event stream consumer type to the cluster-state-service. [#77](https://github.com/blox/blox/pull/77)
* Added support to the cluster-state-service API for combining multiple task filters. [#138](https://github.com/blox/blox/pull/138)
* Added support to the cluster-state-service API for filtering tasks by startedBy. [#122](https://github.com/blox/blox/pull/122)
* Added support to the daemon-scheduler API for filtering environments by cluster ARN. [#87](https://github.com/blox/blox/pull/87)
* Added versioned streaming to the cluster-state-service API. [#143](https://github.com/blox/blox/pull/143)
* Modified the cluster-state-service v1 API response objects and URIs. These changes to the API are not backwards compatible with Blox v0.1.0. [#93](https://github.com/blox/blox/pull/93),[#143](https://github.com/blox/blox/pull/143)
* Refactored the cluster-state-service and daemon-scheduler swagger artifact generation to be consistent. [#119](https://github.com/blox/blox/pull/119)
* Updated the daemon-scheduler to make pending deployments asynchronous. [#145](https://github.com/blox/blox/pull/145)

### Notable bug fixes
* Fixed the CloudFormation template to dynamically select availability zone. [#94](https://github.com/blox/blox/pull/94)
* Fixed the daemon-scheduler environment APIs to return the TaskDefinition field correctly. [#68](https://github.com/blox/blox/pull/68)
* Fixed the null resource information in the cluster-state-service instance API responses. [#101](https://github.com/blox/blox/pull/101)

## 0.1.0

* Initial release of cluster-state-service and daemon-scheduler.
