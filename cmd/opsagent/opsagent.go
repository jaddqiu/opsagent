package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof" // Comment this line to disable pprof endpoint.
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"

	"github.com/jaddqiu/opsagent/agent"
	"github.com/jaddqiu/opsagent/internal"
	"github.com/jaddqiu/opsagent/internal/config"
	"github.com/jaddqiu/opsagent/logger"
	_ "github.com/jaddqiu/opsagent/plugins/aggregators/all"
	"github.com/jaddqiu/opsagent/plugins/tasks"
	_ "github.com/jaddqiu/opsagent/plugins/tasks/all"
	"github.com/jaddqiu/opsagent/plugins/inputs"
	_ "github.com/jaddqiu/opsagent/plugins/inputs/all"
	"github.com/jaddqiu/opsagent/plugins/outputs"
	_ "github.com/jaddqiu/opsagent/plugins/outputs/all"
	_ "github.com/jaddqiu/opsagent/plugins/processors/all"
	"github.com/kardianos/service"
)

var fDebug = flag.Bool("debug", false,
	"turn on debug logging")
var pprofAddr = flag.String("pprof-addr", "",
	"pprof address to listen on, not activate pprof if empty")
var fQuiet = flag.Bool("quiet", false,
	"run in quiet mode")
var fTest = flag.Bool("test", false, "gather metrics, print them out, and exit")
var fConfig = flag.String("config", "", "configuration file to load")
var fConfigDirectory = flag.String("config-directory", "",
	"directory containing additional *.conf files")
var fVersion = flag.Bool("version", false, "display the version and exit")
var fSampleConfig = flag.Bool("sample-config", false,
	"print out full sample configuration")
var fPidfile = flag.String("pidfile", "", "file to write our pid to")
var fTaskFilters = flag.String("task-filter", "",
	"filter the tasks to enable, separator is :")
var fTaskList = flag.Bool("task-list", false,
	"print available task plugins.")
var fInputFilters = flag.String("input-filter", "",
	"filter the inputs to enable, separator is :")
var fInputList = flag.Bool("input-list", false,
	"print available input plugins.")
var fOutputFilters = flag.String("output-filter", "",
	"filter the outputs to enable, separator is :")
var fOutputList = flag.Bool("output-list", false,
	"print available output plugins.")
var fAggregatorFilters = flag.String("aggregator-filter", "",
	"filter the aggregators to enable, separator is :")
var fProcessorFilters = flag.String("processor-filter", "",
	"filter the processors to enable, separator is :")
var fUsage = flag.String("usage", "",
	"print usage for a plugin, ie, 'opsagent --usage mysql'")
var fService = flag.String("service", "",
	"operate on the service (windows only)")
var fServiceName = flag.String("service-name", "opsagent", "service name (windows only)")
var fRunAsConsole = flag.Bool("console", false, "run as console application (windows only)")

var (
	version string
	commit  string
	branch  string
)

var stop chan struct{}

func reloadLoop(
	stop chan struct{},
	taskFilters []string,
	inputFilters []string,
	outputFilters []string,
	aggregatorFilters []string,
	processorFilters []string,
) {
	reload := make(chan bool, 1)
	reload <- true
	for <-reload {
		reload <- false

		ctx, cancel := context.WithCancel(context.Background())

		signals := make(chan os.Signal)
		signal.Notify(signals, os.Interrupt, syscall.SIGHUP,
			syscall.SIGTERM, syscall.SIGINT)
		go func() {
			select {
			case sig := <-signals:
				if sig == syscall.SIGHUP {
					log.Printf("I! Reloading Opsagent config")
					<-reload
					reload <- true
				}
				cancel()
			case <-stop:
				cancel()
			}
		}()

		err := runAgent(ctx, taskFilters, inputFilters, outputFilters)
		if err != nil {
			log.Fatalf("E! [opsagent] Error running agent: %v", err)
		}
	}
}

