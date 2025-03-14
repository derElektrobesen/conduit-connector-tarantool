# Conduit Connector for <!-- readmegen:name -->Tarantool<!-- /readmegen:name -->

[Conduit](https://conduit.io) connector for <!-- readmegen:name -->Tarantool<!-- /readmegen:name -->.

<!-- readmegen:description -->
A destination connector for [Tarantool](https://www.tarantool.io).

Source connector isn't implemented yet.<!-- /readmegen:description -->

## Source

A source connector pulls data from an external resource and pushes it to
downstream resources via Conduit.

### Configuration

<!-- readmegen:source.parameters.yaml -->
```yaml
version: 2.2
pipelines:
  - id: example
    status: running
    connectors:
      - id: example
        plugin: "tarantool"
        settings:
```
<!-- /readmegen:source.parameters.yaml -->

## Destination

A destination connector pushes data from upstream resources to an external
resource via Conduit.

### Configuration

<!-- readmegen:destination.parameters.yaml -->
```yaml
version: 2.2
pipelines:
  - id: example
    status: running
    connectors:
      - id: example
        plugin: "tarantool"
        settings:
          # password specifies a password to access tarantool
          # Type: string
          # Required: yes
          password: ""
          # replicasets specifies tarantool cluster configuration
          # Example configuration:
          #   replicasets: |-
          #     - name: "replicaset_1"                             # optional
          #       uuid: "ff69c808-039f-4478-9f28-27a487b3d1d3"     # optional
          #       replicas:                                        # *required*
          #         - addr: "127.0.0.1:1001"                       # *required*
          #           name: "1_1"                                  # optional
          #           uuid: "15299fc8-fb53-44dc-9a84-2332fad9687c" # optional
          #         - addr: "127.0.0.1:1002"
          #           name: "1_2"
          #           uuid: "15299fc8-fb53-44dc-9a84-2332fad9687c"
          #     - name: "replicaset_2"
          #       uuid: "29bbdcc4-2aa8-475f-b660-4fdc20dd5052"
          #       replicas:
          #         - addr: "127.0.0.1:1003"
          #           name: "2_1"
          #         - addr: "127.0.0.1:1004"
          #           name: "2_2"
          # Type: string
          # Required: yes
          replicasets: ""
          # total_buckets specifies a number of buckets in tarantool cluster
          # Type: int
          # Required: yes
          total_buckets: "0"
          # user specifies username to access tarantool
          # Type: string
          # Required: yes
          user: ""
          # shard_function specifies a shard function to be used during
          # sharding. At now only default sharding function is supported.
          # Type: string
          # Required: no
          shard_function: "default"
          # Maximum delay before an incomplete batch is written to the
          # destination.
          # Type: duration
          # Required: no
          sdk.batch.delay: "0"
          # Maximum size of batch before it gets written to the destination.
          # Type: int
          # Required: no
          sdk.batch.size: "0"
          # Allow bursts of at most X records (0 or less means that bursts are
          # not limited). Only takes effect if a rate limit per second is set.
          # Note that if `sdk.batch.size` is bigger than `sdk.rate.burst`, the
          # effective batch size will be equal to `sdk.rate.burst`.
          # Type: int
          # Required: no
          sdk.rate.burst: "0"
          # Maximum number of records written per second (0 means no rate
          # limit).
          # Type: float
          # Required: no
          sdk.rate.perSecond: "0"
          # The format of the output record. See the Conduit documentation for a
          # full list of supported formats
          # (https://conduit.io/docs/using/connectors/configuration-parameters/output-format).
          # Type: string
          # Required: no
          sdk.record.format: "opencdc/json"
          # Options to configure the chosen output record format. Options are
          # normally key=value pairs separated with comma (e.g.
          # opt1=val2,opt2=val2), except for the `template` record format, where
          # options are a Go template.
          # Type: string
          # Required: no
          sdk.record.format.options: ""
          # Whether to extract and decode the record key with a schema.
          # Type: bool
          # Required: no
          sdk.schema.extract.key.enabled: "true"
          # Whether to extract and decode the record payload with a schema.
          # Type: bool
          # Required: no
          sdk.schema.extract.payload.enabled: "true"
```
<!-- /readmegen:destination.parameters.yaml -->

## Development

- To install the required tools, run `make install-tools`.
- To generate code (mocks, re-generate `connector.yaml`, update the README,
  etc.), run `make generate`.

## How to build?

Run `make build` to build the connector.

## Testing

Run `make test` to run all the unit tests. Run `make test-integration` to run
the integration tests.

The Docker compose file at `test/docker-compose.yml` can be used to run the
required resource locally.

## How to release?

The release is done in two steps:

- Bump the version in [connector.yaml](/connector.yaml). This can be done
  with [bump_version.sh](/scripts/bump_version.sh) script, e.g.
  `scripts/bump_version.sh 2.3.4` (`2.3.4` is the new version and needs to be a
  valid semantic version). This will also automatically create a PR for the
  change.
- Tag the connector, which will kick off a release. This can be done
  with [tag.sh](/scripts/tag.sh).

## Known Issues & Limitations

- Known issue A
- Limitation A

## Planned work

- [ ] Item A
- [ ] Item B
