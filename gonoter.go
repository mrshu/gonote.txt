package main

import  (
        "fmt"
        "os"
        "github.com/spf13/cobra"
        "github.com/rakyll/globalconf"
)

func main() {

        conf, _ := globalconf.New("gonoter")

        var GonoterCmd = &cobra.Command{
            Use:   "gonoter",
            Short: "gonoter is a go implementation of note.txt specification.",
            Long: `A small, fast and fun implementation of note.txt`,
            Run: func(cmd *cobra.Command, args []string) {
            },
        }
}
