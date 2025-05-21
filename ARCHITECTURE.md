# Stragollum

This project is a client for a HTTP API, designed
to make user code able to interact with the "Astra DB Data API".

This package is being designed and created as we go, so here for now is a tentative basic architectural plan and an action plan (an ordered list of TODO items).

## Architectural plan

User code would always begin by creating an instance of class `DataAPIClient`.

A `DataAPIClient` has a `GetDatabase` method to spawn a `Database` instance.

Database methods (e.g. `ListCollectionNames`) send a request with a certain JSON payload and returns
(a manipulated) response to its caller.

For convenience, since these requests share certain properties,a helper `DataAPICommander` class is created,
not for direct use by the user code.

The `Database` needs a `CreateCollection` method. Its parameters are `name` (string) and `definition`, and for now it will return nothing. The `definition` is an object of a new struct type, `CollectionDefinition`, which must capture the structure in the following example (but all fields are optional!):
```
{
    "defaultId": {
        "type": "string"
    },
    "indexing": {"k": "v"},
    "lexical": {
        "analyzer": "string",
        "enabled": true
    },
    "rerank": {
        "enabled": true,
        "service": {
            "authentication": {"k": "v"},
            "modelName": "string",
            "parameters": {"k": "v"},
            "provider": "string"
        }
    },
    "vector": {
        "dimension": 999,
        "metric": "string",
        "service": {
            "authentication": {"k": "v"},
            "modelName": "string",
            "parameters": {"k": "v"},
            "provider": "string"
        },
        "sourceModel": "string"
    }
}
```
In the above, `{"k": "v"}` represents an arbitrary key-value mapping (with string keys).

The `createCollection` method, with parameters `name` and `definition` outlined above, must create the full payload for a POST request to the Data API (done through the APICommander just like listCollectionNames).
The payload has the form: `{"createCollection": {"name": <name>, "options": <definition>}}`; the response must simply be: `{"status": {"ok": 1}}` (if it's not like this, an error must be raised).

## Action plan

(Done items are marked with [X] in the following)

1. [X] Create a skeletal DataAPIClient class and a simple unit test for it
2. [X] Write a draft for class `Database`
3. [X] Add `GetDatabase`, with appropriate tests to check the parameter behaviour under all conditions
4. [X] Create the APICommander class (with getters out of convenience) and simple unit tests
5. [X] Hook the APICommander instance creation into the constructor of Database, enriching the Database unit tests accordingly
6. [X] Add an actually-working `Request` method to the APICommander, which sends a POST request to the URL and returns the response (if successful). This may require new dependencies to handle http requests.
7. [X] Work on the `ListCollectionNames` method of `Database`
8. [X] Add an integration test that simply creates a client->database and runs ListCollectionNames on a real, actual database.
9. [X] Prepare, in a separate source file, the `CollectionDefinition` type (in the future, special constructors and helpers will also be there)
10. [X] Add the `Database.CreateCollection` method as described in the specs.