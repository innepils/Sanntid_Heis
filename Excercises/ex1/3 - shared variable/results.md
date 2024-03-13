// Excercise 1

* Task 3
We discovered that each execute of the program results in random values being printed, due to race conditions since we are not locking the cricital sections. 

* Task 4
We think Mutex is the right synchronization mechanism in this scenario, over Semaphore. This is because we are working with a single resource, which we only want one of the threads to be working on at a time. 




