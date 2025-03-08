#!/bin/bash

# Empty directories cannot be commited to source control, so we need to make sure these directories exist
scriptDir=$(dirname "$(realpath $0)")
mkdir -p $scriptDir/expected/command1/empty_dir
mkdir -p $scriptDir/expected/childCommand1/empty_dir



# Clear out anything that has been created by a previous scaffolding test
ignoredNames="scaff|scaff.exe|scaff.json|scaff_files|my_templates"
scaffoldingRunDir=$scriptDir/environment/child_dir/grandchild_dir

itemPaths=($(find $scaffoldingRunDir -mindepth 1 -maxdepth 1 -print))

for itemPath in "${itemPaths[@]}"; do
    itemNameOnly=$(basename "$itemPath")  # Extract just the file/directory name
    
    if [[ ! "$itemNameOnly" =~ $ignoredNames ]]; then
        rm -rf $itemPath
    fi
done
