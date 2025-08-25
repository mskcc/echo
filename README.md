# Echo - Copy Service

The Echo is a the application used for coping files from one location to another. You can comunicate with the service using RabbitMQ or Rest API.

## Usage

```
git clone https://github.com/mskcc/echo/
cd echo
# Create configuration .env file
docker-compose up
```

### REST API

Service status
```
GET /status
Response 200 
{
    "message": "Server is running"
}
```

Copy file request
```
POST /copy
Body: {
    "source": "/path/to/source/source_file.txt", 
    "destination": "/path/to/source/source_file.txt"
}
Response 202
{
    "id": "`<request_id>`",
    "message": "File copy request accepted with id: `<request_id>`"
    "task": "COPY"
}
```

Delete file request
```
POST /delete
Body: {
    "source": "/path/to/source/source_file.txt"
}
Response 202
{
    "id": "`<request_id`",
    "message": "Delete file request accepted with id: `<request_id`",
    "task": "DELETE"
}

```

### RabbitMQ Messages

You can also submit COPY and DELETE requests

Copy Message

```
{
    "id": "`<request_id>`",
    "task": "COPY",
    "source": "/path/to/source/source_file.txt", 
    "destination": "/path/to/source/source_file.txt"
}
```

Delete Message

```
{
    "id": "`<request_id>`",
    "task": "DELETE",
    "source": "/path/to/source/source_file.txt", 
}
```

Response Message
```
{
    "id": "`<request_id>`",
    "status": `<task_status>`, # success or fail
    "message": `<message>`
}
```

## Configuration

RABBITMQ_URL=amqp://{username}:{password}@{url}:{port}/<br>
API_TOKEN={api_token}<br>
FILE_TASK_QUEUE={task_queue}<br>
CONFIRMATION_QUEUE={confirmation_queue}<br>
NUMBER_OF_WORKERS={number_of_workers}<br>
SERVER_PORT={server_port}<br>
