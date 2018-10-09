Event sourcing with Kafka in Go
===============================
Interactive demo accompanying talk at Stockholm Go Conference 2018 and
[blog post](https://www.formulate.app/formulate-hq-blog/event-sourcing-with-kafka-in-go).


## Requirements
Go and Kafka. On mac:

```
brew intsall golang
brew install kafka
brew services start zookeeper
brew services start kafka
```


## How to run
1. Clone.

2. Build

   ```
   cd orchestrator
   make all
   ```

3. Start tracker and worker from the `bin` directory.
   Send batch jobs to the tracker using `bin/cli --create --job-count 10`.
   See the code for other arguments to `bin/cli`.


## Notes on debugging
The makefile `orchestrator/kafka.mk` contains some useful commands for inspecting and interacting with Kafka, e.g. create topics and list consumer group offsets.

### If execution freezes
Make sure to wait a couple of minutes without interacting with Kafka before doing any debugging.
Sometimes it takes Kafka a while to rebalance the worker consumer group (changing which worker reads from which partition) during which time no reading is done.

If it is still blocked try resetting the consumer group offsets or delete and recreating topics and the consumer groups.
