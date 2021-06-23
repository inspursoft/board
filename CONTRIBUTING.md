# Contributing to Board

## 

## Welcome

Board is developed in the open, and is constantly being improved by our **users, contributors, and maintainers**. It is because of you that we can bring great software to the community.

This guide provides information on filing issues and guidelines for open source contributors. **Please leave comments / suggestions if you find something is missing or incorrect.**

Contributors are encouraged to collaborate using the following resources in addition to the GitHub [issue tracker](https://github.com/inspursoft/board/issues):

- [Monthly public community meetings](https://github.com/inspursoft/board/wiki/community-meeting-schedule)
- Chat with us on the Slack ([login here](https://openboard-workspace.slack.com/) )

## Getting Started

Once you've chosen what to code(feature/enhancement) or fix(issue/broken link/typo), you can begin the step-by-step process below:

- **Steps 1 through 6** are setup steps, meaning you only have to do them once for Board project.
- **Steps 7 through 19** should be repeated for each contribution to Board project.

### Step 1: Sign into GitHub

Sign into your GitHub account, or [create a free GitHub account](https://github.com/join) if you don't have one.

### Step 2: Fork the project repository

Find the project's repository on GitHub, and then "fork" it by clicking the **Fork** button in the upper right corner. This creates a copy of the project repository in your GitHub account. In the upper left corner, you will see that you are now looking at a  repository in your account.

### Step 3: Clone your fork

While still in your repository, click the green **Clone or download** button and then copy the HTTPS URL:  `git clone https://github.com/[YOURGITHUBID]/board.git`.

Using Git on your local machine, clone your fork using the URL you just copied.

Cloning copies the repository files (and commit history) from GitHub  to your local machine. The repository will be downloaded into a  subdirectory of your working directory, and the subdirectory will have  the same name as the repository.

(If you run into problems during this step, read the [Set up Git](https://help.github.com/en/github/getting-started-with-github/set-up-git) page from GitHub's documentation.)

### Step 4: Navigate to your local repository

Since the clone was downloaded into a subdirectory of your working directory, you can navigate to it using: `cd board`.

```
board
|-- LICENSE
|-- Makefile
|-- README.md
|-- README_zh_CN.md
|-- VERSION
|-- docs             <--- document
|-- make             <--- CI/CD script
|-- src              <--- source code
|-- tests            <--- test script
`-- tools            <--- related tools
```

### Step 5: Check that your fork is the "origin" remote

You are going to be synchronizing your local repository with both the project repository (on GitHub) and your fork (also on GitHub). The URLs that point to these repositories are called "remotes". More  specifically, the project repository is called the "upstream" remote,  and your fork is called the "origin" remote.

When you cloned your fork, that should have automatically set your fork as the "origin" remote. Use `git remote -v` to show your current remotes. You should see the URL of your fork (which you copied in step 3) next to the word "origin".

If you don't see an "origin" remote, you can add it using: `git remote add origin https://github.com/[YOURGITHUBID]/board.git`.

```
# git remote -v
origin	https://github.com/[YOURGITHUBID]/board.git (fetch)
origin	https://github.com/[YOURGITHUBID]/board.git (push)
```

(If you run into problems during this step, read the [Managing remote repositories](https://help.github.com/en/github/using-git/managing-remote-repositories) page from GitHub's documentation.)

### Step 6: Add the project repository as the "upstream" remote

Go to your fork on GitHub, and click the "forked from" link to return to the Board project repository. Add the project repository as the "upstream" remote using: `git remote add upstream https://github.com/inspursoft/board.git`.

Use `git remote -v` to check that you now have two  remotes: an origin that points to your fork, and an upstream that points to the project repository.

```
#git remote -v
origin	https://github.com/[YOURGITHUBID]/board.git (fetch)
origin	https://github.com/[YOURGITHUBID]/board.git (push)
upstream	https://github.com/inspursoft/board.git (fetch)
upstream	https://github.com/inspursoft/board.git (push)
```

### Step 7: Pull the latest changes from upstream into your local repository

Before you start making any changes to your local files, it's a good  practice to first synchronize your local repository with the project  repository. Use `git pull upstream master` to "pull" any changes from the "master" branch of the "upstream" into your local repository.

If you forked and cloned the project repository just a few minutes  ago, it's very unlikely there will be any changes, in which case Git  will report that your local repository is "already up to date". But if  there are any changes, they will automatically be merged into your local repository.

### Step 8: Create a new branch

Rather than making changes to the project's "master" branch, it's a  good practice to instead create your own branch. This creates an  environment for your work that is isolated from the master branch.

Use `git checkout -b BRANCH_NAME` to create a new branch  and then immediately switch to it. The name of the branch should briefly describe what you are working on, and should not contain any spaces.

The branch should be named  `XXX-description` where XXX is  the number of the issue. PR should be rebased on top of master without  multiple branches mixed into the PR. If your PR do not merge cleanly,  use commands listed below to get it up to date.

```
cd $working_dir/board
git fetch origin
git checkout master
git rebase board/master
```

​                    

Branch from the updated `master` branch:

```
git checkout -b XXX-description master
```

### Step 9: Make changes in your local repository

Use a text editor or IDE to make the changes you planned to the files in your local repository. Because you checked out a branch in the  previous step, any edits you make will only affect that branch.

The coding style used in Harbor is suggested by the Golang community. See the [style doc](https://github.com/golang/go/wiki/CodeReviewComments) for details.

Try to limit column width to 120 characters for both code and markdown documents such as this one.

As we are enforcing standards set by [golint](https://github.com/golang/lint), please always run golint on source code before committing your changes. If it reports an issue, in general, the preferred action is to fix the  code to comply with the linter's recommendation because golint gives suggestions according to the stylistic conventions  listed in [Effective Go](https://golang.org/doc/effective_go.html) and the [CodeReviewComments](https://github.com/golang/go/wiki/CodeReviewComments).

Unit test cases should be added to cover the new code. Unit test framework for backend services is using [go testing](https://golang.org/doc/code.html#Testing). The UI library test framework is built based on Angular, please refer to [Angular Testing](https://angular.io/guide/testing) for more details.

To build the code, please refer to [build](https://github.com/inspursoft/board/blob/dev/docs/compile_guide.md) guideline.

Board uses [swagger](https://swagger.io) to define API. To add or change the APIs, first update the `docs/swagger.yaml` file.

### Step 10: Commit your changes

After you make a set of changes, use `git add -A` to stage your changes and `git commit -m "DESCRIPTION OF CHANGES"` to commit them.

If you are making multiple sets of changes, it's a good practice to make a commit after each set.

### Step 11: Push your changes to your fork

When you are done making all of your changes, upload these changes to your fork using `git push origin BRANCH_NAME`. This "pushes" your changes to the "BRANCH_NAME" branch of the "origin" (which is your fork on GitHub).

For example, I used `git push origin XXX-description`.

### Step 12: Begin the pull request

Return to your fork on GitHub, and refresh the page. You may see a highlighted area that displays your recently pushed branch.

Click the green **Compare & pull request** button to begin the pull request.

(Alternatively, if you don't see this highlighted area, you can switch to your branch using the **Branch** button and then click the **New pull request** button.)

### Step 13: Create the pull request

When opening a "pull request", you are making a "request" that the  project repository "pull" changes from your fork. You will see that the  project repository is listed as the "base repository", and your fork is  listed as the "head repository".

Before submitting the pull request, you first need to describe the  changes you made (rather than asking the project maintainers to figure  them out on their own). You should write a descriptive title for your  pull request, and then include more details in the body of the pull  request. If there are any related GitHub issues, make sure to mention  those by number. The body can include Markdown formatting, and you can  click the **Preview** tab to see how it will look.

On the right side, you may see a link to the project's **Contributing** guidelines. This is primarily worth reading through if you are  submitting substantial code (rather than just fixing a typo), but it may still be worth scanning through at this point.

Below the pull request form, you will see a list of the commits you  made in your branch, as well as the "diffs" for all of the files you  changed.

If everything looks good, click the green **Create pull request** button!

### Step 14: Review the pull request

You have now created a pull request, which is stored in the project's repository (not in your fork of the repository). It's a good idea to  read through what you wrote, as well as clicking on the **Commits** tab and the **Files changed** tab to review the contents of your pull request.

If you realize that you left out some important details, you can click  the **'...'** button in the upper right corner to edit your pull request  description.

### Step 15: Add more commits to your pull request

You can continue to add more commits to your pull request even after  opening it! For example, the project maintainers may ask you to make  some changes, or you may just think of a change that you forgot to  include.

Start by returning to your local repository, and use `git branch` to see which branch is currently checked out. If you are currently in  the master branch (rather than the branch you created), then use `git checkout XXX-description` to switch. 

Then, you should repeat steps 9 through 11: make changes, commit them, and push them to your fork.

Finally, return to your open pull request on GitHub and refresh the  page. You will see that your new commits have automatically been added  to the pull request.

### Step 16: Discuss the pull request

If there are questions or discussion about your pull request from the project maintainers, you can add to the conversation using the comment  box at the bottom of the pull request.

If there are inline comments about specific changes you made, you can respond to those as well.

Click the **Resolve conversation** button once you have addressed any specific requests.

### Step 17: Delete your branch from your fork

If the project maintainer accept your pull request  (congratulations!), we will merge your proposed changes into the  project's master branch and close your pull request.

You will be given the option to delete your branch from your fork, since it's no longer of any use.

### Step 18: Delete your branch from your local repository

You should also delete the branch you created from your local  repository, so that you don't accidentally start working in it the next  time you want to make a contribution to this project.

First, switch to the master branch: `git checkout master`.

Then, delete the branch you created: `git branch -D XXX-description`. 

### Step 19: Synchronize your fork with the project repository

At this point, your fork is out of sync with the project repository's master branch.

To get it back in sync, you should first use Git to pull the latest  changes from "upstream" (the project repository) into your local  repository: `git pull upstream master`.

Then, push those changes from your local repository to the "origin" (your fork): `git push origin master`.

If you return to your fork on GitHub, you will see that the master branch is "even" with the project repository's master branch.



## Reporting issues

It is a great way to contribute to Board by reporting an issue.  Well-written and complete bug reports are always welcome! Please open an issue on GitHub and follow the template to fill in required  information.

Before opening any issue, please look up the existing [issues](https://github.com/inspursoft/board/issues) to avoid submitting a duplication. If you find a match, you can "subscribe" to it to get notified on  updates. If you have additional helpful information about the issue,  please leave a comment.

When reporting issues, always include:

- Version of docker engine and docker-compose
- Configuration files of Board
- Log files in /var/log/b/oard

Because the issues are open to the public, when submitting the log  and configuration files, be sure to remove any sensitive information,  e.g. user name, password, IP address, and company name. You can replace those parts with "REDACTED" or other strings like "****".

Be sure to include the steps to reproduce the problem if applicable. It can help us understand and fix your issue faster.



## Documenting

Update the documentation if you are creating or changing features. Good documentation is as important as the code itself.

The main location for the documentation is the [document directory](https://github.com/inspursoft/board/docs). The images referred to in documents can be placed in `docs/img` in that repo.

Documents are written with Markdown. See [Writing on GitHub](https://help.github.com/categories/writing-on-github/) for more details.



## Congratulations!

Congratulations on making your first Board contribution! 🎉 If you ran into any unexpected problems, we would love to hear about it so that we can continue to improve this guide.