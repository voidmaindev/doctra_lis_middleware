# ASTM E1394-97 query handshake

Applies to driver `astm_E1394-97` (`driver/driver_astm_e1394_97`), used with the Snibe Maglumi.

## Device → middleware

The analyzer sends a query for the orders of one sample:

```
ENQ
<STX>1H|\^&||PSWD|MaglumiBio|||||Lis||P|E1394-97|20261021<CR><ETX>92<CR><LF>
<STX>2Q|1|^998877||ALL||||||||O<CR><ETX>6A<CR><LF>
<STX>3L|1|N<CR><ETX>06<CR><LF>
EOT
```

The middleware ACKs every frame. After EOT, `Unmarshal` detects the Q record and returns
`additionalData` with the parsed query (`QRY`) and the parsed frames (`MSGS`) instead of lab
data. `processDeviceMessage` then calls `PostUnmarshalActions`, which runs `doQuery`.

## Middleware → device

`buildQueryReply` asks the query service for the sample's orders and formats the reply with
`generateASTMMessagesFromQuery`; `sendASTMFrames` transmits it:

```
ENQ                      -> wait for ACK
<STX>1H|\^&||PSWD|MaglumiBio|||||Lis||P|E1394-97|<date><CR><ETX><cs><CR><LF>   -> wait for ACK
<STX>2P|1<CR><ETX><cs><CR><LF>                                                  -> wait for ACK
<STX>3O|1|998877||^^^ALP\^^^TSH|P|R<CR><ETX><cs><CR><LF>                        -> wait for ACK   (only if orders exist)
<STX>4L|1|N<CR><ETX><cs><CR><LF>                                                -> wait for ACK
EOT
```

The checksum `<cs>` is the modulo-256 sum of frame number, record, CR and ETX, as two
uppercase hex digits (`formatASTMFrame`).

### No orders for the barcode

When the query host returns an empty indicator list the middleware still answers with
`H`, `P|1`, `L|1|N` (no `O` record). Before this change nothing was sent at all and the
Maglumi waited about 12 seconds, then closed the connection (`_samples/snibe240p3.txt`).

### Errors

- Query host unreachable / HTTP error / non-JSON 200 response: `doQuery` returns an error and
  logs `failed to query the service`. No reply is sent to the device.
- Empty indicator list: a warn log line `query host returned no indicators: ...` includes the
  request body, status, content type and the raw response body (truncated), then the empty
  reply is sent.

## Result upload (for contrast)

A result upload (H, P, O, R, L records) only requires an ACK per frame. The middleware stores
the parsed results in `lab_data` and sends nothing else back.
