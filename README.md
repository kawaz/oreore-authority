# DNS Server for oreorecert

- wildcard patterns for `*.ipv4.oreore.net`
  - `/(?:\w+-)*?(\d+)-(\d+)-(\d+)-(\d+)$/`
  - `/(?:\w+-)*?([0-9a-f]{8})$/`
  - AAAA も ::ffff:xxxx:xxxx で返す

- wildcard patterns for `*.ipv6.oreore.net`
  - `/(?:\w+-)*?(\d+)-(\d+)-(\d+)-(\d+)$/`
  - `/(?:\w+-)*?([0-9a-f]{8})$/`
  - `/^[0-9a-f]{4}-?[0-9a-f]{4}$/` は A も返す
