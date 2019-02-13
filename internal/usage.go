// +build !windows

package internal

const Usage = `Opsagent, The plugin-driven server agent for collecting and reporting metrics.

Usage:

  opsagent [commands|flags]

The commands & flags are:

  config              print out full sample configuration to stdout
  version             print the version to stdout

  --aggregator-filter <filter>   filter the aggregators to enable, separator is :
  --config <file>                configuration file to load
  --config-directory <directory> directory containing additional *.conf files
  --debug                        turn on debug logging
  --input-filter <filter>        filter the inputs to enable, separator is :
  --input-list                   print available input plugins.
  --output-filter <filter>       filter the outputs to enable, separator is :
  --output-list                  print available output plugins.
  --pidfile <file>               file to write our pid to
  --pprof-addr <address>         pprof address to listen on, don't activate pprof if empty
  --processor-filter <filter>    filter the processors to enable, separator is :
  --quiet                        run in quiet mode
  --sample-config                print out full sample configuration
  --test                         gather metrics, print them out, and exit;
                                 processors, aggregators, and outputs are not run
  --usage <plugin>               print usage for a plugin, ie, 'opsagent --usage mysql'
  --version                      display the version and exit

Examples:

  # generate a opsagent config file:
  opsagent config > opsagent.conf

  # generate config with only cpu input & influxdb output plugins defined
  opsagent --input-filter cpu --output-filter influxdb config

  # run a single opsagent collection, outputing metrics to stdout
  opsagent --config opsagent.conf --test

  # run opsagent with all plugins defined in config file
  opsagent --config opsagent.conf

  # run opsagent, enabling the cpu & memory input, and influxdb output plugins
  opsagent --config opsagent.conf --input-filter cpu:mem --output-filter influxdb

  # run opsagent with pprof
  opsagent --config opsagent.conf --pprof-addr localhost:6060
`
