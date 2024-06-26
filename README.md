# bosh-disaster-recovery-acceptance-tests (B-DRATs)

Tests a given BOSH director can be backed up and restored using [`bbr`](https://github.com/cloudfoundry-incubator/bosh-backup-and-restore).

The test runner provides hooks around `bbr director backup` and `bbr director restore`.

## Dependencies

1. Install [Go](https://golang.org/doc/install)
1. Download [`bbr` CLI](https://github.com/cloudfoundry-incubator/bosh-backup-and-restore/releases)

## Running B-DRATs in your pipelines

We encourage you to use our [`run-b-drats` CI task](ci/tasks/run-b-drats) to run B-DRATS in your Concourse pipeline.

Please refer to our b-drats [pipeline definition](ci/pipelines/b-drats/pipeline.yml) for a working example.

## Running B-DRATs locally

1. Clone this repo
    ```bash
    $ git clone git@github.com:cloudfoundry/bosh-disaster-recovery-acceptance-tests.git
    $ cd bosh-disaster-recovery-acceptance-tests
    ```
1. Create an `integration-config.json` file, for example:
    ```json
    {
      "bosh_host": "director-address",
      "bosh_agent_endpoint": "https://mbus:AGENT_PASSWORD@10.0.0.3:6868",
      "bosh_agent_certificate": "bosh-agent-ca",
      "bosh_client": "bosh-client-name",
      "bosh_client_secret": "bosh-client-secret",
      "bosh_ca_cert": "bosh-ca-cert",
      "credhub_client":"credhub-client",
      "credhub_client_secret":"credhub-client-secret",
      "credhub_ca_cert":"credhub-ca-cert",
      "timeout_in_minutes": 30,
      "stemcell_src": "stemcell-path.tgz",
      "include_credhub_testcase": true,
      "include_deployment_testcase": true,
      "include_truncate_db_blobstore_testcase": true
    }
    ```
1. Export `INTEGRATION_CONFIG_PATH` to be path to `integration-config.json` file you just created.
1. Export `BBR_BINARY_PATH` to the path to the BBR binary.
1. Run acceptance tests
    ```bash
    $ ./scripts/_run_acceptance_tests.sh
    ```

## Integration Config Variables

* `bosh_host` - the BOSH director hostname
* `bosh_agent_endpoint` - the BOSH director agent endpoint. Typically https://mbus:AGENT_PASSWORD@DIRECTOR_IP:6868
* `bosh_agent_certificate` - the CA used by the BOSH agent on the BOSH director VM
* `bosh_client` - the BOSH director API client
* `bosh_client_secret` - the BOSH director API client secret 
* `bosh_ca_cert` - the BOSH director API CA certificate
* `credhub_client` - the UAA client that is used to authenticate with CredHub
* `credhub_client_secret` - the UAA secret that is used to authenticate with CredHub
* `credhub_ca_cert` - the mTLS certificate that is used to communicate with CredHub
* `timeout_in_minutes` - default ginkgo `Eventually` timeout in minutes, defaults to `30`
* `stemcell_src` - absolute path to stemcell tarball
* `include_<testcase-name>` - flag for whether to run a given testcase, if omitted defaults to `false`
* `deployment_vm_type` - cloud config vm_type to be passed into config, if omitted defaults to `default` (used by example deployment_testcase)
* `deployment_network` - cloud config network to be passed into config, if omitted defaults to `default` (used by example deployment_testcase)
* `deployment_az` - cloud config az to be passed into config, if omitted defaults to `z1` (used by example deployment_testcase)

## Contributing to B-DRATs

B-DRATS runs a collection of test cases against a bosh director.

Test cases should be used for checking that BOSH director components' data has been backed up and restored correctly. For example, if your release backs up a database during `bbr director backup`, and the database is altered after taking the backup. Then after a successful `bbr director restore`, the content of the database will be restored to its original state.

To add extra test cases, create a new test case that implements the [`TestCase` interface](runner/testcase.go).

The methods that need to be implemented are:
* `Name() string` - should return name of the test case.
* `BeforeBackup(Config)` - runs before the backup is taken, and should create state in the BOSH director to be backed up.
* `AfterBackup(Config)` - runs after the backup is complete but before the restore is started.
* `AfterRestore(Config)` - runs after the restore is complete, and should assert that the state in the restored BOSH director matches that created in `BeforeBackup(Config)`.
* `Cleanup(Config)` - should clean up the state created in the BOSH director through the test.

`Config` contains the config for the BOSH Director and for the CF deployments to backup and restore.

### Creating a new test case

1. Create a new test case in the [testcases package](testcases).
1. Add the newly created test case to `[]runner.TestCase` in [`acceptance_test.go`](acceptance/acceptance_test.go).

