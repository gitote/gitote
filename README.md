<div align="center">
  <br>
  <img
    alt="Gitote"
    src="https://i.imgur.com/mB7XYLs.png"
    width=180px
  />
  <br/>
  <h1>👩‍💻👨‍💻 Gitote 👩‍💻👨‍💻</h1>
  <p align="center">
  	<strong>
  		<a href="https://gitote.com">Website</a>
  		•
  		<a href="https://docs.gitote.com">Docs</a>
  		•
  		<a href="https://github.com/gitote/gitote/issues">Issues</a>
  	</strong>
  </p>
  <p align="center">
    <a href="https://github.com/gitote/gitote"><img
  		alt="Go Version"
  		src="https://img.shields.io/badge/Go-v1.11-brightgreen.svg">
  	</a>
  	<a href="https://ci.appveyor.com/project/yogicodes/gitote"><img
  		alt="Appveyor Status"
  		src="https://ci.appveyor.com/api/projects/status/lhunpdfuay0oy8fj?svg=true">
  	</a>
  	<a href="https://goreportcard.com/report/github.com/gitote/gitote"><img
  		alt="Go Report"
  		src="https://goreportcard.com/badge/github.com/gitote/gitote">
  	</a>
  	<a href="https://codebeat.co/projects/github-com-gitote-gitote-master"><img
  		alt="Codebeat Report"
  		src="https://codebeat.co/badges/7ec50fde-7899-4673-ad9e-3ce7740de99d">
  	</a>
  </p>
</div>
<br/>

