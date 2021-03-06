package conf

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

type mapLoader struct {
	layers M
}

func TestLoad(t *testing.T) {
	configProc := TNewProcessor()

	tConfig, err := configProc.Load(
		M{
			"paramA": "default:valA",
			"paramZ": "default:valZ",
		},

		"test:foo",
		"test:bar",
	)

	if err != nil {
		t.Error(err)
		return
	}

	eConfig := M{
		"paramA": "foo:valA",
		"paramB": "bar:valB",
		"paramC": "bar:valC",

		"paramD": M{
			"paramDA": "foo:valDA",
			"paramDB": "bar:valDB",
			"paramDC": "bar:valDC",
			"paramDE": "foo:bar:valDC",

			"paramDF": S{
				"foo:valDFA",
				"foo:valDFB",
				"foo:foo:valDA",
			},
		},

		"paramE": S{
			"bar:valEA",
			"bar:valEB",
		},

		"paramF": "foo:bar:valB",
		"paramG": "bar:foo:valDA",
		"paramH": "foo:bar:valEA",
		"paramI": "bar:foo:bar:valEA",
		"paramJ": "foo:bar:foo:bar:valEA",
		"paramK": "bar:foo:valDFB:foo:bar:valDC",
		"paramL": "foo:${paramD.paramDE}:${}:${paramD.paramDA}",

		"paramM": M{
			"paramDA": "foo:valDA",
			"paramDB": "bar:valDB",
			"paramDC": "bar:valDC",
			"paramDE": "foo:bar:valDC",

			"paramDF": S{
				"foo:valDFA",
				"foo:valDFB",
				"foo:foo:valDA",
			},
		},

		"paramN": M{
			"paramNA": "foo:valNA",
			"paramNB": "foo:valNB",

			"paramNC": M{
				"paramNCA": "foo:valNCA",
				"paramNCB": "bar:valNCB",
				"paramNCC": "bar:valNCC",
				"paramNCD": "bar:foo:valNCA",
				"paramNCE": "foo:valNB",
			},
		},

		"paramO": M{
			"paramOA": "moo:valOA",
			"paramOB": "jar:valOB",
			"paramOC": "jar:valOC",

			"paramOD": M{
				"paramODA": "moo:valODA",
				"paramODB": "jar:valODB",
				"paramODC": "jar:valODC",
				"paramODD": "jar:bar:valNCB",
			},

			"paramOE": S{
				"zoo:valA",
				"zoo:valB",
			},
		},

		"paramP": M{
			"paramODA": "moo:valODA",
			"paramODB": "jar:valODB",
			"paramODC": "jar:valODC",
			"paramODD": "jar:bar:valNCB",
		},

		"paramS": "bar:valS",
		"paramT": "bar:valY",
		"paramY": "bar:valY",
		"paramZ": "default:valZ",
	}

	if !reflect.DeepEqual(tConfig, eConfig) {
		t.Errorf("unexpected configuration returned: %#v", tConfig)
	}
}

func TestDisableProcessing(t *testing.T) {
	configProc := NewProcessor(
		ProcessorConfig{
			DisableProcessing: true,
		},
	)

	tConfig, err := configProc.Load(
		M{
			"paramA": "coo:valA",
			"paramB": "coo:${paramA}",
		},
	)

	if err != nil {
		t.Error(err)
		return
	}

	eConfig := M{
		"paramA": "coo:valA",
		"paramB": "coo:${paramA}",
	}

	if !reflect.DeepEqual(tConfig, eConfig) {
		t.Errorf("unexpected configuration returned: %#v", tConfig)
	}
}

func TestDecode(t *testing.T) {
	type testConfig struct {
		ParamA string `conf:"test_paramA"`
		ParamB int    `conf:"test_paramB"`
		ParamC []string
		ParamD map[string]bool
	}

	configRaw := M{
		"test_paramA": "foo:val",
		"test_paramB": 1234,
		"paramC":      []string{"moo:val1", "moo:val2"},
		"paramD": map[string]bool{
			"zoo": true,
			"arr": false,
		},
	}

	var tConfig testConfig
	Decode(configRaw, &tConfig)

	eConfig := testConfig{
		ParamA: "foo:val",
		ParamB: 1234,
		ParamC: []string{"moo:val1", "moo:val2"},
		ParamD: map[string]bool{
			"zoo": true,
			"arr": false,
		},
	}

	if !reflect.DeepEqual(tConfig, eConfig) {
		t.Errorf("unexpected configuration returned: %#v", tConfig)
	}
}

