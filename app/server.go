package app

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/crowdflux/angel/app/api"
	"github.com/crowdflux/angel/app/services/flu_svc"
)

func Start() {

	/**
	Go to https://github.com/codegangsta/cli for more help regard 'help'
	*/
	app := cli.NewApp()

	app.Name = "Goplay"
	app.Usage = "Play with me"
	app.Author = "himanshu144141"
	app.Copyright = "Playment Inc."
	app.Email = "support@playment.in"
	app.Version = "0.0.1"
	app.Action = func(c *cli.Context) {

		println("Angel server started!")
		println("Run angel -h for help")
		println(
			`		       .__                                      __
		______ |  | _____  ___.__. _____   ____   _____/  |_
		\____ \|  | \__  \<   |  |/     \_/ __ \ /    \   __\
		|  |_> >  |__/ __ \\___  |  Y Y  \  ___/|   |  \  |
		|   __/|____(____  / ____|__|_|  /\___  >___|  /__|
		|__|             \/\/          \/     \/     \/      `)
		println("")
	}

	app.Run(os.Args)

	flu_svc.StartFeedLineSync()

	api.Build()

}
