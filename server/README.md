# API for qualification

## Useful links
- All the steps to create your micro-service : https://www.notion.so/side/Steps-to-create-a-micro-service-in-Golang-b40d35b2abc94d7786ece9629a4336ab

## Installation

- Get all references to keyword: `TOSET` and replace it with your project's name

```
  ag TOSET
```

- Go in the settings of your gitlab project. Click on Pipeline & add a new variable (unprotected is OK) with key: `PROJECT_NAME`, and value: `YOUR_SERVICE_NAME`

- If you want to have your own db, don't forget to replace the `MONGO_DB_NAME` and to set environment vars

- Careful, in init.go you'll have to chose between two ways of authentication. By default, it will be the `apiKey` solution

- You should now be able to set your first routes in `./src/routes` and call them from `./src/main.go`

### Run

- Get all the golang packages (indispensable to start the API):
```
  make vendor_get
```

- Launch the API on qualification.jumperX.side :
```
	make run
```

- Launch the API in local (localhost:8080):
```
	make && ./bin/qualification
```
