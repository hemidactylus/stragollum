# Stragollum

This project is a client for a HTTP API, designed
to make user code able to interact with the "Astra DB Data API".

This package is being designed and created as we go, so here for now is a tentative basic architectural plan and an action plan (an ordered list of TODO items).

## Architectural plan

User code would always begin by creating an instance of class `DataAPIClient`.
When doing so, _optionally_, an "environment" (one of an enum DEV/TEST/PROD) is specified; also optionally an authentication token (a string) can be supplied.

A `DataAPIClient` has a `GetDatabase` method to spawn a `Database` instance.
This method has a required parameter, `apiEndpoint` (string), and an optional `token` which if supplied replaces the DataAPIClient setting.
Another optional parameter is `keyspace`, which if omitted has a hardcoded default.

Databases have methods that trigger HTTP POST requests to the Data API, with a certain payload,
and parse the response to construct the method return value. For convenience, since these requests 
require certain headers and share a common path prefix, a helper `DataAPICommander` class is created,
not for direct use by the user code.

The DataAPICommander constructor accepts the URL (required string) and the token (optional) parameters.
- The commander URL is built by the Database and is in the form "<api endpoint>/api/json/v1/<keyspace>"
- If token is provided, it translates into a "Token: <token>" header to POST requests (no headers oterwise)
- The Database creates an instance of the APICommander and keeps it ready to use each time it needs to perform a request.

## Action plan

(Done items are marked with [X] in the following)

1. [X] Create a skeletal DataAPIClient class and a simple unit test for it
2. [X] Write a draft for class `Database`
3. [X] Add `GetDatabase`, with appropriate tests to check the parameter behaviour under all conditions
4. [ ] Create the APICommander class (with getters out of convenience) and simple unit tests
5. [ ] Hook the APICommander instance creation into the constructor of Database, enriching the Database unit tests accordingly