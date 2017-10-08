HockeyApp webhook script that filters and posts every new app build to OK.ru/TamTam chat via [OK Bot API](https://apiok.ru/dev/graph_api/bot_api).

Example of docker build & run:
`docker build -t ok-hockeyapp .`
`docker run -d -p 8111:8080 ok-hockeyapp`
