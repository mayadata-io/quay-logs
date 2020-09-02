# Quay-logs - How to Contribute

- Find an issue to work on or create a new issue. The issues are maintained at [mayadata-io/quay-logs/issues](https://github.com/mayadata-io/quay-logs/issues).
- Claim your issue by commenting on your intent to work on it to avoid duplication of efforts.
- Fork the repository on GitHub.
- Create a branch from where you want to base your work (usually master).
- Make your changes. If you are working on code contributions, please see **Setting up the Development Environment** below.
- Relevant coding style guidelines are the [Go Code Review Comments](https://code.google.com/p/go-wiki/wiki/CodeReviewComments) and the _Formatting and style_ section of Peter Bourgon's [Go: Best Practices for Production Environments](http://peter.bourgon.org/go-in-production/#formatting-and-style).
- Commit your changes by making sure the commit messages convey the need and notes about the commit.
- Push your changes to the branch in your fork of the repository.

## Setting up your Development Environment

This project is implemented using Go and uses the standard golang tools for development and build. In addition, this project relies on Quay and Kubernetes. It is expected that the contributors:

- are familiar with working with Go;
- are familiar with Docker containers;

Run the command `go run cmd/main.go` with the `quay access token` and a `quay namespace` to create an executable.

- For more information related to quay refer [quay_faq.md](https://github.com/mayadata-io/quay-logs/blob/master/quay_faq.md).

Download the required linting tools of your choice in your code editor for proper code formatting, then start making your changes.

## Committing your work

All the commits made to the repository should be signed by your MayaData email ID. For signing-off according to the [DCO standards](http://developercertificate.org/) use -s flag.

**Example:**

> `git commit -s -m "Your commit message"`

- Commits should be as small as possible. Each commit should follow the checklist below:
  - For code changes, add tests relevant to the fixed bug or new feature.
  - Pass the compile and tests - includes spell checks, formatting, etc.
  - Commit header (first line) should convey what changed.
  - Commit body should include details such as why the changes are required and how the proposed changes.
  - DCO Signed.

## Sending Pull Requests

- Rebase to the current master branch before submitting your pull request.
- Provide appropriate comments and mention what changes have been done briefly.
- Add at least one reviewer or more from the contributors.

## Resources

- Refer the full contributing guidelines here - [Contributing to OpenEBS Maya](https://github.com/openebs/maya/blob/master/CONTRIBUTING.md#contributing-to-openebs-maya)
- Refer troubleshooting.md for more information.
- For quay-logs FAQs refer to this [link](https://github.com/mayadata-io/quay-logs/blob/master/quay_faq.md).
