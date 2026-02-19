Exercise 1 - Theory questions
-----------------------------

### Concepts

What is the difference between *concurrency* and *parallelism*?
> Concurrency is making it look like tasks are happening at the same time, but in reality we are just switching fast between the tasks. Parallelism is that the tasks are actually happening at the same time. 

What is the difference between a *race condition* and a *data race*? 
> Race condition is that code is depended of the order it is being executed. Data race is whenever race conditions happend on shared reources (at least one of the threads are writing to it). 
 
*Very* roughly - what does a *scheduler* do, and how does it do it?
> Scheduler chooses which thread that is going to be executed. It chooses thread by random among the runnable threads.


### Engineering

Why would we use multiple threads? What kinds of problems do threads solve?
> Multiple threads make sure that there are less "dødtid" where nothing happens. Easier to change and easier to read code. Separate code that doesnt have anything to do with eachother. Untangles the code that should be unrelated. 

Some languages support "fibers" (sometimes called "green threads") or "coroutines"? What are they, and why would we rather use them over threads?
> Using either blocking operations or nonblocking operations to deal with the fact that hardware is slow. Blocking returns when there is data, Nonblocking returns immediately, but not necessarily with data. Nonblocking threads cant use preemptive scheduling. Fiber is a way of dealing with this con (ulempe). is Fibers are connected to nonblocking threads.
Lightweight (uses less memory), cooperative threads. Fibers are not interrupted at random times.  
Fibers: 
Does creating concurrent programs make the programmer's life easier? Harder? Maybe both?
> Maybe both

What do you think is best - *shared variables* or *message passing*?
> Message passing is probably better. It copies and sends values when needed to share. 


