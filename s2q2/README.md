# Section 2 - Question 2

This is the solution for the queue challenge.

## Approach

```
                       ┌───────────┐
                       │ Invisible │
                       │           │ ----+
                       └───────────┘     |
                         ^      |        |
                         | View | 1s     |
                         |      V        |
┌───────────┐          ┌───────────┐     |
│ Not In    │   Add    │ In        │     | Delete
│ Queue     │ -------> │ Queue     │     |
└───────────┘          └───────────┘     |
      ^                      |           |
      |                      | Delete    |
      |                      V           |
      |                ┌───────────┐     |
      |      1s        │ Deleted   │     |
      +--------------- │           │ <---+
                       └───────────┘
```
