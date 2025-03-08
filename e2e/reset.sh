#!/bin/bash

# Empty directories cannot be commited to source control, so we need to make sure these directories exists
scriptDir=$(dirname "$(realpath $0)")
mkdir -p $scriptDir/expected/command1/empty_dir
mkdir -p $scriptDir/expected/childCommand1/empty_dir