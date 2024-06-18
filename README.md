# khranity
Khranity is a service for easy and fast put/get objects to S3/cloud storage

## features
- works with S3 API storage
- support ignore file/list

## usage

If necessary, you can have a file `.khranityignore` by analogy `.gitignore`, for example:
```
# git
.git

# dependencies
node_modules
resources

# next.js
.next
out

# production
build

# debug
npm-debug.log*
yarn-debug.log*
yarn-error.log*
.pnpm-debug.log*

# misc
.contentlayer
```

### settings

An example of settings in `lore.yml`:
```
jobs:
  - name: khranity
    path: /root/khranity
    token: "file:khranity.token"
    schedule: "0 1 * * 0,1,2,3,4,5,6"
    tz: "Europe/Moscow"
    cloud: yandex
    archiver: native
    ignore:
      - .git
      - build
      - logs

clouds:
  - name: selectel
    method: aws
    url: "https://s3.storage.selcloud.ru"
    region: "ru-1"
    bucket: backups
    access_id:  "file:selectelID.token"
    secret_key: "file:selectelKey.token"
  - name: yandex
    method: aws
    url: "https://storage.yandexcloud.net"
    region: "ru-central1"
    bucket: backups
    access_id:  "file:yandexID.token"
    secret_key: "file:yandexKey.token"

setup:
  os: nix
```

### docker

An example of the necessary parameters to run docker container:
```
  --volume ~/.khranity/logs:/exec/logs				    \  
  --read-only --volume ~/.khranity/config:/exec/config	\
  --read-only --volume ~/:/exec/data					\
```

| folder in docker    | purpose           |
| ------------------- | ----------------- |
| /exec/logs          | logs              |
| /exec/config        | config and tokens |
| /exec/data          | data              |

## related libs

- [walle/targz](https://github.com/walle/targz)
- [mholt/archiver](https://github.com/mholt/archiver)
- [sabhiram/go-gitignore](https://github.com/sabhiram/go-gitignore)

## lisence

Licensed under MIT license. See [LICENSE](LICENSE) for more information.