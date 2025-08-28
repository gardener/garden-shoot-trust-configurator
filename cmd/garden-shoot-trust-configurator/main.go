// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"os"

	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/gardener/garden-shoot-trust-configurator/cmd/garden-shoot-trust-configurator/app"
)

func main() {
	cmd := app.NewCommand()

	if err := cmd.ExecuteContext(ctrl.SetupSignalHandler()); err != nil {
		ctrl.Log.Error(err, "Failed to run application", "name", cmd.Name())
		os.Exit(1)
	}
}
