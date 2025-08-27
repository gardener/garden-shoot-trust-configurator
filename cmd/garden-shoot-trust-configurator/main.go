// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"os"

	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/gardener/garden-shoot-trust-configurator/cmd/garden-shoot-trust-configurator/app"
)

func main() {
	fmt.Println("Starting garden-shoot-trust-configurator...")
	cmd := app.NewCommand()

	if err := cmd.ExecuteContext(ctrl.SetupSignalHandler()); err != nil {
		ctrl.Log.Error(err, "Failed to run application", "name", cmd.Name())
		os.Exit(1)
	}
}
