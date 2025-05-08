data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "-mod=mod",
    "ariga.io/atlas-provider-gorm",
    "load",
    "--path", "./model",  
    "--dialect", "postgres",          
  ]
}

env "gorm" {
  src = data.external_schema.gorm.url
  dev = "postgresql://postgres:Phongsql123@localhost:5432/golang_example?sslmode=disable"  
  url = "postgresql://postgres:Phongsql123@localhost:5432/golang_example?sslmode=disable"

  migration {
    dir = "file://migrations"    
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"  
    }
  }
}
