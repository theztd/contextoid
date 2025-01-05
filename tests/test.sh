#!/bin/sh


curl -X POST http://localhost:8080/v1/gpt/analyze_comments \
    --header "Content-Type: application/json" \
    -d @test.json