func runAgent(ctx context.Context,
	taskFilters []string,
	inputFilters []string,
	outputFilters []string,
) error {
	// Setup default logging. This may need to change after reading the config
	// file, but we can configure it to use our logger implementation now.
	logger.SetupLogging(false, false, "")
	log.Printf("I! Starting Opsagent %s", version)

	// If no other options are specified, load the config file and run.
	c := config.NewConfig()
	c.TaskFilters = taskFilters
	c.OutputFilters = outputFilters
	c.InputFilters = inputFilters
	err := c.LoadConfig(*fConfig)
	if err != nil {
		return err
	}

	if *fConfigDirectory != "" {
		err = c.LoadDirectory(*fConfigDirectory)
		if err != nil {
			return err
		}
	}
	if !*fTest && len(c.Outputs) == 0 {
		return errors.New("Error: no outputs found, did you provide a valid config file?")
	}
	if len(c.Tasks) == 0 {
		return errors.New("Error: no tasks found, did you provide a valid config file?")
	}

	if len(c.Inputs) == 0 {
		return errors.New("Error: no inputs found, did you provide a valid config file?")
	}

	if int64(c.Agent.Interval.Duration) <= 0 {
		return fmt.Errorf("Agent interval must be positive, found %s",
			c.Agent.Interval.Duration)
	}

	if int64(c.Agent.FlushInterval.Duration) <= 0 {
		return fmt.Errorf("Agent flush_interval must be positive; found %s",
			c.Agent.Interval.Duration)
	}

	ag, err := agent.NewAgent(c)
	if err != nil {
		return err
	}

	// Setup logging as configured.
	logger.SetupLogging(
		ag.Config.Agent.Debug || *fDebug,
		ag.Config.Agent.Quiet || *fQuiet,
		ag.Config.Agent.Logfile,
	)

	if *fTest {
		return ag.Test(ctx)
	}

	log.Printf("I! Loaded tasks: %s", strings.Join(c.TaskNames(), " "))
	log.Printf("I! Loaded inputs: %s", strings.Join(c.InputNames(), " "))
	log.Printf("I! Loaded aggregators: %s", strings.Join(c.AggregatorNames(), " "))
	log.Printf("I! Loaded processors: %s", strings.Join(c.ProcessorNames(), " "))
	log.Printf("I! Loaded outputs: %s", strings.Join(c.OutputNames(), " "))
	log.Printf("I! Tags enabled: %s", c.ListTags())

	if *fPidfile != "" {
		f, err := os.OpenFile(*fPidfile, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Printf("E! Unable to create pidfile: %s", err)
		} else {
			fmt.Fprintf(f, "%d\n", os.Getpid())

			f.Close()

			defer func() {
				err := os.Remove(*fPidfile)
				if err != nil {
					log.Printf("E! Unable to remove pidfile: %s", err)
				}
			}()
		}
	}

	return ag.Run(ctx)
}

func usageExit(rc int) {
	fmt.Println(internal.Usage)
	os.Exit(rc)
}

type program struct {
	taskFilters      []string
	inputFilters      []string
	outputFilters     []string
	aggregatorFilters []string
	processorFilters  []string
}

func (p *program) Start(s service.Service) error {
	go p.run()
	return nil
}
func (p *program) run() {
	stop = make(chan struct{})
	reloadLoop(
		stop,
		p.taskFilters,
		p.inputFilters,
		p.outputFilters,
		p.aggregatorFilters,
		p.processorFilters,
	)
}
func (p *program) Stop(s service.Service) error {
	close(stop)
	return nil
}

func formatFullVersion() string {
	var parts = []string{"Opsagent"}

	if version != "" {
		parts = append(parts, version)
	} else {
		parts = append(parts, "unknown")
	}

	if branch != "" || commit != "" {
		if branch == "" {
			branch = "unknown"
		}
		if commit == "" {
			commit = "unknown"
		}
		git := fmt.Sprintf("(git: %s %s)", branch, commit)
		parts = append(parts, git)
	}

	return strings.Join(parts, " ")
}

