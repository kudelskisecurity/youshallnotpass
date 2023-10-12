# Testing


## Unit Tests


### Go Unit Testing

For each file in the Go project, there is (usually) a second file with the name {module}_test.go. These files contain all of the unit tests for their respective main file.

#### Running Unit Tests

To run the all of the unit tests, run testing/test.sh unit. This triggers the bash script found at [testing/unit/test.sh](testing/unit/test.sh) to run `go test` on each of the packages individually. To run a specific package's unit tests run testing/test.sh unit {package_name}.  This will trigger the [testing/unit/test.sh](testing/unit/test.sh) script to run `go test` for the specific package

#### Notes About Unit Tests

The main reason the bash script is used to trigger unit tests is to not trigger the client integration tests when running `go test`. Because these tests are also of the form {package}_test.go, they would also be covered if go test was used to systematically test the entire project.

#### Contributing Unit Tests

For any packages and checks added to the go codebase, a corresponding file should be created following the unit testing file format. This file should contain a set of test cases to test the various functionalities of the Go file/package.

## Integration Tests

### Types of Integration Tests

#### Client Integration Tests

Client integration tests are basically an extra check that the clients, used to contact external services, have the capabilities used by YouShallNotPass.


##### Existing Client Integration Tests

Currently, the following Client Integration tests exist:

1. Mattermost Client Integration Tests
    - The Mattermost Client Integration Tests make sure that the Mattermost client has the capabilities of writing messages to a given Mattermost channel on a test Mattermost Instance (hosted on a docker container).
2. Hashicorp Vault Client Integration Tests
    - The Hashicorp Vault Client Integration Tests make sure that the Hashicorp client can read the YouShallNotPass configuration and whitelist files from a test Hashicorp Vault Instance (hosted on a docker container).


#### Client-Platform-Case Integration (e2e) Tests

Client-Platform-Case Integration (e2e) Tests are basically a set of tests that test YouShallNotPass's performance given a specific vault client (i.e. Hashicorp, ...) paired with a specific CI Platform (i.e. GitHub, GitLab...) given a specific situation. For example, the hashicorp_gitlab_auth_timeout test uses the Hashicorp Vault Client (Client) and the GitLab CI Platform (Platform) to check that the desired behavior happens when YouShallNotPass authentication timeout is reached (Case).


##### Existing Client-Platform-Case Integration Tests

Currently, the following Client-Platform-Case Integration Tests exist:

1. Hashicorp Gitlab Auth Timeout Test
    - The Hashicorp GitLab Auth Timeout Test tests that when YouShallNotPass times out the GitLab executor should not execute the CI/CD Job.
2. Hashicorp GitLab Automatic Test
    - The Hashicorp GitLab Automatic Test tests that when YouShallNotPass's checks all pass the GitLab runner should execute the CI/CD Job automatically.
3. Hashicorp GitLab Bash Test
    - The Hashicorp GitLab Bash Test tests that when a bash script is supplied as part of the GitLab CI/CD Job, YouShallNotPass should automatically expand the bash script before executing script checks.
4. Hashicorp GitLab Fail Test
    - The Hashicorp GitLab Fail Test tests that when checks with abortOnFail=true fail YouShallNotPass should prevent the GitLab CI/CD Job from being run.
5. Hashicorp GitLab Mfa Test
    - The Hashicorp GitLab MFA Test test that when checks with mfaOnFail=true fail YouShallNotPass should require user multi-factor authentication to succeed before GitLab may run the CI/CD Job.


##### How Client-Platform-Case Integration Tests Work

Most of the Client-Platform-Case Integration Tests work with a set of docker-compose files (currently organized into runner-compose.yml for the YouShallNotPass Runner/Executor and vault-compose.yml for the vault setup).

In general, the following steps are performed to complete a Client-Platform-Case Integration Test:

1. Docker-Compose spins up a Hashicorp vault instance
2. The vault instance is initialized by the vault_init service which creates a youshallnotpass-demo role in the vault.
    - This role has access to the necessary youshallnotpass whitelist, config, and scratch vault locations.
3. The vault_init container writes a set of whitelist images and scripts as well as configuration to the vault instance
4. When the vault_init container is done initializing the vault instance, it exits, which is used to spin up the runner-compose.yml runner into the same docker network.
    - The runner-compose.yml file contains the following services:
        - youshallnotpass_builder_daemon: a daemon that continually rebuilds the youshallnotpass binary so the docker image doesn't have to be continually rebuilt.
        - git_repo_init: an alpine image that creates a git repository in the git_repo volume.
        - gitlab_runner: the gitlab runner that will use the youshallnotpass custom executor to execute the jobs in the git_repo volume.
5. The gitlab_runner service uses the custom executor to run a specific job from the git_repo volume.
6. When the gitlab_runner is finished executing the job it will return a status code of 0 on success and 1 on failure. This exit code is then verified to determine whether or not the expected behavior was completed by the YouShallNotPass executor.


##### Contributing Client-Platform-Case Integration Tests

To contribute a Client-Platform-Case Integration Test, the current format for a runner-compose.yml and vault-compose.yml file does not necessarily have to followed if a better workflow is discovered, however it is necessary (for a bit of uniformity) to add the test to [testing/integration/test.sh](testing/integration/test.sh).  In this file tests are run based on their client, platform, and case keywords. For example, all Hashicorp client tests should be run when testing/test.sh integration Hashicorp is run. Therefore, please add your test in a way that it will only be run when the client, platform, and/or case is present in the argument to the `test.sh` file.
