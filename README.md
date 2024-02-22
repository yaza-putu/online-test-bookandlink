# Online Test Bookandlink
![golang](https://upload.wikimedia.org/wikipedia/commons/thumb/0/05/Go_Logo_Blue.svg/1200px-Go_Logo_Blue.svg.png)

Project layout design influenced by [standard go project layout](https://github.com/golang-standards/project-layout)
## How to start

- install depedency
  ```bash
  make tidy
  # or
  go mod tidy
  ```
- make environment 
  ```bash
  make config
  #or
  cp .env.example .env
  ```

- generate key
  ```bash
    make key
  ```
  note : don't forget copy key to .env


- run migration
  ```bash
  make migrate-up
  ```

- run dev mode
  ```bash
    make run
  ```
- build
  ```bash
  make build
  ```

- run test
  ```bash
   make gotest
  ```

## Command Available
- make migration
  ```bash
   make migration table="name_of_table"
  ```
  
- run migration
  ```bash
   make migrate-up
  ```

## Validation Unique With Struct Tag
- unique
```go
type v struct {
	Name string `validate:"unique=table_name:column_name"`
}
// ecample
type v struct {
Name string `validate:"unique=users:name"`
}
```
- unique with ignore
```go
type v struct {
Name string `validate:"unique=table_name:column_name:ignore_with_field_name"`
ID   string `validate:"required"`
}
// example
type v struct {
Name string `validate:"unique=users:name:ID"`
ID   string `validate:"required" json:"id"`
}
```
## Stack 
- [Echo](https://echo.labstack.com)
- [Gorm](https://gorm.io)
- [Env](https://github.com/spf13/viper)
- [Redis](https://github.com/redis/go-redis)

