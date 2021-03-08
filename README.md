## Gold Rush

### View documentation
docker run -d \
    -p 8080:8080 \
    --name gr-swagger \
    -e SWAGGER_JSON=/goldrush/swagger.yaml \
    -v `pwd`:/goldrush \
    swaggerapi/swagger-ui

