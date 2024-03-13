Exercise 1 - Theory questions
-----------------------------

### Concepts

What is the difference between *concurrency* and *parallelism*?
> Concurrency is switching between operations so that it SEEMS like they run in parallel, while operations in parallelism are actually parallel.

What is the difference between a *race condition* and a *data race*? 
> Race conditions are issues related to incorrect timing or sequence of events, while data race is when two threads acces the same variable concurrently without proper synchronization mechanisms in place.
 
*Very* roughly - what does a *scheduler* do, and how does it do it?
> The scheduler decides which thread to run next. This is done by using a queue. 


### Engineering

Why would we use multiple threads? What kinds of problems do threads solve?
>  Multiple threads allows us to do different things at the same time, and not wait for one another.
    If a task has to wait/do nothing, other tasks would have to wait as well.

Some languages support "fibers" (sometimes called "green threads") or "coroutines"? What are they, and why would we rather use them over threads?
> Fibers are lightweight threads that use cooperative (not preemptive) multitasking. Fibers yield themselves so  that other fibers can run as well. Fibers are generally faster and efficient. "Fine-grained".

Does creating concurrent programs make the programmer's life easier? Harder? Maybe both?
> It is harder but, easier than solving the same problems, without concurrent-programming. Blinking leds.

What do you think is best - *shared variables* or *message passing*?
> That depends on the usecase. Shared variables is fast and timing-wise easier, but more susceptible for race conditions.


