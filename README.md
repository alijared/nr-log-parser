NR Log Parser
=============

Currently there is only 1 function, search, which will search for lines matching
criteria, and then write those lines to a specified output file.

Usage
-----

### Search ###

It is possible to match log lines using at least one of:
* log level: 
    This will match lines with a specific log level, e.g.
    ```shell script
    nrlp search -f newrelic-infra.log --level error
    ```
* component:
    This will match lines with a specific component, e.g.
    ```shell script
    nrlp search -f newrelic-infra.log --component Config
    ```
* custom:
    Custom queries are look at substrings within the line. They can be
    comma separated to search for multiple substrings, e.g.
    ```shell script
    nrlp search -f newrelic-infra.log --custom 'msg="Using Network Interface Filter Set",filter=index-1,interface=tun' 
    ```
* time filter:
    Time filter queries can be used to only evaluate log lines generated before,
    after, or before + after (in between) a specified UTC time, e.g.
    ```shell script
    nrlp search -f newrelic-infra.log --before "2019-10-18 15:05:00" --after "2019-10-18 15:00:00"
    ```
  
You can use multiple filters in order to find very specific logs.
