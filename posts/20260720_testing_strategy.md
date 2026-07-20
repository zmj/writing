# How to Test

This is an attempt to capture my current thinking on automated testing best practices. This will be the template for a policy document in a codebase I work on.

First, what's in scope? These thoughts are relevant to "backend services", meaning applications that:
* present an API: structured, machine-parseable inputs and outputs
* depend on external state: probably a database, maybe some peer services
* respond to requests: blending internal logic and external state

TL;DR: three layers. Three complementary kinds of tests.

## API tests

These tests call the service through its API, exactly like its consumers do. If I could only write one kind of test, it would be these. They might be called end-to-end tests, or integration tests - describing the intent to exercise the complete call stack, including external interactions. The execution environment for these tests should replicate production as faithfully as possible.

API tests alone can suffice for small services. As the codebase and test suite grow, two problems surface:
* the tests are slow. At minute durations, they drag on the development inner loop. At hour durations, they bottleneck deployment.
* the tests are flaky. Each test can fail for many reasons unrelated to the behavior it verifies. Retries make the speed problem worse.

API tests are a finite resource.

## Logic tests

These tests are the opposite of API tests - isolated and fast. They're commonly called 'unit tests', where a 'unit' is some meaningful subset of the service's functionality. This kind of test is ideal for purely functional logic: the same input always produces the same output. The challenge is testing logic that depends on both the input and external state. What stands in for the external authority?

The usual approach is for the logic to define an interface covering all potential interactions with external state. Logic tests set up expected interactions in terms of that interface. The tests verify that the expected interactions happened, and that the output was correct conditional on them. This is called "mock testing", and it's the wrong approach.

There are two problems with mock tests:
* within a single test, they specify too much. There are often many possible interactions with external state that all produce the correct output. A mock test incorrectly flags some of them as failures.
* across all tests, they specify too little. The behavior of the external authority is described by each test. Nothing ensures that behavior is consistent or complete.

Instead, satisfy the logic's interface with a "fake": a second implementation with in-memory internal state. All tests use the same fake implementation, ensuring consistency. The tests are agnostic how the logic interacts with the fake; they only verify the output.

## Boundary tests

These tests verify the service's interactions with external state. They run against the same interfaces that support the logic tests. For example, tests for a database boundary should verify:
* serialization: do writes and reads round-trip the same values?
* filters, sorting: do computations within the database correctly affect query results?
* preconditions and postconditions: across a series of interactions, does the database's state change as expected?

Boundary tests run twice: once against the actual external authority, and once against its fake implementation. The same interactions and assertions run against both implementations. That consistency grants credibility to the logic tests: they verify production behavior, because external state acts as it would in production.

## The final picture

A well-tested service has three test suites:
* logic tests on most of the code, with fakes at external interaction boundaries
* boundary tests, keeping the fakes and production implementations correct and consistent
* API tests, confirming that logic and boundaries are correctly wired end-to-end