Welcome to the [gitote](https://gitote.com) codebase. We are so excited to have you. With your help, we can build out Gitote to be more stable and better serve our platform.

## What is gitote?

[Gitote](https://gitote.com) is an open source end-to-end software development platform with built-in version control, issue tracking, code review, and more.

## Contributing

We expect contributors to abide by our underlying [code of conduct](docs/CONDUCT.md). All conversations and discussions on GitHub (issues, merge requests) and across Gitote must be respectful and harassment-free.

### Where to contribute

When in doubt, ask a [core team member](#core-team)! You can mention us in any issues . Any issue with `Good first Issue` tag is typically a good place to start.

**Refactoring** code, e.g. improving the code without modifying the behavior is an area that can probably be done based on intuition and may not require much communication to be merged.

**Fixing bugs** may also not require a lot of communication, but the more the better. Please surround bug fixes with ample tests. Bugs are magnets for other bugs. Write tests near bugs!

**Building features** is the area which will require the most communication and/or negotiation. Every feature is subjective and open for debate. The [product roadmap](https://github.com/gitote/gitote/projects/1) should be a good guide to follow. As always, when in doubt, ask!

### How to contribute

1.  Fork the project & clone locally. Follow the initial setup [here](#getting-started).
2.  Create a branch, naming it either a feature or bug: `git checkout -b feature/that-new-feature` or `bug/fixing-that-bug`
3.  Code and commit your changes. Bonus points if you write a [good commit message](https://chris.beams.io/posts/git-commit/): `git commit -m 'Add some feature'`
4.  Push to the branch: `git push origin feature/that-new-feature`
5.  [Create a merge request](#create-a-merge-request) for your branch 🎉

## Contribution guideline

### Create an issue

Nobody's perfect. Something doesn't work? or could be done better? Let us know by creating an issue.

PS: a clear and detailed issue gets lots of love, all you have to do is follow the issue template!

#### Clean code with tests

Some existing code may be poorly written or untested, so we must have more scrutiny going forward. We test with **go test**, let us know if you have any questions about this!

#### Create a megre request

* Try to keep the megre requests small; a megre request should try its very best to address only a single concern.
* Make sure all tests pass and add additional tests for the code you submit.
* Document your reasoning behind the changes. Explain why you wrote the code in the way you did; the code should explain what it does.
* If there's an existing issue related to the megre request, reference to it by adding something like `References/Closes/Fixes/Resolves #305`, where 305 is the issue number. [More info here](https://github.com/blog/1506-closing-issues-via-pull-requests)
* If you follow the megre request template, you can't go wrong.

## Codebase

### The stack

We run on a **Go** backend.

Additional technologies and services are listed on [our docs](https://docs.dev.to).

## Features

- Activity timeline
- SSH and HTTP/HTTPS protocols
- Account/Organization/Repository management
- Add/Remove repository collaborators
- Repository/Organization webhooks (including Slack and Discord)
- Repository Git hooks/deploy keys
- Repository issues, pull requests, wiki and protected branches
- Migrate and mirror repository and its wiki
- Web editor for repository files and wiki
- Jupyter Notebook
- Two-factor authentication
- Gravatar and Federated avatar with custom source
- Mail service

## Getting Started

This section provides a high-level requirement & quick start guide.

### Hardware Requirements

- A **Raspberry Pi** or $5 **Digital Ocean Droplet** is more than enough to get you started. Some even use 64MB RAM Docker [CaaS](https://blog.docker.com/2016/02/containers-as-a-service-caas/).
- 2 CPU cores and 512MB RAM would be the baseline for teamwork.
- Increase CPU cores when your team size gets significantly larger, memory footprint remains low.

### Prerequisites

##### [Go Lang](https://golang.org)

```sh
curl -O https://storage.googleapis.com/golang/go1.11.linux-amd64.tar.gz
sha256sum go1.11.linux-amd64.tar.gz
tar -xvf go1.11.linux-amd64.tar.gz
sudo chown -R root:root ./go
sudo mv go /usr/local
export GOPATH=$HOME/go
export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin
```

##### [PostgreSQL](https://www.postgresql.org/) (9.4 or higher)

### Installation

0.  Make sure all the prerequisites are installed.
1.  Fork **Gitote** repository, ie. https://github.com/gitote/gitote/fork

```sh
mkdir -p $GOPATH/src/gitote.com/gitote
cd $GOPATH/src/gitote.com/gitote
git clone https://github.com/gitote/gitote.git
cd gitote
go build
./gitote web -p 8080
```

#### Compiling the LESS

We're mostly a go app, with a bit of **LESS** sprinkled in. **For most cases, simply running `make` will do.** If you're working with LESS though, you'll need to run the following:

* Run **`make less`** to compile all LESS files.

## Product Roadmap

Our new product roadmap can be found [here](https://github.com/gitote/gitote/projects/1). Many notes need to be converted to issues but this should provide an overview of features we plan to work on, as well as features we are considering.

Core team members will move issues along the project board as they progress.

* Backlog: an accumulation of uncompleted work or matters needing to be dealt with.
* Bug: bugs to be fixed.
* Feature Requests: features up for discussion.
* Closed: features we're committed to building -- free for contributors to work on, but please communicate with the owner beforehand.

## Core team

<table>
  <tr>
    <td>
      <a href='https://twitter.com/yogicodes'>
        <img src='https://s.gravatar.com/avatar/554506e208edaf6c95dc896105b898f0?s=100'>
      </a>
      <h4 align='center'><a href='https://twitter.com/yogicodes'>Yoginth</a></h4>
    </td>
    <td>
      <a href='https://twitter.com/GitoteHQ'>
        <img src='https://s.gravatar.com/avatar/894559413cde06d4aacc1a05bd7f3205?s=100'>
      </a>
      <h4 align='center'><a href='https://twitter.com/GitoteHQ'>Gitote</a></h4>
    </td>
  </tr>
</table>

## License

This program is free software: you can redistribute it and/or modify it under the terms of the **MIT License**. Please see the [LICENSE](https://github.com/gitote/gitote/blob/master/LICENSE) file in our repository for the full text.

<br/>

<div align="center">
  <img
    alt="sloan"
    width=220px
    src="https://i.imgur.com/WAppLz3.jpg"
  />
  <br/>
  <strong>Happy Coding</strong> ❤️
</div>