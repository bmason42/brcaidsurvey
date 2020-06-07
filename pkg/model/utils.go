/*
 * Copyright (c) 2020.  This software is made for the Black Rock City Aid group and is provided AS IS with no support or liability under the Apache 2 license.
 */

package model

import "os"

func GetEnvVar(key string, defaultValue string) string {
	v := os.Getenv(key)
	if len(v) == 0 {
		v = defaultValue
	}
	return v
}
