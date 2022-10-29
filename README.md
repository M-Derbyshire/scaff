# SCAFF

SCAFF (shortening of "SCAFFold") is a command-line tool that allows you to predefine file/folder structures, and then generate these structures in your current working directory (whilst also utilising variable tags to populate the file/folder names -- and the file contents -- with custom values).

## Installing SCAFF:

In the `distributions` directory, you can find the executable for your operating system. You will need to place this executable in a directory of your choice on your machine, and then add that directory to your PATH environment variable. If you don't know how to do edit your PATH variable, see one of the below articles:

 - [Windows](https://www.computerhope.com/issues/ch000549.htm)
 - [Linux](https://www.howtogeek.com/658904/how-to-add-a-directory-to-your-path-in-linux/)
 - [Mac OS](https://osxdaily.com/2014/08/14/add-new-path-to-path-command-line/)

If the correct executable for your operating system isn't available in the `distributions` directory, you will need to compile this yourself (using a Go compiler, and the project's makefile).

## Using SCAFF:

SCAFF creates files/directories in your current working directory, based on the "command" you have called.

When you call a "command", SCAFF will start to move up the directory tree (starting from your current working directory, and ending at the root of the current drive), looking for *scaffconfig.json* files. When it finds one of these files in a directory, it will open it and check to see if it contains the requested command. If it doesn't, the process will continue until the command is found.

So, if you were working on a project with a group of people, the root of your repository could contain a *scaffconfig.json* file (and a templates directory) with commands specific to that project. Then, further up the directory tree, you may have another *scaffconfig.json* file (say, in your user directory), that contains your personal commands.

### Example call to a SCAFF command:

`scaff my_command var1=my_value var2="my longer value"`

Here, `my_command` is the name of the command you want to execute. `var1=my_value` declares a variable named "var1", with the value "my_value". `var2="my longer value"` declares a variable named "var2", with the value "my longer value".

The variables can be used in file/directory names, and also in file templates, via tags. If a variable is required, but not provided, SCAFF will prompt the user to provide it.

### Using SCAFF variable tags:

Your file/directory names (and the templates used to generate file contents) can contain "tags" that SCAFF will replace with variable values. Below is an example:

`{: var1 :}` - This refers to a variable named "var1".

Variable tags start with "{:", and end with ":}". If you want to escape a tag, you can do so by replacing the opening with "{\\:".

### Setting up SCAFF commands:

A *scaffconfig.json* file contains a JSON object, with 2 properties:
 - `commands` is an array of command objects.
 - `templateDirectoryPath` is the path to the directory that contains the file templates for these commands (this path is relative to the location of the *scaffconfig.json* file).

Each command object has 2 properties:
 - `name` is the name of the command.
 - `directories` is an array of directory objects.

Each directory object has 3 properties:
 - `name` is the name that the directory should be created with. This can contain variable tags.
 - `directories` is an array of directory objects.
 - `files` is an array of file objects.

Each file object has 2 properties:
 - `name` is the filename (including file extension) that the file should be created with. This can contain variable tags.
 - `templatePath` is the path to the template for this file (this path is relative to the `templateDirectoryPath`, defined in the config object).

#### Example *scaffconfig.json* file:

```
{
    "templateDirectoryPath": "my_templates/some_templates",
    "commands": [
        {
            "name": "cmd1",
            "directories": [
                {
                    "name": "{: var1 :}_dir",
                    "directories": [
                        {
                            "name": "my_empty_dir",
                            "directories": [],
                            "files": []
                        }
                    ],
                    "files": [
                        {
                            "name": "my_{:var1:}_file.txt",
                            "templatePath": "fileTemplate1.txt"
                        },
                        {
                            "name": "my_{:var2:}_file.txt",
                            "templatePath": "fileTemplate2.txt"
                        }
                    ]
                }
            ]
        },
        {
            "name": "cmd2",
            "directories": [
                {
                    "name": "empty_dir",
                    "directories": [],
                    "files": []
                }
            ]
        }
    ]
}
```

Executing the `cmd1` command in the above file will generate the below files/directories in your current working directory (where var1="val1" and var2="val2"):

 - `./val1_dir`
 - `./val1_dir/my_empty_dir`
 - `./val1_dir/my_val1_file.txt`
 - `./val1_dir/my_val2_file.txt`

The 2 files will be populated with the below templates (if the *scaffconfig.json* file was located in `C:/stuff`):

- `C:/stuff/my_templates/some_templates/fileTemplate1.txt`
- `C:/stuff/my_templates/some_templates/fileTemplate2.txt`

### File templates:

File templates are simply text files, but their contents can include variable tags.

```
This is an example file template.
My name is {: user_name :} and my age is {: user_age :}
My favorite ice cream is {:favorite_ice_cream:}
```

[My Twitter: @mattdarbs](http://twitter.com/mattdarbs)  
[My Portfolio](http://md-developer.uk)
