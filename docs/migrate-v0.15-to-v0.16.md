# Migration Guide: v0.15 to v0.16

v0.16.0 introduces a number of API changes, which should be simple to migrate.
Just follow this guide and if you still encounter problems, ask for help on Discord or feel free to create an issue.

<!-- toc -->

-   [Upstream API changes](#upstream-api-changes)
-   [GetPullRequestDiff: add PullRequestDiffOption parameter (#542)](#getpullrequestdiff)

<!-- tocstop -->

## Upstream API changes

As we aim to track API changes in Khulnasoft 1.16 with this SDK release, you may find this [summary listing of changes](https://khulnasoft.com/khulnasoft/go-sdk/issues/558) helpful.

## GetPullRequestDiff
 Added new parameter `opts PullRequestDiffOption`. Khulnasoft 1.16 will default to omit binary file changes in diffs; if you still need that information, set `opts.Binary = true`.
 Related PRs:
 - [go-sdk#542](https://khulnasoft.com/khulnasoft/go-sdk/pulls/542)
 - [khulnasoft#17158](https://github.com/go-khulnasoft/khulnasoft/pull/17158)

## ReadNotification, ReadNotifications, ReadRepoNotifications
The function now has a new return argument. The read notifications will now be returned by Khulnasoft 1.16. If you don't require this information, use a blank identifier for the return variable. 

Related PRs:
 - [go-sdk#590](https://khulnasoft.com/khulnasoft/go-sdk/pulls/590)
 - [khulnasoft#17064](https://github.com/go-khulnasoft/khulnasoft/pull/17064)
