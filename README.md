# md2csv

Converts an IETF standards document (RFC) into a CSV file, ignoring common sections such as the `Introduction`, `Definitions`, `Acknowledgements`, etc.

This application was written with the intention to create a quick compliance check list from an RFC document.

## Install
```bash
go install github.com/digitorus/md2csv
```

## Usage
```bash
md2csv {url}
md2csv {url} {url} {url} ...

md2csv https://raw.githubusercontent.com/mozilla/pkipolicy/master/rootstore/policy.md
md2csv https://raw.githubusercontent.com/cabforum/documents/master/docs/BR.md
md2csv https://raw.githubusercontent.com/cabforum/documents/master/docs/BR.md https://raw.githubusercontent.com/cabforum/documents/master/docs/EVG.md
```

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