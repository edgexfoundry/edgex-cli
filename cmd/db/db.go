// Copyright © 2019 VMware, INC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// The purgedb command purges the entire Database. It performs the same action as the
// clean_mongo.js developer script. Unlike the clean_mongo.js, this command purges the
// database using API calls only. clean_mongo.js accesses the DB directly, which might
// always be possible using the CLI.
package db

import (
	purgedb "github.com/edgexfoundry/edgex-cli/cmd/db/purge"

	"github.com/spf13/cobra"
)

// NewCommand returns the db command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "db",
		Short: "Purges entire EdgeX Database. [USE WITH CAUTION]",
		Long:  `Purge DB`,
	}
	cmd.AddCommand(purgedb.NewCommand())
	return cmd
}
