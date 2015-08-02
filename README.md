## Periodic File Fetcher

Periodic file fetcher is a small program written in Go. It allows you to define
external resources (distant with a URL for example) in a configuration file. The
program will then parse the configuration file in which you can define an interval
for each resource and start the file fetching, periodically checking whether or not
the distant resource changed or not. If the resource changed, the old fetched file
is renamed so you can keep an history of the changes made to this resource.
