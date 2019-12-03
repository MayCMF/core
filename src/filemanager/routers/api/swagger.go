/*
Package api "Generate swagger document"

For document rules, please refer to: https://github.com/swaggo/swag#declarative-comments-format

How to use:

	go get -u github.com/swaggo/swag/cmd/swag
	swag init -g ./src/account/routers/api/swagger.go -o ./docs/swagger*/

package api

// @title MayCMF
// @version 0.1.0
// @description Serverless CMF with Full Rest API and RBAC(Role Base Control Access) System.
// @schemes http https
// @host 127.0.0.1:8088
// @basePath /
// @contact.name eneus
