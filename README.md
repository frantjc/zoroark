# zoroark

```sh
# No auth
docker run -e DISPLAY=host.docker.internal:0 xeyes

# For auth...
# Get cookie
xauth list
$HOSTNAME:0  MIT-MAGIC-COOKIE-1  $COOKIE
...

# Add entry in container (somehow; could mount or run this command while inside)
xauth add host.docker.internal:0  MIT-MAGIC-COOKIE-1  $COOKIE
```

NO way to use graphical Steam installation's authenticated session with `steamcmd`: must prompt for re-login.
