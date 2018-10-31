# Golang distribured nodes implementation

## Getting started
You'll need to set up the first node

In your 1st terminal ``go run main.go`` 
Example response from 1st terminal:
```
Launching server...
Using port: 35569
```
Launching 2nd terminal need to specify "ip:port" of the node you want connect to
``go run main.go -d 127.0.0.1:39237``

Response will be the same as from 1st terminal.

After connection to 1st node:
```
2018/10/31 03:21:51 Got a new stream!

```

After retreiving new data from connected node you will see:
```
{
  "1": "2",
  "4": "5",
  "key": "value"
}
```
