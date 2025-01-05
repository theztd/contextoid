# Contextoid

Check if given text is still relevant to given comments and marks comments that changing context.



## Quick start

```bash
go run ./cmd/
```

## Example

Sent example data to api

```bash
curl -X POST http://localhost:8080/v1/gpt/analyze_comments \
    --header "Content-Type: application/json" \
    -d @data.json
```

data.json
```json
{
  "task_text": "Implement the API for handling user comments.",
  "comments": [
    "I think we need to add error handling.",
    "This API should support authentication.",
    "The API must integrate with a third-party service for notifications.",
    "Update the documentation to reflect recent changes.",
    "The API must handle 5000 requests per second.",
    "Can we discuss the timeline for this task? It seems unclear.",
    "This is unrelated to the API development and should be handled separately."
  ]
}
```


Response should look like this
```json
{
  "Data": {
    "comments_analysis": [
      {
        "changes_scope": true,
        "comment_id": 1
      },
      {
        "changes_scope": true,
        "comment_id": 2
      },
      {
        "changes_scope": true,
        "comment_id": 3
      },
      {
        "changes_scope": true,
        "comment_id": 4
      },
      {
        "changes_scope": true,
        "comment_id": 5
      },
      {
        "changes_scope": false,
        "comment_id": 6
      },
      {
        "changes_scope": false,
        "comment_id": 7
      }
    ],
    "new_description": "Implement the API for handling user comments. The API should support authentication, add error handling, and handle 5000 requests per second. The API must integrate with a third-party service for notifications. After all updates, the documentation needs to be updated to reflect the changes."
  },
  "Usage": {
    "completion_tokens": 213,
    "completion_tokens_details": {
      "accepted_prediction_tokens": 0,
      "audio_tokens": 0,
      "reasoning_tokens": 0,
      "rejected_prediction_tokens": 0
    },
    "prompt_tokens": 261,
    "prompt_tokens_details": {
      "audio_tokens": 0,
      "cached_tokens": 0
    },
    "total_tokens": 474
  },
  "finish_reason": "stop"
}

```

