package cmd

import "github.com/spf13/cobra"

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "identify",
		Short: "identify the carbionio device",
		Long: `identify flashes the LEDs of the carbionio parent device, making it
easier to recognize which device is being controlled.`,
		Run: identify,
	})
}

func identify(cmd *cobra.Command, args []string) {
	// // LEDs.
	// ls := []*leds.LED{leds.Power, leds.Status, leds.Mute}
	// for _, l := range ls {
	//  fmt.Println(l)
	// }
	// fmt.Println("Toggling LEDs…")
	// for _, l := range ls {
	//  l.SetState(leds.Off)
	// }
	// time.Sleep(2 * time.Second)
	// for _, l := range ls {
	//  l.SetState(leds.On)
	// }
}
