package main

type pipelineFile struct {
	Pipeline []struct {
		Name      string   `yaml:"name"`
		Exe       string   `yaml:"exe"`
		Msg       string   `yaml:"message"`
		Args      []string `yaml:"args"`
		Exception bool     `yaml:"exception"`
		Timeout   int      `yaml:"timeout"`
	} `yaml:"pipeline"`
}
