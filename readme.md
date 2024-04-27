# Carbide Registry API

The database of hardened images for rancher government carbide resides in our secured cloud. This stateless API should act as a simple interface to such.

## Build
```bash
# compile the binary
make
# for more options try
make help
```
## Develop
```bash
# This will start the api and hot reload such whenever changes are saved to a .go file.
./reload.sh
```
## Deploy
Bare-metal, docker, and k8s deployments are all supported - though k8s is recommended.
```bash
helm install carbide-registry-api ./chart --values <values-file>
```
## Testing
For now an insomnia file can be found at `./example/exampleQueries.json`.
If you prefer to use the OSS version try [insomnium](https://github.com/ArchGPT/insomnium).
If you prefer to use cURL or a programming language to test the endpoints the queries can be exported to those as well.

If you don't have access to a mysql database available for testing, you can deploy one to your local cluster with the [mysql operator](https://github.com/mysql/mysql-operator)
You can either follow their instructions or run `./bin/mysql.sh` which should deploy their operator and create a test cluster.

## Environment
| Variable              | Description               | Optional                  |
| --------------------  | -----------               | --------                  |
| DBUSER                | MySQL/MariaDB username    | false                     |
| DBPASS                | MySQL/MariaDB password    | false                     |
| DBHOST                | MySQL/MariaDB host        | false                     |
| DBPORT                | MySQL/MariaDB port        | false                     |
| DBNAME                | MySQL/MariaDB name        | false                     |
| PORT                  | port to serve api         | true (defaults to 5000)   |
| JWTSECRET             | JWT privatekey            | false                     |

## REST Schema
- /user
    - POST: accepts username and password
        - creates user account
        - only returns cookie first time user is created
    - DELETE: accepts username and password
        - deletes user account
- /login
    - POST: should accept username and password
        - returns authentication cookie/token (stored in browser by frontend)  
>The following require the user to have provide their auth token via cookie:
- /product
    - GET: get all products
    - POST: create new product
- /product/{product_name}
    - GET: get product
    - PUT: update product
    - DELETE: delete product
- /product/{product_name}/release
    - GET: get all releases for product
    - POST: create new release for product
- /product/{product_name}/release/{release_name}
    - GET: get release
    - PUT: update release
    - DELETE: delete release
- /image
    - GET: get all images
    - POST: create new image
- /image/{image_id}
    - GET: get image
    - PUT: update image
    - DELETE: delete image
- /releaseImageMapping
    - GET: get all releaseimgmappings
    - POST: create new releaseimgmapping
    - DELETE: delete releaseimgmapping
#### Query Parameters
The following query parameters are supported:
| Parameter           | Description            | Default  | Max  |
| ------------------- | -----------            | -------- | ---  |
| page                | page of data to return | 1        | N/A  |
| pageSize            | size of data page      | 10       | 50   |
> Pagination is supported on all endpoints that return multiple objects. These are: `/product`, `/product/{product_name}/release`, `/image`, `/releaseImageMapping`

## Misc TODO notes
- move entire DB schema to this api (for portability)
- product names should be unique

#### Possible Endpoint Additions:
- /product/{product_name}/release/{release_name}/image
    - GET: get all images for product release
    - POST: create new image for product release
- /product/{product_name}/release/{release_name}/image/{image_name}
    - GET: get image from product release
    - PUT: update image from product release
    - DELETE: delete image from product release
- /release
    - GET: get all releases
    - POST: create new release
- /release/{release_name}
    - GET: get release
    - PUT: update release
    - DELETE: delete release

