# JC lab analyzers

How each analyzer at the JC lab talks to the device server, and which driver to set on its
device model. Captures of real sessions were used to build and test the drivers.

| Analyzer | Driver (device model) | Protocol |
|---|---|---|
| HumaCount 5D (Dymind DF5x) | `hl7_231` | HL7 2.3.1 `ORU^R01` over MLLP |
| Wondfo Accre 8 | `astm_wondfo_accre_8` | ASTM-like upload, control characters sent as text |

Every analyzer must also exist as a device whose `net_address` is the analyzer's IP; otherwise
the device server logs `failed to get a device by network address` and drops the message.

## Wondfo Accre 8 — `astm_wondfo_accre_8`

Package `driver/driver_astm_wondfo_accre_8`. Reference: Tisenc "ACCRE" LIS interface manual
(ASTM E1394-97 chapter). The analyzer only uploads results; it never queries the host.

### Wire format

The analyzer connects over TCP, sends the text `Network is connected!`, and later sends one
block per sample:

```
<ENQ>\x021H|\x5e&||| ACCRE8 |||||Host||1|20090119131415<CR>\x0353<CR><LF>
\x022P|1|||||||||||||^<CR>\x03E9<CR><LF>
\x023O|1|128123^nadiraze^1^2|||R|20261103 17:30:48|||||||||Serum||||||||||O<CR>\x0355<CR><LF>
\x024R|1|^^^D-Dimer|46.33|μg/mL(FEU)|^\x5e|N||F||||20261103 17:30:48<CR>\x0332<CR><LF>
\x026L|1|N<CR>\x030A<CR><LF><EOT>
```

Deviations from ASTM that the generic `astm` drivers cannot handle:

- `ENQ` and `EOT` are real control bytes, but `STX`, `ETX` and `^` inside the block arrive as the
  **literal text** `\x02`, `\x03` and `\x5e`. The driver converts them while parsing (real control
  bytes are accepted too) and stores the raw data exactly as received.
- Checksums, the header timestamp and the first component of O-3 (`128123`) are constants copied
  from the vendor's example message. Checksums are therefore not verified.
- The whole block, often several blocks, is sent in one go without waiting for per-frame ACKs.
  The driver answers every received chunk with a single `ACK`.
- A block may contain an order without any result record; it produces no lab data and no error.

### Field mapping

| `lab_data` | Source |
|---|---|
| `barcode` | O-4 (scanned specimen barcode) when set, otherwise the 2nd component of O-3 (sample ID typed on the analyzer) |
| `param` | 4th component of R-3 (`^^^D-Dimer`) |
| `result` | R-4 |
| `unit` | R-5 (UTF-8, may be empty) |
| `completed_date` | R-13, `YYYYMMDD HH:MM:SS` (the manual's `YYYYMMDDHHMMSS` is accepted too) |
| `index` | position of the result within its order, starting at 1 |

Result records without a test name or value are skipped. A result with an unparseable date, or
without a preceding order, fails the message (`raw_data.processed = false`).

### Analyzer-side setup

- LIS protocol: ASTM, network (TCP) mode, host IP and port of the device server.
- The operator must enter (or scan) the LIS barcode as the **sample ID**; that value becomes
  `barcode`.
