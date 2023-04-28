# Contributing to terraform-provider-for-vmware-cloud-director-availability

We welcome contributions from the community and first want to thank you for taking the time to contribute!

Please familiarize yourself with the [Code of Conduct](https://github.com/vmware/.github/blob/main/CODE_OF_CONDUCT.md) before contributing.

Before you start working with terraform-provider-for-vmware-cloud-director-availability, please read our [Developer Certificate of Origin](https://cla.vmware.com/dco). All contributions to this repository must be signed as described on that page. Your signature certifies that you wrote the patch or have the right to pass it on as an open-source patch.

## Community

Terraform Provider for VMware Cloud Director Availability contributors can discuss matters here:
https://vmwarecode.slack.com, channel `#vcda-terraform-dev`

## Ways to contribute

We welcome many different types of contributions and not all of them need a Pull request. Contributions may include:

* New features and proposals
* Documentation
* Bug fixes
* Issue Triage
* Answering questions and giving feedback
* Helping to onboard new contributors
* Other related activities

## Getting started

See [**Developing the provider**](README.md#developing-the-provider) for more advice on how to write code and run tests for this project.


## Contribution Flow

This is a rough outline of what a contributor's workflow looks like:

* Make a fork of the repository within your GitHub account
* Create a topic branch in your fork from where you want to base your work
* Make commits of logical units (don't forget to add or modify tests too)
* Make sure your commit messages are with the proper format, quality and descriptiveness (see below)
* Fetch changes from upstream and resolve any merge conflicts so that your topic branch is up-to-date
* Push your changes to the topic branch in your fork
* Create a pull request containing that commit

We follow the GitHub workflow and you can find more details on the [GitHub flow documentation](https://docs.github.com/en/get-started/quickstart/github-flow).

### Pull Request Checklist

Before submitting your pull request, we advise you to use the following:

1. Update Go modules files `go.mod` and `go.sum` if you're changing dependencies.
2. Check if your code changes will pass both code linting checks and unit tests.
3. Ensure your commit messages are descriptive. We follow the conventions on [How to Write a Git Commit Message](http://chris.beams.io/posts/git-commit/). Be sure to include any related GitHub issue references in the commit message. See [GFM syntax](https://guides.github.com/features/mastering-markdown/#GitHub-flavored-markdown) for referencing issues and commits.
4. Check the commits and commits messages and ensure they are free from typos.

## Reporting Bugs and Creating Issues

For specifics on what to include in your report, please follow the guidelines in the issue and pull request templates when available.
Anyone can log a bug using the GitHub 'New Issue' button.  Please use
a short title and give as much information as you can about what the
problem is, relevant software versions, and how to reproduce it. If you
know of a fix or a workaround include that too.

## Ask for Help

The best way to reach us with a question when contributing is to ask on:

* The original GitHub issue
* Our Slack channel: https://vmwarecode.slack.com, channel `#vcda-terraform-dev`
