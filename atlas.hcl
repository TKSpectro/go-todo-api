data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "-mod=mod",
    "./loader/atlasGorm.go",
  ]
}

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

env "local" {
    url = "maria://root:root@localhost:3307/go_api"
    dev = "docker://maria/latest/dev"
    migration {
        dir = "file://migrations?format=golang-migrate"
    }
}