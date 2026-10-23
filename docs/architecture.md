# Architecture

Doctra LIS middleware sits between laboratory analyzers and the Doctra LIS (a 1C system).
It is a Go application with two Cobra commands (`cmd/`):

- **device server** (`app/device_server.go`): a TCP listener that talks to analyzers.
- **api server** (`app/api_server.go`): a Fiber HTTP API over the stored data (devices, device models, raw data, lab data, users).

Both share the same GORM store (`store/`) and config files in `config/*.json`.

## Device server flow

```
analyzer  --TCP-->  tcp.TCP  --RcvData-->  app.manageMessage  -->  driver.Driver  -->  store
                                                                      |
                                                                      +--> services.DeviceQueryService --HTTP--> LIS query host
```

1. `tcp.AcceptConnections` accepts a socket per analyzer and starts `ReadMessages` for it.
   Each read is pushed on `RcvChannel`; a `sync.WaitGroup` per connection blocks the next read
   until the message has been fully processed, so a driver may read from the socket itself
   (for example to wait for an ACK) without racing the reader.
2. `app.manageMessage` looks the device up by its IP (`devices.net_address`), creates the driver
   named by the device model (`driver.NewDriver`), and calls `processDeviceMessage`.
3. `processDeviceMessage` strips characters listed by `Driver.DataToBeReplaced`, answers ACKs,
   and uses `driver.GetRawDatas` to accumulate frames between the driver's start and end
   strings (ENQ/EOT for ASTM). Partial data is kept per connection in `tcp.PrevData`.
4. A complete raw message is stored in `raw_data`, then `Driver.Unmarshal` turns it into
   `lab_data` rows (barcode, param, result, unit, completed date).
5. `Driver.PostUnmarshalActions` runs afterwards. For ASTM drivers this is where a device
   **query** (Q record) is answered; see [astm-query.md](astm-query.md).

## Drivers

`driver/driver.go` defines the `Driver` interface and the registry of names
(`hl7_231`, `astm`, `astm_E1394-97`, `text_*`). Each driver lives in its own package with
`<name>.go` (protocol handling) and `<name>_unmarshal.go` (parsing). Driver names are normalized
(lowercase, punctuation removed) before lookup.

## Query service

`services/deviceQuery.go` posts `{"barcode", "hardware_sn"}` to `DBSettings.QueryHost`
(`config/db_config.json`) and expects `{"indicators":[{"indicator": ...}]}`. A 2xx response
with a non-JSON content type is treated as an error; an empty indicator list is logged at
warn level with the request, status, content type and raw body so misconfiguration is visible.

## Storage

GORM with AutoMigrate on startup (`store/*.go`). Supported DB drivers: PostgreSQL, MySQL,
SQL Server (`DBSettings.DriverName`).

## Logging

`log/log.go` wraps zerolog; output and format come from `config/log_config.json`.
