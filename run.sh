
#!/bin/bash
CURL='/usr/bin/curl'
URL='https://verify.rest/ping'
RESP="$($CURL $CURLARGS $URL)"
if [ $RESP = "Pong" ]
then
  echo "Application is running"
else
  cd /home/verify/public_html && make start
fi
