# md2csv

Converts a Markdown document into a CSV file based on the document index, ignoring common sections such as the `Introduction`, `Definitions`, `Acknowledgements`, etc.

This application was written with the intention to create a quick compliance check list from requirements written in markdown such as the [CA/Browser Forum Documents](https://github.com/cabforum/documents) and the [Mozilla PKI Policy](https://github.com/mozilla/pkipolicy/).

## Install
```bash
go install github.com/digitorus/md2csv
```

## Usage
```bash
md2csv {url|filename}
md2csv {url|filename} {url|filename} {url|filename} ...

md2csv document.md
md2csv https://raw.githubusercontent.com/mozilla/pkipolicy/master/rootstore/policy.md
md2csv https://raw.githubusercontent.com/mozilla/pkipolicy/master/rootstore/policy.md document.md
md2csv https://raw.githubusercontent.com/cabforum/documents/master/docs/BR.md
md2csv https://raw.githubusercontent.com/cabforum/documents/master/docs/BR.md https://raw.githubusercontent.com/cabforum/documents/master/docs/EVG.md
```

Using [pdf2md](http://pdf2md.morethan.io/) you can convert a PDF document to Markdown so that you can use it with `md2csv`. 

## Using Docker
```bash
docker run -ti digitorus/md2csv sh
```
After this you can use the `md2csv` command from above.

## Reconnect Docker image
```bash
docker ps -a
docker start {container id}
docker attach {container id}
```