
# Dolus Mock Server

## Overview

Dolus is a mock server designed to simulate HTTP-based APIs. It's useful in testing scenarios where you need to mimic the behavior of external services. Dolus allows you to define expectations for incoming requests and specify the responses that should be returned.

## Key Features

- **Expectation Engine**: Dolus uses an expectation engine to match incoming requests against a set of predefined expectations. Each expectation can specify the HTTP method, path, headers, and query parameters of the request.
- **Dynamic Response Generation**: Dolus can generate dynamic responses based on the incoming request and the matched expectation. This includes generating random data for certain fields in the response body.
- **Schema Generation**: Dolus includes functionality for generating schemas from different types of input, including OpenAPI 3 examples and arbitrary data structures.
- **Logging**: Dolus includes a logging package for detailed logging of the server's operations.

## How to Use

1. Define your expectations in a `.cue` file. Each expectation should be an object that includes the following fields:

    - `priority`: The priority of the expectation. If multiple expectations match a request, the one with the highest priority is used.
    - `request`: An object that describes the request. This includes the following fields:
        - `path`: The path of the request.
        - `method`: The HTTP method of the request.
        - `headers`: An object that maps header names to expected values. The values can be either strings or arrays of strings.
        - `params`: An object that includes `query` and `path` fields. Both of these are objects that map parameter names to expected values.
    - `response`: An object that describes the response. This includes the following fields:
        - `body`: An object that describes the body of the response. This can include fields that are generated dynamically, such as `id: task.#GenInt32 & {min: 80, max: 100}` which generates a random `int32` number between 80 and 100 for the `id` field.
        - `status`: The HTTP status code of the response.

    Here's an example:

    ```cue
    expectation.#Expectation & {
        priority: 1
        request: expectation.#Request & {
            path:   "/store/order/{orderId}"
            method: "GET"
            headers: {
                "Content-Type": ["application/json"]
            }
            params: {
                query: {
                    value: [3,7]
                }
                path: {
                    orderId: 3
                }
            }
        }
        response: expectation.#Response & {
            body: {
                petId: {
                    id: task.#GenInt32 & {min: 80, max: 100}
                }
                complete: true
                status:   "good day"
            }
            status: 200
        }
    }
    ```

2. Pass your `.cue` file into the Dolus server. The server will parse the `.cue` file and load the expectations.

3. Start the Dolus server. The server will listen for incoming requests and match them against the predefined expectations. If a match is found, the server will return the corresponding response.

---

Please note that this is a basic outline and might not cover all aspects of your project. You should expand upon this outline with more specific information about your project, such as how to install and run it, how to contribute, and any other relevant details.