# Contributing to IBMCloud Cost Estimator SDK

**First:** if you're unsure  _anything_, just ask or submit the issue or pull request anyways. We appreciate any sort of contributions.

However, for those individuals who want a bit more guidance on the best way to contribute to the project, read on. This document will cover what we're looking for. By addressing all the points we're looking for, it raises the chances we can quickly merge or address your contributions.

Specifically, we have provided checklists below for each type of issue and pull request that can happen on the project. These checklists represent everything we need to be able to review and respond quickly.

## Issues



#### Bug Reports

 - [ ] __Test against latest release__: Make sure you test against the latest released version. It is possible we already fixed the bug you're experiencing.

 - [ ] __Search for possible duplicate reports__: It's helpful to keep bug reports consolidated to one thread, so do a quick search on existing bug reports to check if anybody else has reported the same thing. You can scope searches by the label "bug" to help narrow things down.

 - [ ] __Include steps to reproduce__: Provide steps to reproduce the issue, along with your `.tf` `plan.json` files, with secrets removed, so we can try to reproduce it. Without this, it makes it much harder to fix the issue.

#### Feature Requests

 - [ ] __Search for possible duplicate requests__: It's helpful to keep requests consolidated to one thread, so do a quick search on existing requests to check if anybody else has reported the same thing. You can scope searches by the label "enhancement" to help narrow things down.

 - [ ] __Include a use case description__: In addition to describing the behavior of the feature you'd like to see added, it's helpful to also lay out the reason why the feature would be important and how it would benefit Terraform users.


## Pull Requests

Thank you for contributing! Here you'll find information on what to include in your Pull Request to ensure it is accepted quickly.

 * For pull requests that follow the guidelines, we expect to be able to review and merge very quickly.
 * Pull requests that don't follow the guidelines will be annotated with what they're missing. A community or core team member may be able to swing around and help finish up the work, but these PRs will generally hang out much longer until they can be completed and merged.

### Checklists for Contribution

There are several different kinds of contribution, each of which has its own standards for a speedy review. The following sections describe guidelines for each type of contribution.


#### Enhancement/Bugfix to a Resource

Working on implementing cost estimate logic for existing resources is a great way to get started as a sdk contributor because you can work within existing code and tests to get a feel for what to do.

 - [ ] __Well-formed Code__: Do your best to follow existing conventions you see in the codebase, and ensure your code is formatted with `go fmt`. (The Travis CI build will fail if `go fmt` has not been run on incoming code.) The PR reviewers can help out on this front, and may provide comments with suggestions on how to improve the code.

#### Cost estimate for New Resource

Implementing a cost estimator logic for new resource is a good way to learn more about how Terraform plan configuration is responsible for calculating cost using the plan.json input parameters. There are plenty of examples to draw from in the existing resources, but you still get to implement something completely new.

 - [ ] __Minimal LOC__: It can be inefficient for both the reviewer and author to go through long feedback cycles on a big PR with many resources. We therefore encourage you to only submit **1 resource at a time**.

 - [ ] __Well-formed Code__: Do your best to follow existing conventions you see in the codebase, and ensure your code is formatted with `go fmt`. (The Travis CI build will fail if `go fmt` has not been run on incoming code.) The PR reviewers can help out on this front, and may provide comments with suggestions on how to improve the code.


### Writing Acceptance Tests

Cost Estimator SDK includes an acceptance test harness that does most of the repetitive work involved in testing a resource.

### Steps for writting and running the acceptance test

 * Generate plan.json for the terraform template for which you are calculating the resource. To  generate plan.json run the following command
 ```bash
 terraform init

 terraform plan --out tfplan.binary

 terraform show -json tfplan.binary > tfplan.json
 ``` 
 * Update the testplan.json with updated plan.json and verify the estimated cost for current state and previous state(create,update and destroy scenarios).

 * Note: Doing a terraform apply might cost you money as it would provision the infrastructure on your account but a terraform plan will not incurr any cost.

 * Note  before doing release do go get github.com/IBM-Cloud/terraform-cost-estimator
 then do go mod vendor inside tfcost directory
 for Private Repo export GOPRIVATE="github.com/IBM-Cloud/terraform-cost-estimator"
 * Note after adding a new feature in the /ibm (sdk) you need to raise a Pr. After your PR is merged to add the new features to the tfcost tool you need to go get the latest /ibm inside the tfcost. 
 Do the following inside /tfcost after new features has been merged.
 ```
 cd tfcost
 go get github.com/IBM-Cloud/terraform-cost-estimator
 go mod vendor
 cd ..
 go build 
 ```