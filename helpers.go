package main

import "os"

func getEnvOrDefault(llave string, defecto string) string {
	valor := os.Getenv(llave)
	if valor == "" {
		return defecto
	}
	return valor
}