func TestPanic(t *testing.T) {
	t.Run("no_locators",
		func(t *testing.T) {
			defer func() {
				err := recover()
				errStr := fmt.Sprintf("%v", err)

				if err == nil {
					t.Error("no error happened")
				} else if strings.Index(errStr, "no configuration locators") == -1 {
					t.Error("other error happened:", err)
				}
			}()

			configProc := TNewProcessor()
			configProc.Load()
		},
	)
}

func TestErrors(t *testing.T) {
	configProc := TNewProcessor()

	t.Run("empty_locator",
		func(t *testing.T) {
			_, err := configProc.Load("")

			if err == nil {
				t.Error("no error happened")
			} else if strings.Index(err.Error(), "empty configuration locator") == -1 {
				t.Error("other error happened:", err)
			}
		},
	)

	t.Run("invalid_locator",
		func(t *testing.T) {
			_, err := configProc.Load(42)

			if err == nil {
				t.Error("no error happened")
			} else if strings.Index(err.Error(), "configuration locator must be of type") == -1 {
				t.Error("other error happened:", err)
			}
		},
	)

	t.Run("missing_loader",
		func(t *testing.T) {
			_, err := configProc.Load("foo")

			if err == nil {
				t.Error("no error happened")
			} else if strings.Index(err.Error(), "missing loader name") == -1 {
				t.Error("other error happened:", err)
			}
		},
	)

	t.Run("loader_not_found",
		func(t *testing.T) {
			_, err := configProc.Load("etcd:foo")

			if err == nil {
				t.Error("no error happened")
			} else if strings.Index(err.Error(), "loader not found") == -1 {
				t.Error("other error happened:", err)
			}
		},
	)

	t.Run("invalid_config_type",
		func(t *testing.T) {
			_, err := configProc.Load("test:zoo")

			if err == nil {
				t.Error("no error happened")
			} else if strings.Index(err.Error(), "has invalid type") == -1 {
				t.Error("other error happened:", err)
			}
		},
	)

	t.Run("invalid_ref",
		func(t *testing.T) {
			_, err := configProc.Load("test:invalid_ref")

			if err == nil {
				t.Error("no error happened")
			} else if strings.Index(err.Error(), "invalid _ref directive") == -1 {
				t.Error("other error happened:", err)
			}
		},
	)

	t.Run("invalid_ref_name",
		func(t *testing.T) {
			_, err := configProc.Load("test:invalid_ref_name")

			if err == nil {
				t.Error("no error happened")
			} else if strings.Index(err.Error(), "reference name must be of type") == -1 {
				t.Error("other error happened:", err)
			}
		},
	)

	t.Run("invalid_ref_first_defined",
		func(t *testing.T) {
			_, err := configProc.Load("test:invalid_ref_first_defined")

			if err == nil {
				t.Error("no error happened")
			} else if strings.Index(err.Error(), "firstDefined list must be of type") == -1 {
				t.Error("other error happened:", err)
			}
		},
	)

	t.Run("invalid_ref_first_defined_argument",
		func(t *testing.T) {
			_, err := configProc.Load("test:invalid_ref_first_defined_argument")

			if err == nil {
				t.Error("no error happened")
			} else if strings.Index(err.Error(), "reference name in firstDefined") == -1 {
				t.Error("other error happened:", err)
			}
		},
	)

	t.Run("invalid_include",
		func(t *testing.T) {
			_, err := configProc.Load("test:invalid_include")

			if err == nil {
				t.Error("no error happened")
			} else if strings.Index(err.Error(), "invalid _include directive") == -1 {
				t.Error("other error happened:", err)
			}
		},
	)

	t.Run("invalid_index",
		func(t *testing.T) {
			_, err := configProc.Load("test:invalid_index")

			if err == nil {
				t.Error("no error happened")
			} else if strings.Index(err.Error(), "invalid slice index") == -1 {
				t.Error("other error happened:", err)
			}
		},
	)

	t.Run("index_out_of_range",
		func(t *testing.T) {
			_, err := configProc.Load("test:index_out_of_range")

			if err == nil {
				t.Error("no error happened")
			} else if strings.Index(err.Error(), "index out of range") == -1 {
				t.Error("other error happened:", err)
			}
		},
	)
}

func TNewProcessor() *Processor {
	mapLdr := TNewLoader()

	configProc := NewProcessor(
		ProcessorConfig{
			Loaders: map[string]Loader{
				"test": mapLdr,
			},
		},
	)

	return configProc
}

