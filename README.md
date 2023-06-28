# Shono
Shono is a golang SDK for building transactional streaming systems through business-like abstractions.

Wow, quite a mouth full that is.

Just as a tl;dr, Shono allows you to describe what your organization looks like and what should happen when an event
occurs. It then generates the necessary artifacts to run your system.

## Why Shono?
Event-driven architectures provide enormous benefits to software systems. They are highly scalable, resilient, and 
flexible. However, they are also difficult to build and maintain. There are many factors to take into consideration 
when building an event-driven system, such as:

    - partitioning schemes
    - event ordering
    - event delivery and processing guarantees

All of the above cause a hurdle for developers who want to build event-driven systems. Shono aims to solve this problem
by providing a set of abstractions that are familiar to developers as well as business users and allow them to think 
about cause and effect in their systems.

## How does it work?
In essence, Shono generates artifacts based on an inventory you maintain. This inventory contains the scopes, concepts,
events and reactors that make up your system. Shono then generates artifacts based on this inventory that you can run
and test through the Shono CLI.

More information on how to use Shono can be found in the [documentation](https://docs.shono.io).

## Installation
The shono library can be added to your project by executing the following command:

```shell
go get github.com/shono-io/shono
```

## What's next?
Well, we are only at the start of our journey, but we have a lot of ideas on how to improve Shono. This is what we
have lined up for the time being:

- [ ] Add support for more systems to inject from or extract to
- [ ] Create a YAML declarative language to define the inventory
- [ ] Allow ReST APIs to automatically be generated based on a concept
- [ ] Extend the DSL to allow for more complex scenarios
- [ ] Generate documentation based on the inventory