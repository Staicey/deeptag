# DeepTag
[![Docker](https://img.shields.io/badge/docker-ghcr.io%2Fstaicey%2Fdeeptag-blue)](https://ghcr.io/staicey/deeptag)

Middleware for linking [deepcrate](https://github.com/jordojordo/deepcrate) & [wrtag](https://github.com/sentriz/wrtag) instances together 

## Features

- **Wishlist Management**: Re-queues failed downloads back into your deepcrate wishlist
- **wrtagweb Integration**: Allows deepcrate to call wrtagweb for media tagging

## Quick Start

### Prerequisites

- [deepcrate](https://github.com/jordojordo/deepcrate) running with API enabled
- [wrtagweb](https://github.com/sentriz/wrtag) running with API enabled

### Installation 

```yaml
services:
  deeptag:
    image: ghcr.io/staicey/deeptag
    container_name: deeptag
    environment:
      - WRTAG_OP="move"      # type of operation to trigger (move|copy|reflink)
      - WRTAG_URL=           # full URL for wrtagweb (e.g. http(s)://wrtag.local)
      - WRTAG_TOKEN=         # wrtag api key
      - WRTAG_DL_DIR=        # path where unorganised downloads are stored
      - DEEPCRATE_URL=       # full URL for deepcrate (e.g. http(s)://deepcrate.local)
      - DEEPCRATE_USER=      # username/password to use for deepcrate wishlist management
      - DEEPCRATE_PASS=
    expose:
      - 8080
    restart: unless-stopped
```