func main() {
	flag.Usage = func() { usageExit(0) }
	flag.Parse()
	args := flag.Args()

	taskFilters, inputFilters, outputFilters := []string{}, []string{}, []string{}
	if *fTaskFilters != "" {
		taskFilters = strings.Split(":"+strings.TrimSpace(*fTaskFilters)+":", ":")
	}
	if *fInputFilters != "" {
		inputFilters = strings.Split(":"+strings.TrimSpace(*fInputFilters)+":", ":")
	}
	if *fOutputFilters != "" {
		outputFilters = strings.Split(":"+strings.TrimSpace(*fOutputFilters)+":", ":")
	}

	aggregatorFilters, processorFilters := []string{}, []string{}
	if *fAggregatorFilters != "" {
		aggregatorFilters = strings.Split(":"+strings.TrimSpace(*fAggregatorFilters)+":", ":")
	}
	if *fProcessorFilters != "" {
		processorFilters = strings.Split(":"+strings.TrimSpace(*fProcessorFilters)+":", ":")
	}

	if *pprofAddr != "" {
		go func() {
			pprofHostPort := *pprofAddr
			parts := strings.Split(pprofHostPort, ":")
			if len(parts) == 2 && parts[0] == "" {
				pprofHostPort = fmt.Sprintf("localhost:%s", parts[1])
			}
			pprofHostPort = "http://" + pprofHostPort + "/debug/pprof"

			log.Printf("I! Starting pprof HTTP server at: %s", pprofHostPort)

			if err := http.ListenAndServe(*pprofAddr, nil); err != nil {
				log.Fatal("E! " + err.Error())
			}
		}()
	}

	if len(args) > 0 {
		switch args[0] {
		case "version":
			fmt.Println(formatFullVersion())
			return
		case "config":
			config.PrintSampleConfig(
				taskFilters,
				inputFilters,
				outputFilters,
				aggregatorFilters,
				processorFilters,
			)
			return
		}
	}

	// switch for flags which just do something and exit immediately
	switch {
	case *fOutputList:
		fmt.Println("Available Output Plugins:")
		for k := range outputs.Outputs {
			fmt.Printf("  %s\n", k)
		}
		return
	case *fTaskList:
		fmt.Println("Available Task Plugins:")
		for k := range tasks.Tasks{
			fmt.Printf("  %s\n", k)
		}
		return
	case *fInputList:
		fmt.Println("Available Input Plugins:")
		for k := range inputs.Inputs {
			fmt.Printf("  %s\n", k)
		}
		return
	case *fVersion:
		fmt.Println(formatFullVersion())
		return
	case *fSampleConfig:
		config.PrintSampleConfig(
			taskFilters,
			inputFilters,
			outputFilters,
			aggregatorFilters,
			processorFilters,
		)
		return
	case *fUsage != "":
		err0 := config.PrintTaskConfig(*fUsage)
		err := config.PrintInputConfig(*fUsage)
		err2 := config.PrintOutputConfig(*fUsage)
		if err0 != nil && err != nil && err2 != nil {
			log.Fatalf("E! %s, %s and %s", err0, err, err2)
		}
		return
	}

	shortVersion := version
	if shortVersion == "" {
		shortVersion = "unknown"
	}

	// Configure version
	if err := internal.SetVersion(shortVersion); err != nil {
		log.Println("Opsagent version already configured to: " + internal.Version())
	}

	if runtime.GOOS == "windows" && !(*fRunAsConsole) {
		svcConfig := &service.Config{
			Name:        *fServiceName,
			DisplayName: "Opsagent Data Collector Service",
			Description: "Collects data using a series of plugins and publishes it to" +
				"another series of plugins.",
			Arguments: []string{"--config", "C:\\Program Files\\Opsagent\\opsagent.conf"},
		}

		prg := &program{
			taskFilters:       taskFilters,
			inputFilters:      inputFilters,
			outputFilters:     outputFilters,
			aggregatorFilters: aggregatorFilters,
			processorFilters:  processorFilters,
		}
		s, err := service.New(prg, svcConfig)
		if err != nil {
			log.Fatal("E! " + err.Error())
		}
		// Handle the --service flag here to prevent any issues with tooling that
		// may not have an interactive session, e.g. installing from Ansible.
		if *fService != "" {
			if *fConfig != "" {
				(*svcConfig).Arguments = []string{"--config", *fConfig}
			}
			if *fConfigDirectory != "" {
				(*svcConfig).Arguments = append((*svcConfig).Arguments, "--config-directory", *fConfigDirectory)
			}
			err := service.Control(s, *fService)
			if err != nil {
				log.Fatal("E! " + err.Error())
			}
			os.Exit(0)
		} else {
			err = s.Run()
			if err != nil {
				log.Println("E! " + err.Error())
			}
		}
	} else {
		stop = make(chan struct{})
		reloadLoop(
			stop,
			taskFilters,
			inputFilters,
			outputFilters,
			aggregatorFilters,
			processorFilters,
		)
	}
}
