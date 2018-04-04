# Ontology Test Framework
Ontology Test Framework is a light-weight test framework for ontology. Integration ontology-sdk to run Ontology test case.

## How to use?

1. copy wallet file from bookkeeper ontology node to ontology-test. Otherwise some testcase will failed because of balance of ont is zero.
2. Set rpc server address of ontology, wallet file and password in config_test.json config file.

```
{
  "JsonRpcAddress":"http://localhost:20336",
  "RestfulAddress":"http://localhost:20334",
  "WebSocketAddress":"http://localhost:20335",
  "WalletFile":"./wallet.dat",
  "Password":"wangbing"
}
```

Then start to run.


