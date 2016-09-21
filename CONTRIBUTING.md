# How to contribute to a nullstyle/go project

There are a few guidelines that we
ask contributors to follow so that we can merge your changes quickly.

## Getting Started

* Make sure you have a [GitHub account](https://github.com/signup/free)
* Create a GitHub issue for your contribution, assuming one does not already exist.
  * Clearly describe the issue including steps to reproduce if it is a bug.
* Fork the repository on GitHub

## Finding things to work on
The first place to start is always looking over the current github issues for the project you are interested in contributing to. Issues marked with [help wanted](https://github.com/issues?q=is%3Aopen+is%3Aissue+user%3Anullstyle+label%3A%22help+wanted%22) are usually pretty self contained and a good place to get started.

Of course feel free to make your own issues if you think something needs to added or fixed.

## Making Changes

* Create a topic branch from where you want to base your work.
  * This is usually the master branch.
  * Please avoid working directly on the `master` branch.
* Make sure you have added the necessary tests for your changes and make sure all tests pass.

## Submitting Changes

* [Sign the Contributor License Agreement](TODO)
* Push your changes to a topic branch in your fork of the repository.
 * Include a descriptive [commit message](https://github.com/erlang/otp/wiki/Writing-good-commit-messages).
 * Changes contributed via pull request should focus on a single issue at a time.
 * Rebase your local changes against the master branch. Resolve any conflicts that arise.

At this point you're waiting on us. We may suggest some changes, improvements or alternatives.

## Minor Changes

### Documentation
For small changes to comments and documentation, it is not
always necessary to create a new GitHub issue. In this case, it is
appropriate to start the first line of a commit with 'doc' instead of
an issue number.

# Additional Resources
* [Contributor License Agreement](TODO)


This document is inspired by:

https://github.com/stellar/docs/blob/master/CONTRIBUTING.md

https://github.com/puppetlabs/puppet/blob/master/CONTRIBUTING.md

https://github.com/thoughtbot/factory_girl_rails/blob/master/CONTRIBUTING.md

https://github.com/rust-lang/rust/blob/master/CONTRIBUTING.md