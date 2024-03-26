## Setta-auth is a JWT-microservice for [Setta](https://github.com/rosenthall/setta-core) app.

### Basic info

It provides GRPC api, that described in `protos/auth.proto`

Available methods are :  token validating, refreshing, generating and data extracting.
Service stores the JWT refresh-sessions via Redis.  
RS256 signing alghoritm(sha256 with RSA 2048 bits keys) is used to provide better security.

### Libraries stack

The official [GRPC](https://github.com/grpc/grpc-go/tree/master) library and [grpc-middleware](https://github.com/grpc-ecosystem/go-grpc-middleware/tree/main) to create grpc server
[Viper](https://www.notion.so/23-12-13-738f6c3280d94fc7a3ad1c7270b35f30?pvs=21) framework to provide simple project configuration in toml file

[Zap](https://github.com/uber-go/zap) for logging, [redis](http://github.com/go-redis/redis/v8) for interaction with redis db

[testify/mock](http://github.com/stretchr/testify/mock) for creating mock objects for testing

[golang-jwt/jwt](golang-jwt/jwt) as library for working with jwt :>
