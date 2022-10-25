docker stop bothoi
docker rm bothoi
docker run -d --name bothoi --restart unless-stopped bothoi
