# BlaBlaBot

BlaBlaBot is a Telegram bot to receive alert for trips in Blablacar :)


### Installation

First of all install and configure Golang, then:

1. Clone this repository (in $GOPATH/src) and rename the folder to BlaBlaBot if it is not the name.
2. Install and configure RedisDB
3. Run the install script `install.sh` to install all the dependencies.
4. Compile using Makefile. For example: `make install` will compile and generate a binary called BlaBlaBot in **$GOPATH/bin**.
5. Configure the bot using the configure_example.json file and renaming to config.json.

Then, you can run using `-c` param to provide a configuration file if the name is not "config.json" or it isn't in the actual directory.
Example: `bin/blablabot -c path/to/config.json`

