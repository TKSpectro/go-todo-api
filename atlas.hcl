data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "-mod=mod",
    "./loader/atlasGorm.go",
  ]
}

# This environment is used for generating the migrations based on the gorm schemas
env "gorm" {
  src = data.external_schema.gorm.url
  dev = "docker://maria/latest/dev"
  migration {
    dir = "file://migrations?format=golang-migrate"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}

# This environment is used for local development. See: https://atlasgo.io/concepts/dev-database
env "local" {
    url = "maria://root:root@localhost:3307/go_api"
    dev = "docker://maria/latest/dev"
    migration {
        dir = "file://migrations?format=golang-migrate"
    }
}