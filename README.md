# GoLang saga

This package provides utilities which allow usage of such pattern as Saga. This
pattern is usually used to control data in distributed systems and is often
used in microservices architecture.

Sagas themselves should follow the requirements described by the **ACID** 
abbreviation:

- **Atomicity** states, that each transaction will be completed fully or not 
completed at all;
- **Consistency** follows from Atomicity. Transaction should not allow the 
existence of intermediate results. So, it should grant consistency;
- **Isolation** states, that there could not be several same transactions which
are currently executing at the same time;
- **Durability** grants, that transaction was completed and changes done will
not be canceled.

## How it works

Each saga is described by some count of steps, which are able to roll back. 
Let's imagine an operation of buying an item in some shop, it could consist of
these steps:

1. Decrease the user balance;
2. Create a record in database which states, that the item belongs to the user;
3. Increase item purchase counter to probably hide it from the store in case,
the item ran out of stock.

In case, some of these steps failed, we should roll back all previously 
completed steps to follow the **Atomicity** rule. So, we should describe 
rollback action for each step:

1. Increase user balance;
2. Remove created record in database;
3. Decrease purchase counter.

Additionally, all these steps could be executed by several services. So, how
do we solve this problem via this library? We should create a saga, which 
describes all steps above and perform it.

## Saga

It is highly recommended to create some sort of factory which 
contains all saga dependencies and can execute sagas with some custom arguments.
Dependencies are usually service interfaces which are able to perform some
service related actions. Custom arguments are dynamic parameters which could
differ from one saga to another.

Current library requires usage of structure, which implements `Saga` 
interface (see `saga.go`). Let's describe each of its methods.

### `Run`

This method describes all steps contained by current saga. It accepts following
arguments:

1. **Contexts**: 
   1. **Execution context**. This is the main context used by all saga steps. You
   could probably pass request context here which could have deadline or just 
   be canceled;
   2. **Rollback context**. This is the separate context due to the reason,
   execution context could be timed out and as the result, one of the steps 
   would fail. This will lead the saga to start rollback phase which will 
   instantly fail as long as context is dead. This is why you want rollback 
   context to be separate entity - to control rollback lifecycle separately;
2. **Saga runner**. This structure allows developer to define, run and rollback 
steps.

As the result, function returns value of specified generic type and 2 errors:
1. **Execution error**. The first occurred error during saga execution;
2. **Rollback error**. Error occurred while rolling saga. You want this error
to be returned separately as long rollback error is rather important as it
will lead to saga memoization and rollback retry in the future.

### `Lock`

This method locks saga and prevents it to be concurrently called.

### `Unlock`

Unlocks saga.

### `ShouldRetryUnlock`

Determines if saga unlock should be retried. To learn more, take a look at
[`ShouldRetruFunc`](https://github.com/heyqbnk/go-saga/blob/master/retry.go#L5).

## Example

You can find good example covering all library aspects 
[here](https://github.com/heyqbnk/go-saga/blob/master/example/main.go). To 
launch an example:

```bash
go run example/main.go
```

To try other cases, try returning errors from the saga dependencies: 
`user-items`, `wallet`, `locker` and `store`. 