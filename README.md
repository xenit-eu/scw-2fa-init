# scw-2fa-init
A tool to initialize the Scaleway CLI with a short lived API key using 2 factor authentication.

The tool will ask your email, password, token, organization and the duration of the API Key. All parameters except password and token can be provided as parameters.


## Installation

As a prerequisite you will first need to install the [Scaleway CLI](https://github.com/scaleway/scaleway-cli).
To run download a binary found under [releases](https://github.com/thijslemmens/scw-2fa-init/releases).


## Usage

```
$ ./scw-2fa-init -h
Usage of ./scw-2fa-init:
  -duration int
        Duration for the API key (in hours, max 8)
  -email string
        Email for the scaleway client
  -org string
        Organization for the scaleway client
```
