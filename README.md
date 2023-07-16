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

### Migrations

We use [Atlas](https://atlasgo.io/) for schema based migrations.
Because we use GORM as our ORM, we can use the <https://atlasgo.io/guides/orms/gorm> package to generate migrations directly from our models.
The configuration for this happens in `atlas.hcl` and `loader/atlasGorm.go`

We wrap the most common Atlas commands in the Makefile, so that we can easily run them.
The commands are:

```bash
# Generate a new migration file based on the current models
make migrate-gen name=<migration-name>

# Generate a new empty migration file
make migrate-new name=<migration-name>

# Apply all migrations up to the latest version
make migrate-up

# Reverse all migrations down to the given version (version is the timestamp of the migration file)
make migrate-down version=<version>
```
