# JustKV

JustKV is Just a Key-Value store. Nothing more, nothing less. 

## Functions:

- `/put/<key>`: PUT the body of the request to the key.
- `/get/<key>`: GET the value of the key.
- `/delete/<key>`: DELETE the key.

## Usage:
```
go build -o justkv
./justkv
```

you can also set the port by setting the `PORT` environment variable.
```
PORT=8080 ./justkv
```

## Extensions:

You can install extensions to JustKV by adding them at compile time via tags. For example, the 
ttl extension can be added by running the following command:

```
go build -tags "ttl" -o justkv
```

The ttl extension adds the following functions:

- `/put/<key>?ttl=<duration>`: PUT the body of the request to the key with a time-to-live.

