### Install it
```go install github.com/blakgeek/gsmject@latest```

### Use It
```gsmject command [command args...]```

### Gsmject Some Secrets 
```shell
# expose latest version of secret from the current project as environment variable ENV_SECRET1
ENV_SECRET1=secret:SOME_SECRET

# expose version 1 of secret from the current project as environment variable ENV_SECRET2
ENV_SECRET2=secret:OTHER_SECRET:1

# expose the secret SOME_SECRET from the project some-project as environment variable ENV_SECRET3. 
# you must use a valid secret path as defined by Google
ENV_SECRET3=secret:projects/some-project/secretes/SOME_SECRET/versions/latest

# expose latest version of secret from the current project as a file at /secrets/SECRET1
MOUNTED_SECRET1=secret:SOME_SECRET@/secrets/SECRET1

# expose version 1 of secret from the current project as a file at /SECRET2/key
MOUNTED_SECRET2=secret:OTHER_SECRET:1@/SECRET2/key

# expose the secret SOME_SECRET from the project some-project as a file at /MY_SECRET. 
# must use a valid secret path as defined by Google
MOUNTED_SECRET3=secret:projects/some-project/secretes/SOME_SECRET/versions/latest@/MY_SECRET
```
