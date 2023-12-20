// SPDX-FileCopyrightText: 2023 INDITEX S.A
//
// SPDX-License-Identifier: Apache-2.0

package logging

import (
	"fmt"
	"os"
)

const envDebug = "SHERPA_DEBUG"

func Debug(message ...string) {
	if os.Getenv(envDebug) == "" {
		return
	}

	fmt.Printf("DEBUG: %s\n", message)
}

func Debugf(message string, args ...interface{}) {
	if os.Getenv(envDebug) == "" {
		return
	}

	fmt.Printf("DEBUG: %s\n", fmt.Sprintf(message, args...))
}

func Info(message string) {
	fmt.Println(PaintInfo(message))
}

func Error(message string) {
	fmt.Println(PaintError("ERROR: " + message))
}

func Errorf(message string, args ...interface{}) {
	fmt.Println(PaintError("ERROR: " + fmt.Sprintf(message, args...)))
}

func PrintCommandHeader(command string) {
	fmt.Printf("\n=> Running %s command in %s.\n\n",
		PaintInfo(command),
		PaintInfo("Sherpa"),
	)
}
