docker stop bothoi
docker rm bothoi
docker run -d --name bothoi --restart unless-stopped -v D:\database:/usr/local/bin/database bothoi
