# kincong

Утилита для загрузки маршрутов для соединения в роутеры Keenetic по rci интерфейсу

## Возможности
- Загрузка маршрутов из конфигурационного файла по домену, IP-адресу, IP-подсети, ASN
- Просмотр текущих загруженных маршрутов
- Просмотр доступных соединений, через которые можно проложить маршрут

## Установка

```bash
go install github.com/kincong/kincong@latest
```

## Использование

```bash
kincong --action.upload -f <config_file>
```

## Справочник параметров

```bash
Application Options:
      --debug                enable debug mode [$DEBUG]
  -i, --interface=           select interface
  -f, --file=                select file
      --dns=                 use for custom dns resolver

Main Actions:
      --action.upload        upload file params to rci (default action)

File Actions:
      --action.file.inspect  validate file
      --action.file.preview  preview file

RCI Actions:
      --action.rci.list      list avaliable rci interfaces
      --action.rci.inspect   list avaliable static routes
      --action.rci.clear     clear all static routes

RCI Params:
      --rci.host=            RCI host (default: 192.168.1.1) [$RCI_HOST]
      --rci.cookie-name=     RCI cookie name [$RCI_COOKIE_NAME]
      --rci.cookie-value=    RCI cookie value [$RCI_COOKIE_VALUE]
```

## Пример файла конфигурации

```yaml
- title: Youtube
  values:
    - AS11344
    - AS36040
    - AS36561
    - AS43515
- title: Cloudflare
  values:
    - 103.21.244.0/22
    - 103.22.200.0/22
    - 103.31.4.0/22
    - 104.16.0.0/13
    - 104.24.0.0/14
    - 108.162.192.0/18
    - 131.0.72.0/22
    - 141.101.64.0/18
    - 162.158.0.0/15
    - 172.64.0.0/13
    - 173.245.48.0/20
    - 188.114.96.0/20
    - 190.93.240.0/20
    - 197.234.240.0/22
    - 198.41.128.0/17
- title: Rutracker
  values:
    - https://rutracker.org
```
