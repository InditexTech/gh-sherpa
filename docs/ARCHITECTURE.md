# Architecture

## Use case flows

### Create branch

<details open>
<summary>Show/hide</summary>

```mermaid
graph TD;
  Start([Start]) --> getIssue

  getIssue -->|Not found| getIssueError([Error])
  getIssue(Get issue from issue tracker) --> getIssueType

  getIssueType -->|Unknown type| getIssueTypeErr([Error])
  getIssueType(Get issue type from issue **) --> generateBranchName
  
  generateBranchName -->|Branch already exists| generateBranchNameErr([Error])
  generateBranchName(Generate branch name from Issue **) --> getBaseBranch

  getBaseBranch(Get base branch from repository) -->  checkoutBranch

  checkoutBranch(Create branch from origin **) --> End

  End([End])

```
*\*\* The user may need to enter input manually in interactive mode*

</details>

### Create pull request

<details open>
<summary>Show/Hide</summary>

```mermaid
graph TD;
  Start([Start]) --> getBranchName

  getBranchName(Get branch name **) --> checkLocalBranch 
  
  checkLocalBranch(Check local branch exists)
  checkLocalBranch -->|Branch exists| checkoutBranch
  checkLocalBranch -->|Missing branch| createLocalBranch

  createLocalBranch(Create local branch **) --> checkoutBranch

  checkoutBranch(Checkout branch) --> checkPendingCommits

  checkPendingCommits(Check local pending commits **)
  checkPendingCommits -->|Pending commits| pushBranch
  checkPendingCommits-->|No commits| checkRemoteBranch

  checkRemoteBranch(Check remote branch exists)
  checkRemoteBranch -->|No remote branch| createEmptyCommit
  checkRemoteBranch -->|Remote branch exists| checkPrExists

  createEmptyCommit(Create empty commit **) --> pushBranch

  pushBranch(Push branch to remote **) --> checkPrExists

  checkPrExists(Check open PR already exists)
  checkPrExists -->|Open PR exists| checkPrExistsErr([Error])
  checkPrExists -->|No open PR| createPr

  createPr(Create Pull Request **) --> End

  End([End])
```
*\*\* The user may need to enter input manually in interactive mode*

</details>

## CLI flow

<details open>
<summary>Show/hide</summary>

```mermaid
graph TD;
  Start([Start]) --> 

  readArgs(Read args) --> validateArgs

  validateArgs -->|Invalid| validateArgsErr([Error])
  validateArgs(Validate args) --> validateConfig

  validateConfig -->|Invalid| validateConfigErr([Error])
  validateConfig(Validate configuration) --> checkGhAuth

  checkGhAuth -->|Unauthenticated| checkGhAuthErr([Error])
  checkGhAuth[/Check github authentication/] --> checkGitRepo

  checkGitRepo -->|Not a git repository| checkGitRepoErr([Error])
  checkGitRepo(Check execution within git repository) -->

  E[[Execute use case]] -->

  End([End])
```

</details>