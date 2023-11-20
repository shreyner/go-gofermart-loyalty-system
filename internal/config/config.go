package config

import (
	"net/url"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	Env                  string  `env:"ENV" envDefault:"dev"`
	Address              string  `env:"RUN_ADDRESS" envDefault:":8080"`
	DBURI                string  `env:"DATABASE_URI,required"`
	AccrualSystemAddress url.URL `env:"ACCRUAL_SYSTEM_ADDRESS,required"`

	IsProd bool
	IsDev  bool
}

func (c *Config) Parse() error {
	//address := flag.String("a", c.Address, "адрес и порт запуска сервиса")
	//DbURI := flag.String("d", c.DBURI, "адрес подключения к базе данных")
	//accrualSystemAddress := flag.String("r", c.AccrualSystemAddress.String(), "адрес системы расчёта начислений")

	//flag.Parse()

	//opt := env.Options{
	//Environment: map[string]string{
	//	"RUN_ADDRESS":            *address,
	//	"DATABASE_URI":           *DbURI,
	//	"ACCRUAL_SYSTEM_ADDRESS": *accrualSystemAddress,
	//},
	//}

	if err := env.Parse(c); err != nil {
		return err
	}

	c.IsDev = c.Env == "dev"
	c.IsProd = c.Env != "dev"

	return nil
}