func TNewLoader() Loader {
	return &mapLoader{
		M{
			"foo": M{
				"paramA": "foo:valA",
				"paramB": "foo:valB",

				"paramD": M{
					"paramDA": "foo:valDA",
					"paramDB": "foo:valDB",
					"paramDE": "foo:${.paramDC}",

					"paramDF": S{
						"foo:valDFA",
						"foo:valDFB",
						"foo:${..paramDA}",
					},
				},

				"paramE": S{
					"foo:valEA",
					"foo:valEB",
				},

				"paramF": "foo:${paramB}",
				"paramH": "foo:${paramE.0}",
				"paramJ": "foo:${paramI}",
				"paramL": "foo:$${paramD.paramDE}:${}:$${paramD.paramDA}",

				"paramN": M{
					"paramNA": "foo:valNA",
					"paramNB": "foo:valNB",

					"paramNC": M{
						"paramNCA": "foo:valNCA",
						"paramNCB": "foo:valNCB",
						"paramNCE": M{"_ref": "..paramNB"},
					},
				},

				"paramO": M{
					"_include": S{"test:moo", "test:jar"},
				},
			},

			"bar": M{
				"paramB": "bar:valB",
				"paramC": "bar:valC",

				"paramD": M{
					"paramDB": "bar:valDB",
					"paramDC": "bar:valDC",
				},

				"paramE": S{
					"bar:valEA",
					"bar:valEB",
				},

				"paramG": "bar:${paramD.paramDA}",
				"paramI": "bar:${paramH}",
				"paramK": "bar:${paramD.paramDF.1}:${paramD.paramDE}",
				"paramM": M{"_ref": "paramD"},

				"paramN": M{
					"paramNC": M{
						"paramNCB": "bar:valNCB",
						"paramNCC": "bar:valNCC",
						"paramNCD": "bar:${paramN.paramNC.paramNCA}",
					},
				},

				"paramP": M{"_ref": "paramO.paramOD"},

				"paramS": M{
					"_ref": M{
						"name":    "paramX",
						"default": "bar:valS",
					},
				},

				"paramT": M{
					"_ref": M{
						"firstDefined": S{"paramX", "paramY"},
						"default":      "bar:valT",
					},
				},

				"paramY": "bar:valY",
			},

			"moo": M{
				"paramOA": "moo:valOA",
				"paramOB": "moo:valOB",

				"paramOD": M{
					"paramODA": "moo:valODA",
					"paramODB": "moo:valODB",
				},
			},

			"jar": M{
				"paramOB": "jar:valOB",
				"paramOC": "jar:valOC",

				"paramOD": M{
					"paramODB": "jar:valODB",
					"paramODC": "jar:valODC",
					"paramODD": "jar:${paramN.paramNC.paramNCB}",
				},

				"paramOE": M{
					"_include": S{"test:zoo"},
				},
			},

			"zoo": S{
				"zoo:valA",
				"zoo:valB",
			},

			"invalid_ref": M{
				"paramQ": M{"_ref": 42},
			},

			"invalid_ref_name": M{
				"_ref": M{
					"name":    42,
					"default": "foo",
				},
			},

			"invalid_ref_first_defined": M{
				"_ref": M{
					"firstDefined": 42,
					"default":      "bar:valT",
				},
			},

			"invalid_ref_first_defined_argument": M{
				"_ref": M{
					"firstDefined": S{42},
					"default":      "bar:valT",
				},
			},

			"invalid_include": M{
				"paramQ": M{"_include": 42},
			},

			"invalid_index": M{
				"paramQ": S{"valA", "valB"},
				"paramR": M{"_ref": "paramQ.paramQA"},
			},

			"index_out_of_range": M{
				"paramQ": S{"valA", "valB"},
				"paramR": M{"_ref": "paramQ.2"},
			},
		},
	}
}

func (p *mapLoader) Load(loc *Locator) (interface{}, error) {
	key := loc.BareLocator
	layer, _ := p.layers[key]

	return layer, nil
}

func ExampleDecode() {
	type DBConfig struct {
		Host     string `conf:"server_host"`
		Port     int    `conf:"server_port"`
		DBName   string
		Username string
		Password string
	}

	configRaw := M{
		"server_host": "stat.mydb.com",
		"server_port": 1234,
		"dbname":      "stat",
		"username":    "stat_writer",
		"password":    "some_pass",
	}

	var config DBConfig
	Decode(configRaw, &config)

	fmt.Printf("%v", config)

	// Output: {stat.mydb.com 1234 stat stat_writer some_pass}
}