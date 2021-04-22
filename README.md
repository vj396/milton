# milton
This is Milton, it is a slack bot.

## Local Testing
### Build Docker Container for the bot.
```
docker build -t milton .
```
### Bring Up Stack
```
docker-compose up -d
```

### Teardown Stack
```
docker-compose down
docker volume rm milton_db_data
```
