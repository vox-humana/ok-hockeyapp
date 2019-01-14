HockeyApp webhook script that filters and posts every new app build to OK.ru/TamTam chat via [TamTam BOT API](https://dev.tamtam.chat).

Example of docker build & run:
`docker build -t ok-hockeyapp .`
`docker run -d -p 8111:8080 ok-hockeyapp`
