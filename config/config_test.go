package config

import (
	"errors"
	"testing"
	"github.com/prometheus/common/model"
)

func TestParseCfgFile(t *testing.T) {
	cases := []struct {
		name     string
		fileName string
		err      error
		mrc      []*RelabelConfig
	}{
		{
			name: "no config file provided",
			err:  errors.New("no config file provided"),
		},
		{
			name:     "non existing file",
			fileName: "i-dont-exist.yaml",
			err:      errors.New("valid file not found"),
		},
		{
			name:     "invalid yaml file",
			fileName: "testdata/invalid_yaml.yaml",
			err:      errors.New("yaml: line 1: mapping values are not allowed in this context"),
		},
		{
			name:     "empty file",
			fileName: "testdata/empty_file.yaml",
			err:      errors.New("valid file not found"),
		},
		{
			name:     "directory instead of file",
			fileName: "testdata",
			err:      errors.New("valid file not found"),
		},
		{
			name:     "file with no regex in metricrelabelconfig",
			fileName: "testdata/no_regex.yaml",
			err:      nil,
		},
		{
			name:     "file with no sourcelabels in metricrelabelconfig",
			fileName: "testdata/no_sourcelabels.yaml",
			err:      nil,
		},
		{
			name:     "invalid regex in metricrelabelconfig",
			fileName: "testdata/invalid_regex.yaml",
			err:      errors.New("error parsing regexp: missing closing ): `^(?:$^*()$`"),
		},
		{
			name:     "file with action 'addprefix' but no prefix in metricrelabelconfig",
			fileName: "testdata/no_prefix.yaml",
			err:      errors.New("addprefix action requires prefix"),
		},
		{
			name:     "file with action 'addprefix' and toplevel prefix",
			fileName: "testdata/with_mrc_and_prefix.yaml",
			err:      nil,
			mrc: []*RelabelConfig{
				{
					SourceLabels:model.LabelNames{model.MetricNameLabel},
					Regex:MustNewRegexp("my_too_large_metric"),
					Action:RelabelDrop,
				},
				{
					SourceLabels: model.LabelNames{model.MetricNameLabel},
					Regex: MustNewRegexp(".*"),
					Action: RelabelAddPrefix,
					Prefix: "my-prefix",
				},
			},
		},
		{
			name:     "file with only toplevel prefix",
			fileName: "testdata/with_only_prefix.yaml",
			err:      nil,
			mrc: []*RelabelConfig{
				{
					SourceLabels: model.LabelNames{model.MetricNameLabel},
					Regex: MustNewRegexp(".*"),
					Action: RelabelAddPrefix,
					Prefix: "my-prefix",
				},
			},
		},
	}

	for _, c := range cases {
		cfg, err := ParseCfgFile(c.fileName)
		if c.err != nil && c.err.Error() != err.Error() {
			t.Errorf("case '%s'. Expected %+v, Got %+v", c.name, c.err, err)
		}

		if c.err == nil && err != nil {
			t.Errorf("case '%s'. Expected no error, got: %+v", c.name, err)
		}

		if c.mrc != nil {
			if len(cfg.MetricRelabelConfigs) <= 0 {
				t.Errorf("case '%s'. Expected atleast one MetricRelabelConfig", c.name)
			}

			last := cfg.MetricRelabelConfigs[len(cfg.MetricRelabelConfigs)-1]
			if last.Prefix != "my-prefix" {
				t.Errorf("case '%s'. Expected prefix: %s, got %s", c.name, "my-prefix", last.Prefix)
			}
		}
	}
}
