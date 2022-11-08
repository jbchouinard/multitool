package cmd

import (
	"fmt"
	"log"

	"github.com/gofrs/uuid"
	"github.com/jbchouinard/multitool/clip"
	"github.com/jbchouinard/multitool/history"
	"github.com/spf13/cobra"
)

var uuidV4 bool

var uuidCmd = &cobra.Command{
	Use:   "uuid",
	Short: "Generate a UUID",
	Run: func(cmd *cobra.Command, args []string) {
		var id uuid.UUID
		var err error
		if uuidV4 {
			id, err = uuid.NewV4()
		} else {
			id, err = uuid.NewV1()
		}
		if err != nil {
			log.Fatal(err)
		} else {
			fmt.Println(id)
			clip.WriteBytes(id.Bytes())
			history.Add("uuid", id.String())
		}
	},
}

func init() {
	uuidCmd.Flags().BoolVar(&uuidV4, "v4", false, "Generate v4 UUID (default: v1)")
	rootCmd.AddCommand(uuidCmd)
}
