# Porter
> Ross Gallagher

## Running 
The app can be run with the following make command

```bash
make run-porter
```

## Testing
The project can be tested with the following command

```bash
make test-porter
```

## Overview
The project is implemented in 3 distinct parts:

* Parser - A generic JSON object parser exposing a `Read() (string, *T, error)` func.
  * This allows for reuse of the parser across multiple types assuming the input format is consistent 
* Store -  A simple Redis backed store allowing keyed insert
  * This is a very simple store implementation, it requires callers to correctly key their input.
* Port Service - A Port specific processor service. 
  * On initialization the processor crates a Parser with its Port type
  * On initialisation the processor creates Store to insert read ports
  * On Run the service reads from the Parser until an io.EOF error is returned.
    * It validates the read Port type (a simple validation has been implemented).
    * It writes valid Ports to the Store, using the parsed key with a prefix.
    * An optional flag can be set to ignore non EOF errors from the Parser.

The main func of the project has simple implementation for reading a file for input to the port service.
Configuration is based on environment variables.