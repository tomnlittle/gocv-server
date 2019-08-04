# Setup

## Docker
AWS_PROFILE and AWS_REGION need to be set in your shell, these will be inherited
by the application in order to connect to S3

Then run
```sh
docker-compose -f docker-compose.yml up --build
```

# Interface

## Simple

The simple interface will fetch the image from S3 and can convert it to either
a PNG or JPEG at any compression or quality level - expressed as the optional query parameters shown below

### Example Request with format and quality parameters
```http
GET localhost:8000/v1/s3/:bucket/:key?format=jpeg&quality=100
```

## Complex
The 'complex' interface takes a PUT request and a JSON encoded body.

There are three main components of the json body.

| Key           | Required           | Example Value  |
| :------------- |:-------------:| :-----:|
| encoding      | :heavy_check_mark: | png |
| quality      | :x: | 50 |
| functions | :x: | [] |

### Example Request

```http
PUT localhost:8000/v1/s3/:bucket/:key
```

With JSON Body
```json
{
    "encoding": "jpeg",
    "quality": 100,
    "functions": [
        {
            "id": "resize",
            "parameters": [
                {
                    "key": "width",
                    "value": "1000"
                },
                {
                    "key": "height",
                    "value": "1200"
                }
            ]
        }
    ]
}
```

# Implemented Functions

| OpenCV Function           | Implemented  |
| :-------------: | :-------------: |
| Resize      | :heavy_check_mark: |
| Rotate      | :x: |
| Perspective Transform | :x: |
