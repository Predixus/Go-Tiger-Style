# Go! Tiger Style

A set of ammendmentments to [Tigerstyle](https://github.com/tigerbeetle/tigerbeetle/blob/main/docs/TIGER_STYLE.md),
for Golang.

First, I would read [TigerStyle](https://github.com/tigerbeetle/tigerbeetle/blob/main/docs/TIGER_STYLE.md) fully.

This will provide a solid foundation for the Go! specific ammendements listed herein.

## Why Style?

As taken directly from the TigerStyle docs:

> "...in programming, style is not something to pursue directly. Style is necessary only where
> understanding is missing."

For our operations at Predixus, we are building data driven applications in Go. This, if done earnestly, results
in us going into the unknown and recovering the preverbial gold to distribute. And so, there

And so the goal of a style is to guide and support development with 3 north stars:

- Effectiveness - The software _has_ to work
- Simplicity - The design goals need to come together in as simple a manner as possible
- Robustness - There needs to be a baseline level of robustness

So, that being said, let's get into it.

## Technical Debt

There is nothing to add here. Tigerstyle nailed it. Refer to [their comments](https://github.com/tigerbeetle/tigerbeetle/blob/main/docs/TIGER_STYLE.md#technical-debt) on technical debt.

## Safety

[NASAs Power of 10](https://spinroot.com/gerard/pdf/P10.pdf) still applies. But there are some modifications
that we need to make that are specific to go.
