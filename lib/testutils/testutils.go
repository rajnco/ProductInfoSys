package testutils

import (
	"os"
	"github.com/sirupen/logrus"
)

type envVariables map[string]string

var (
	EnvVariableValues = envVariables{
		"APP_NAME": "ProductInfoSys",
		"APP_ENV": "develop",
		"DBNAME": "products.db",
	}
	TestEnvVariableValues = envVariables{
		"APP_NAME": "ProductInfoSys",
		"APP_ENV": "test",
		"DBNAME": "products_test.db",
	}
)

func LoadEnv(opts ...envVariables){
	loadEnvVariables(EnvVariableValues, opts...)
}

func LoadTestEnv(opts ...envVariables){
	loadEnvVariables(TestEnvVariableValues, opts...)
}

func loadEnvVariables(preset envVariables, opts ...envVariables){
	for k, v := range preset {
		err := os.Setenv(k, v)
		if err != nil {
			logrus.Errorf("error setting environment variable %v", err)
		}
	}

	for _, env := range opts{
		for k, v := range env {
			err := os.Setenv(k,v) 
			if err != nil {
				logrus.Errorf("error setting environment variable %v", err)
			}
		}
	}
}
