package invites

import (
	"context"
	"fmt"

	"github.com/urfave/cli"

	"github.com/arigatomachine/cli/api"
	"github.com/arigatomachine/cli/config"
)

const approveInviteFailed = "Could not approve invitation to org, please try again."

// Approve executes logic required to approve an org invite
func Approve(ctx *cli.Context) error {
	usage := usageString(ctx)

	args := ctx.Args()
	if len(args) < 1 {
		text := "Missing email\n\n"
		text += usage
		return cli.NewExitError(text, -1)
	}
	email := ctx.Args()[0]

	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	client := api.NewClient(cfg)

	org, err := client.Orgs.GetByName(context.Background(), ctx.String("org"))
	if err != nil {
		return cli.NewExitError(approveInviteFailed, -1)
	}
	if org == nil {
		return cli.NewExitError("Org not found", -1)
	}

	states := []string{"accepted"}
	invites, err := client.Invites.List(context.Background(), org.ID, states)
	if err != nil {
		return cli.NewExitError("Failed to retrieve invites, please try again.", -1)
	}

	if len(invites) < 1 {
		return cli.NewExitError("Invite not found", -1)
	}

	var output api.ProgressFunc
	output = func(event *api.Event, err error) {
		if event != nil {
			fmt.Println(event.Message)
		}
	}

	invite := invites[0]
	err = client.Invites.Approve(context.Background(), *invite.ID, &output)
	if err != nil {
		return err
	}

	fmt.Println("")
	fmt.Println("You have approved " + email + "'s invitation.")
	fmt.Println("")
	fmt.Println("They are now a member of the organization!")

	return nil
}