# Welcome Note

Welcome to YouShallNotPass and thank you for contributing your time and expertise to the project.  This document describes the contribution guidelines for the project

* [Setup](#setup)
    * [Environment Setup](#environment-setup)
    * [Testing Setup](#testing-setup)
* [Contributing Steps](#contributing-steps)
* [Contributing Checks](#contributing-checks)
* [Conventions](#conventions)
* [PR Process](#pr-process)
* [Before creating a PR](#what-to-do-before-submitting-a-pull-request)
* [Adding New Checks](#adding-new-checks)


## Setup


### Environment Setup

You must install these tools:

1.  [`git`](https://help.github.com/articles/set-up-git/): For source control
2.  [`go`](https://golang.org/doc/install): You need go version
    [v1.20](https://golang.org/dl/) or higher.
3.  [`docker`](https://docs.docker.com/engine/install/): `v18.9` or higher.


### Testing Setup

A testing file [testing/test.sh](testing/test.sh) is used make the running of tests easier.

For unit testing run:
```sh
testing/test.sh unit
```

Integration testing sets up set of docker images to perform testing on.

To build the docker images run:
```sh
testing/test.sh integration build
```

To run the integration tests run:
```sh
testing/test.sh integration
```


## Contributing Steps

1. Submit an issue describing your proposed change to the repo.
2. For the repo, develop, and test your code changes
3. Submit a pull request


## Contributing Checks

To Contribute Checks to YouShallNotPass follow the following steps:

1. Create an issue describing the new check to be added (with the enhancement label)
2. Create a Fork of the YouShallNotPass repository
3. Create a new module in pkg/checks with the name of the check to be created
    - Please name the new module the check name in lowercase with no separations.
    - Example, Script Hash Check -> pkg/checks/scripthash.
4. Create a struct and initializer to pass in the check configuration and other information (besides the whitelist) that is necessary for the check's completion
    - Example, the Script Hash Check -> `func NewScriptHashCheck(config config.CheckConfig, jobName string, scriptLines []string) ScriptHashCheck`.
5. Implement the IsValidForCheckType function from the Check interface given [here](pkg/checks/checks.go)
    - This function should return true if the checkType is valid for the given check (valid checkTypes are ImageCheck, ScriptCheck, and All).
    - Example, the Script Hash Check is only valid for Script Checks and All Checks -> returns true when `checkType == ScriptCheck` OR `checkType == All`.
6. Implement the IsValidForPlatform function from the Check interface given [here](pkg/checks/checks.go)
    - This function should return true if the CI Platform is valid for a given check (valid ciPlatforms are the lowercase of a given platform - i.e. gitlab, github, ...).
    - For example, a Gitlab Linter Check would only be valid for the "gitlab" platform -> return true if `ciPlatform == "gitlab"`
7. Implement the Check function from the Check interface given [here](pkg/checks/checks.go)
    - This is where the magic of your check happens.  In this function all of the necessary logic that must occur to check whether or not your check has been passed happens.
    - All checks are given access to the whitelist in addition to any information they have stored as a part of their initializer.
    - To make things slightly easier every Check implementation should start with `defer wg.Done()`.  This lets the wait group the check is part of know that the check has completed its necessary tasks and the main thread can read the results of each check.
    - After a check has come to a conclusion as to the success of its execution it will create a [checkResult](pkg/checks/checks.go) and send it through the channel provided through the function.  This CheckResult is then read by the main thread after every check is completed to create a "scorecard" for the YouShallNotPass execution.
8. Create a second file title `{checkname}_test.go` and create a series of unit tests to ensure that the check you created does exactly what you expect it to do.
9. Add the check to the [checkparser](pkg/checkparser/checkparser.go) and [checkparser unit tests](pkg/checkparser/checkparser_test.go).
10. Add the new check to the [unit testing bash script](testing/unit/test.sh)
    - more or less this is basically running go test with the path to the new check `go test "${currentDir}/../../pkg/checks/newcheckname"`
11. Add the check (and configuration options) to the [README](README.md#project-configuration-options)
12. Create a PR documenting the new check and linking it back to the Issue created (or found) in step 1.


## Conventions

* modules should be all lowercase characters
* variables should be named using mixedCase


## PR Process

Every PR should be annotated with an icon indicating whether it's a:
* Breaking change: (:warning:)
* Non-breaking feature: (:sparkles:)
* Patch fix: (:bug:)
* Docs: (:book:)
* Tests/Other: (:seedling:)
* No release note: (:ghost:)

Only the final PR needs to be marked with its icon.


## What to do before submitting a Pull Request

The following tests can be run to test your local changes.

| Command  | Description                                        | Is called in the CI? |
| -------- | -------------------------------------------------- | -------------------- |
| testing/test.sh unit | Runs go test                           | yes                  |
| testing/test.sh integration | Runs integration tests          | partialy             |


## Adding New Tests

When adding new tests please follow the guidelines in [TESTING.md](TESTING.md)