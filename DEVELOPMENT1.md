# Working on the SCAFF Project

## Packages

Each package in this project has a _doc.go_ file. This file will provide a description of the purpose of that package. 

A lot of the package directories also contain files for their corresponding test packages (providing the unit tests for that package).

## Testing

There are 2 sets of tests in this project:
 - Unit Tests -- These can be run with the `make test-unit` command. These tests can be found in the corresponding package's directory.
 - E2E Tests -- These can be run with the `make test-e2e` command. These tests can be found in the _e2e_ directory.

## Building the application

There are 2 build commands that can be used in the project:
 - `make build-dev` -- This builds the application with the debug information included
 - `make build-prod` -- This builds the application without the debug information