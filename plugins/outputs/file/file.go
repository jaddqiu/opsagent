package file

import (
	"fmt"
	"io"
	"os"

	"github.com/jaddqiu/opsagent"
	"github.com/jaddqiu/opsagent/plugins/outputs"
	"github.com/jaddqiu/opsagent/plugins/serializers"
)

type File struct {
	Files []string

	writers []io.Writer
	closers []io.Closer

	serializer serializers.Serializer
}

var sampleConfig = `
  ## Files to write to, "stdout" is a specially handled file.
  files = ["stdout", "/tmp/metrics.out"]

  ## Data format to output.
  ## Each data format has its own unique set of configuration options, read
  ## more about them here:
  ## https://github.com/jaddqiu/opsagent/blob/master/docs/DATA_FORMATS_OUTPUT.md
  data_format = "influx"
`

func (f *File) SetSerializer(serializer serializers.Serializer) {
	f.serializer = serializer
}

func (f *File) Connect() error {
	if len(f.Files) == 0 {
		f.Files = []string{"stdout"}
	}

	for _, file := range f.Files {
		if file == "stdout" {
			f.writers = append(f.writers, os.Stdout)
		} else {
			var of *os.File
			var err error
			if _, err := os.Stat(file); os.IsNotExist(err) {
				of, err = os.Create(file)
			} else {
				of, err = os.OpenFile(file, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
			}

			if err != nil {
				return err
			}
			f.writers = append(f.writers, of)
			f.closers = append(f.closers, of)
		}
	}
	return nil
}

func (f *File) Close() error {
	var errS string
	for _, c := range f.closers {
		if err := c.Close(); err != nil {
			errS += err.Error() + "\n"
		}
	}
	if errS != "" {
		return fmt.Errorf(errS)
	}
	return nil
}

func (f *File) SampleConfig() string {
	return sampleConfig
}

func (f *File) Description() string {
	return "Send opsagent metrics to file(s)"
}

func (f *File) Write(metrics []opsagent.Metric) error {
	var writeErr error = nil
	for _, metric := range metrics {
		b, err := f.serializer.Serialize(metric)
		if err != nil {
			return fmt.Errorf("failed to serialize message: %s", err)
		}

		for _, writer := range f.writers {
			_, err = writer.Write(b)
			if err != nil && writer != os.Stdout {
				writeErr = fmt.Errorf("E! failed to write message: %s, %s", b, err)
			}
		}
	}
	return writeErr
}

func init() {
	outputs.Add("file", func() opsagent.Output {
		return &File{}
	})
}
