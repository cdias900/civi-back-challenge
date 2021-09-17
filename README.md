# civi-back-challenge

API URL: [civi-back-challenge.herokuapp.com](https://civi-back-challenge.herokuapp.com/)

## Routes

- `GET /read` : Returns messages

    Query Params:
    - `page`: number (Optional). Each page contains up to  20 messages, if omitted, defaults to 1.
    
        Example: `/read?page=2`

    Response body format:
    ```json
    {
        "messages": [
            {
                "id": number,
                "timestamp": number,
                "content": {
                    "subject": string,
                    "detail": string
                }
            }
        ]
    }
    ```
---
- `POST /send` : Creates a new message

    Request body format:
    ```json
    {
        "subject": string,
        "detail": string,
    }
    ```
