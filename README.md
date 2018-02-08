A [Concourse](http://concourse.ci/) [resource](http://concourse.ci/resources.html) to interact with the build status API of [Atlassian BitBucket](https://bitbucket.org).

This repo is tied to the [associated Docker image](quay.io/pickledrick/concourse-bitbucket-pullrequest-resource) on quay.io, built from the master branch.
## Resource Configuration


These items go in the `source` fields of the resource type. Bold items are required:
 * **`repo`** - repository name to track
 * **`key`** - OAuth key for Consumer
 * **`secret`** - OAuth Secret for Consumer
 * **`team`** - Team name repository belongs to
 * **`url`** - bitbucket cloud api path (example: `https://api.bitbucket.org`) **Currently only supported**
 * **`version`** - bitbucket API Version (example: `2.0`) **Currently only supported**
 * **`concourse_url`** - concourse url for setting build link in bitbucket (example: `http://ci.example.com`)



## Behavior


### `check`

Checks for a Pull request with a head commit in an untested state.


### `in`

Retrieves a copy of the tracking branch, sets pull request state to IN_PROGRESS.

### `out`

Update the status of a commit.

Parameters:

 * **`commit`** - File containing commit SHA to be updated.
 * **`state`** - the state of the status. Must be one of `success` or `failed`.


## Example

A typical use case is to watch for Pull Requests on a Repository, Run Tests and update the status of a commit.

An example of this is shown in **ci/pipeline.yaml**

## Installation

This resource is not included with the standard Concourse release. Use one of the following methods to make this resource available to your pipelines.


### Deployment-wide

To install on all Concourse workers, update your deployment manifest properties to include a new `groundcrew.resource_types` entry...

    properties:
      groundcrew:
        additional_resource_types:
          - image: "docker:///quay.io/pickledrick/concourse-bitbucket-pullrequest-resource#master"
            type: "pull-request"                   

### Pipeline-specific

To use on a single pipeline, update your pipeline to include a new `resource_types` entry...

    resource_types:
      - name: "pull-request"
        type: "docker-image"
        source:
          repository: "quay.io/pickledrick/concourse-bitbucket-pullrequest-resource"
          tag: "master"


## References

 * [Resources (concourse.ci)](https://concourse.ci/resources.html)
 * [Bitbucket build status API](https://confluence.atlassian.com/bitbucket/use-the-bitbucket-cloud-rest-apis-222724129.html)

## License

[Apache License v2.0]('./LICENSE')
