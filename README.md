# bufwr

bufwr is a buffered writer that accumulates input data in an internal buffer until the buffer is filled up or some defined amount of time is pass.

It especially useful when you need both performance of buffered writing and predictable delay in case of temporally lag of input data.

bufwr is attentive to data bounds and does not split apart an input data chunk on flushing. So you can be sure that the underlying writer will always get a set of the whole data chunks.

## Example

Good example of usage might be when you need to write huge amount of JSON/CSV/... messages to Kafka.
- You want to pack some number of messages together in one block on writing for performance.
- And a message should not be split apart and be written in two different blocks.
- Also you don't want that messages get stuck in a buffer for а long time when the incoming message stream is poor/stopped.
