package version

import (
	"fmt"

	"github.com/nbifrye/rmn/internal/api"
	"github.com/nbifrye/rmn/internal/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdCreate(f *cmdutil.Factory) *cobra.Command {
	var projectID, name, status, sharing, dueDate, description string

	cmd := &cobra.Command{
		Use:     "create",
		Aliases: []string{"new"},
		Short:   "Create a version",
		Long:    "Create a new Redmine project version.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if projectID == "" {
				return fmt.Errorf("--project is required")
			}
			if name == "" {
				return fmt.Errorf("--name is required")
			}

			client, err := f.APIClient()
			if err != nil {
				return err
			}

			params := api.VersionCreateParams{
				Name:        name,
				Status:      status,
				Sharing:     sharing,
				DueDate:     dueDate,
				Description: description,
			}

			version, err := client.CreateVersion(cmd.Context(), projectID, params)
			if err != nil {
				return err
			}

			// GetString cannot error for flags defined on the root command.
			output, _ := cmd.Root().PersistentFlags().GetString("output")
			if output == "json" {
				data, err := marshalJSON(version)
				if err != nil {
					return err
				}
				fmt.Fprintln(f.IO.Out, string(data))
				return nil
			}

			fmt.Fprintf(f.IO.Out, "Created version #%d: %s\n", version.ID, version.Name)
			return nil
		},
	}

	cmd.Flags().StringVarP(&projectID, "project", "p", "", "Project ID or identifier (required)")
	cmd.Flags().StringVar(&name, "name", "", "Version name (required)")
	cmd.Flags().StringVar(&status, "status", "", "Version status (open, locked, closed)")
	cmd.Flags().StringVar(&sharing, "sharing", "", "Version sharing (none, descendants, hierarchy, tree, system)")
	cmd.Flags().StringVar(&dueDate, "due-date", "", "Due date (YYYY-MM-DD)")
	cmd.Flags().StringVarP(&description, "description", "d", "", "Version description")

	return cmd
}
