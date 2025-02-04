// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package projectconfig

import (
	"context"
	"fmt"

	"github.com/charmbracelet/huh"
	apiclient_util "github.com/daytonaio/daytona/internal/util/apiclient"
	"github.com/daytonaio/daytona/pkg/apiclient"
	"github.com/daytonaio/daytona/pkg/views"
	views_util "github.com/daytonaio/daytona/pkg/views/util"
	"github.com/daytonaio/daytona/pkg/views/workspace/selection"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var allFlag bool
var yesFlag bool
var forceFlag bool

var projectConfigDeleteCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"remove", "rm"},
	Short:   "Delete a project config",
	Args:    cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var selectedProjectConfig *apiclient.ProjectConfig
		var selectedProjectConfigName string

		apiClient, err := apiclient_util.GetApiClient(nil)
		if err != nil {
			return err
		}

		if allFlag {
			if !yesFlag {
				form := huh.NewForm(
					huh.NewGroup(
						huh.NewConfirm().
							Title("Delete all project configs?").
							Description("Are you sure you want to delete all project configs?").
							Value(&yesFlag),
					),
				).WithTheme(views.GetCustomTheme())

				err := form.Run()
				if err != nil {
					return err
				}

				if !yesFlag {
					fmt.Println("Operation canceled.")
					return nil
				}
			}

			projectConfigs, res, err := apiClient.ProjectConfigAPI.ListProjectConfigs(context.Background()).Execute()
			if err != nil {
				return apiclient_util.HandleErrorResponse(res, err)
			}

			if len(projectConfigs) == 0 {
				views_util.NotifyEmptyProjectConfigList(false)
				return nil
			}

			for _, projectConfig := range projectConfigs {
				selectedProjectConfigName = projectConfig.Name
				res, err := apiClient.ProjectConfigAPI.DeleteProjectConfig(context.Background(), selectedProjectConfigName).Execute()
				if err != nil {
					log.Error(apiclient_util.HandleErrorResponse(res, err))
					continue
				}
				views.RenderInfoMessage("Deleted project config: " + selectedProjectConfigName)
			}
			return nil
		}

		if len(args) == 0 {
			projectConfigs, res, err := apiClient.ProjectConfigAPI.ListProjectConfigs(context.Background()).Execute()
			if err != nil {
				return apiclient_util.HandleErrorResponse(res, err)
			}

			if len(projectConfigs) == 0 {
				views.RenderInfoMessage("No project configs found")
				return nil
			}

			selectedProjectConfig = selection.GetProjectConfigFromPrompt(projectConfigs, 0, false, false, "Delete")
			if selectedProjectConfig == nil {
				return nil
			}
			selectedProjectConfigName = selectedProjectConfig.Name
		} else {
			selectedProjectConfigName = args[0]
		}

		res, err := apiClient.ProjectConfigAPI.DeleteProjectConfig(context.Background(), selectedProjectConfigName).Force(forceFlag).Execute()
		if err != nil {
			return apiclient_util.HandleErrorResponse(res, err)
		}

		views.RenderInfoMessage("Project config deleted successfully")
		return nil
	},
}

func init() {
	projectConfigDeleteCmd.Flags().BoolVarP(&allFlag, "all", "a", false, "Delete all project configs")
	projectConfigDeleteCmd.Flags().BoolVarP(&yesFlag, "yes", "y", false, "Confirm deletion without prompt")
	projectConfigDeleteCmd.Flags().BoolVarP(&forceFlag, "force", "f", false, "Force delete prebuild")
}
