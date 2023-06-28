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

## The Inventory
The inventory contains representations of the components that make up your system. It allows you to split your system
in so-called "scopes" which are roughly equivalent to bounded contexts in Domain-Driven Design. Each scope contains
concepts, which are the building blocks of your system. Good examples of concepts are 'Employee', 'Product' or 
'Invoice' Each of these concept has a set of events that can be emitted by it.

Concepts also allow for the definition of reactors. Reactors are the components that process events. They are the
building blocks of your system. A reactor can be a simple component that emits an event when it receives one, or it
can contain complex logic to emit multiple events based on the event it receives.

I know by this point you are probably thrilled to know what else is in the inventory, so lets dive a little bit deeper.

### Scope
There was a time when the goal of IT was to come up with one model to describe the whole organization; one model to 
rule them all - of sorts. When talking about a customer, everyone would know what you meant. But this is an illusion. 
The reality is that the same concept can mean different things in different contexts. A customer in the sales 
department is not the same as a customer in the finance department. This is why we need to define scopes.

A scope is a mental boundary that you define in your system. Within this mental boundary, concepts mean the same thing
and any data stored within the scope is internal to that scope; it is not shared across scopes. Events are shared 
instead so that other scopes can react to them and interpret what they mean within their context.

How you structure your scopes is up to you. It is perfectly fine to create a scope for each department of your 
organization since within a department, concepts will mean the same thing. You can ever go finer grained and create
a scope for each team within a department. The important thing is that you understand scopes to be mental boundaries
and only events are share across scopes.

### Concept
A concept has very little meaning on its own. It is a label given to a real-world entity that is relevant within your
scope. A concept can be a customer, an invoice, a product, etc.

The state associated with a concept is private to the scope it is defined in. Concepts themselves are not shared 
across scopes. Instead, events are emitted which can be interpreted by other scopes to adapt their internal state. It
allows for a loose coupling between scopes as well as the ability to interpret what an event means within a specific
context.

While Shono's main focus is on transactional processing, even basic analytics can be done by interpreting events. For
example, it is perfectly fine to have a `CustomerStats` concept that is updated each time a customer related event is
being processed.

### Event
An event is something that happens to a concept. It is a fact that is emitted about the concept and can be picked up
by other scopes to act upon. Events are the only thing that is shared across scopes.

Event names are usually in the past tense. For example, if a customer is created, an event named `CustomerCreated`
could be emitted. If a customer is updated, an event named `CustomerUpdated` could be emitted. Events are not limited
to CRUD operations. They can be anything that happens to a concept. For example, if a customer is moved to a different 
department, an event named `CustomerMoved` could be emitted.

When working with external facing systems like APIs etc, it is common to have events like `CustomerCreationRequested`
because from the perspective of the system, the API is receiving a request to create a customer. I know this might 
sound a bit controversial, but we can always spend some time discussing this over a beer, sigar or good whisky. You're
buying btw; it will be a long discussion.

### Reaktor
Let's move on to the components which are actually performing the work in your system, the first of which are reactors.
Reactors define what happens when an event is received and are always functioning from the perspective of the concept
they are linked to. Some reactors can be relatively simple, like just logging the event and emitting a new one. Others
are more involved and interpret the event to update the state of the concept.

A reactor can listen for any event, from any concept in any scope. However, it can only emit events for the concept it
is linked to. You use a domain specific language to define what happens when an event is received. For example, if a
`CustomerCreationRequested` event is received, we could execute the following steps:

- validate the event against a schema
- check if the customer already exists
- if not, create the customer
- emit a `CustomerCreated` event

But logic doesn't just focus on the success flow; you can also handle errors, all from the comfort of the DSL.

Now that is already a lot of fun, but we can take it a step further. Defining logic is one thing, but how will
we know if the logic we defined is actually correct? Well, that's where tests come in. Aside from the different steps
that make up the processing logic, you can also define tests that will be executed against the logic. These tests can
be run using the CLI to verify that the logic is performing as expected.

### Injector
There are times however when you want to emit events based on state changes from your existing systems. For example,
many organizations already have a CRM system in place that contains all their customer data. You would actually want
to emit an event each time a customer is created in the CRM system so that other systems can react to it.

Meet the injector. The injector can read from existing systems and emit events based on the data it reads. It will
be up to you to define the logic that reads the data and emits the events, but you can do so using the same DSL as
you would for a reactor, including the testing and all. The only difference is that the injector is not linked to a
single concept. It can emit events for any concept within the scope it is linked to.

### Extractor
Injectors provide a way to get events from your existing systems, but there certainly are times you want to have the
opposite as well; reacting to events and writing them to systems already part of your organization. An extractor does
just that. It listens for events and executes the logic you defined to write the data to whatever system you want.

And yes, the same logic DSL is here to help out again.

