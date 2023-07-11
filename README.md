# go-fiber-api

```air``` for development

```swag init``` for generating docs

```make start``` and ```make stop``` for starting/stopping development database

## Knowledge base

### JSON parsing with go

Because the BodyParser will default parse null string fields as empty strings, we need a better solution to get actual null values
With the omitempty tag, we also won't get the desired result, because it will omit the field if it's null (and then use the default value - empty string)

The solution is to use the packages `"gopkg.in/guregu/null.v4/zero"` and `"gopkg.in/guregu/null.v4/null"`
And then use the types `zero.String` and `null.String` instead of `string`

Because these will use a struct under the hood, we also want to overwrite the swagger documentation for these fields, so that it will show the correct type in the docs with `swaggertype:"string"` tag
