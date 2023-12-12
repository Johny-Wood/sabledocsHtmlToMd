## Installation
```
go get github.com/Johny-Wood/sabledocsHtmlToMd
```
## Usage

In the folder with sabledocs output files, run:
```
sabledocsHtmlToMd
```
## Configuration

Create config.toml where the package will be executed (where your sabledocs are). Use config.example.toml as an example for configuring your sabledocsHtmlToMd.

### Options

#### [Settings] ExcludeInputFiles

List the files you don't want to be converted to md files.  For instance, sabledocs builds index.html and search.html that you won't need to be converted. You can easily ecxlude them from conversion:
```
[Settings]
ExcludeInputFiles = [ "index.html", "search.html" ]
```

#### [Translation.TablesT]
Your translations for table headers (Field, Type, Description).
```
[Translation.TablesT]
Field = "Поле"
Type = "Тип"
Description = "Описание"
```

#### [Translation.ReqResT]
Your translations for request & response words in methods (Request, Response).
```
[Translation.ReqResT]
Request = "Запрос"
Response = "Ответ"
```

#### [Translation.EntitiesT]
Your translations for entities (Message AddMesasage).
```
[Translation.EntitiesT]
Message = "Сообщение"
Service = "Сервис"
```
