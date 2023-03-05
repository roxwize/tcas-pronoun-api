# tcas-pronoun-api
RESTful API for fetching the pronouns of Two Cans and String users. Authentication will be implemented eventually; **this is not ready for release and should not be used currently in any circumstance other than for testing purposes.**
## Endpoints
### GET /pronouns
Get all pronouns stored in the database.
### GET /pronouns/{username}
Get the pronouns of a specific user in plaintext.
### PATCH /pronouns/{username}
Set the pronouns of a specific user. Accepts the `pronouns` field.
### POST /pronouns/add
Add a user. Accepts JSON data where `username` is the username of that user (without any special characters) and `pronouns` are their pronouns.
### DELETE /pronouns/{username}
Delete a user from the database.