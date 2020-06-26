# Kitsune
Simple message ingestion and distribution service.

Under development. Currently only inplements a simple in memory repository.

## Goals
Primarily focusing on simplicity in as many aspects as possible:
- Easy to use as a client.
- Easy to deploy.
- Easy to manage.
- Easy to develop and maintain.

Currently the focus is to optimize for deploying as docker container on a cluster. Think of ECS using a spot fleet in AWS. Which
means an instance can go down and up at any time and you typically run multiple instances spread over several zones. The first
consequence of this is that the storage should be decoupled from the application.  

## Publishing message (POST - /publish/{topic})
Messages are published using this simple format:
```
{
    "properties": {
        "key1": "value1",
        "key2": "value2"
    },
    "eventTime": "2020-06-26T18:59:00.266293Z",
    "payload": "Some payload"
}
```
| Field         | Description                                                                       |
| ------------- | --------------------------------------------------------------------------------- |
| properties    | Optional - Key value store. Can be any json structure.                            |
| eventTime     | Optional - Client can set a time for the event. Will not be used for ordering.    |
| paylaod       | Required - Content of the message/event.                                          |

The persisted message is returned on a successful call:
```
{
    "id": "f95400fb-fb4b-4feb-af24-8727fbdbda61",
    "publishedTime": "2020-06-26T18:59:46.045908+02:00",
    "properties": {
        "key1": "value1",
        "key2": "value2"
    },
    "eventTime": "2020-06-26T18:59:00.266293Z",
    "topic": "test",
    "payload": "Some payload"
}
```

| Field         | Description                                   |
| ------------- | ----------------------------------------------|
| id            | Unique message message id set by server.      |
| publishedTime | Time set by server when receiving the message.|
| properties    | Properties set by client.                     |
| eventTime     | Event time set by client.                     |
| topic         | Topic the message was published to.           |
| paylaod       | Payload set by client.                        |

## Fetching individual message (GET - /{topic}/{messageId})
An individual message can be retrieved by performing a GET request with the topic and the message id.

## Polling messages (POST - poll/{topic})
PollRequest:
```
{
	"subscriptionName": "test",
	"topicName": "test",
	"maxNumberOfMessages": 3
}
```
| Field                 | Description                                   |
| --------------------- | ----------------------------------------------|
| subscriptionName      | Name of subscription.                         |
| topicName             | Topic to get messages from.                   |
| maxNumberOfMessages   | Maximum number os messages to be returned.    |

The return value is an array of messages:
```
[
    {
        "id": "7e8bd90d-a833-43bf-8f67-5487f4a909c7",
        "publishedTime": "2020-06-26T17:22:57.318603Z",
        "properties": {
            "key1": "value1",
            "key2": "value2"
        },
        "eventTime": "2020-06-26T18:59:00.266293Z",
        "topic": "test",
        "payload": "Some payload"
    },
    {
        "id": "96dc64be-ef87-4a16-bee1-a2d8ab4835be",
        "publishedTime": "2020-06-26T17:22:59.85086Z",
        "properties": {
            "key1": "value1",
            "key2": "value2"
        },
        "eventTime": "2020-06-26T18:59:00.266293Z",
        "topic": "test",
        "payload": "Some payload"
    },
    {
        "id": "8511210c-9afb-4d8f-987a-d810f12d42bc",
        "publishedTime": "2020-06-26T17:23:00.761936Z",
        "properties": {
            "key1": "value1",
            "key2": "value2"
        },
        "eventTime": "2020-06-26T18:59:00.266293Z",
        "topic": "test",
        "payload": "Some payload"
    }
]
```
What messages are returned are determined on a topic + subscription basis. A topic can have multiple subscribers which are
maintained seperately. Multiple instances of an application can share a subscription and process different messages. Make sure
different applications dont have the same topic + subscription combination.

## Settings
A subscription can be set to start at a specific message. Can be set by either a message id or a time. 
```
{
	"subscriptionName": "test",
	"publishedTime": "2020-06-26T18:59:00.266293Z",
	"messageId": "someId"
}
```
| Field                 | Description                                   |
| --------------------- | ----------------------------------------------|
| subscriptionName      | Name of subscription.                         |
| publishedTime         | Set the subscription back to a specific time. |
| messageId             | Set subscription back to a specific message.  |

When setting the subscription back to a specific message, that message will be the first to be returned when polling. If time is
used the message with the nearest which is not before is picked.

## Planned features
- **Ack/nack functionality**: Polled messages should be unavailable for polling for a set amount of time and be returned in later
calls to poll unless nacked in the span of the set time.
- Actual usable persistance. What will be the primary (first) one to be implemented is not decided.
- Implement a max size for returned result. Want to support "large" messages.
- Authorization.
