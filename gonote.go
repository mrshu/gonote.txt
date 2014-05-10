package main

import  (
        "fmt"
        "os"
        "github.com/spf13/cobra"
        "github.com/rakyll/globalconf"
        "flag"
)

func main() {

        conf, _ := globalconf.New("gonoter")

        var flagNotedir = flag.String("dir", "", "Location of the note.txt directory.")
        var dir string

        var cmdAdd = &cobra.Command{
            Use:   "add [title] [category]",
            Short: "Add a note.",
            Long:  `Add a note into a category.`,
            Run: func(cmd *cobra.Command, args []string) {
                if len(args) < 1 {
                        fmt.Println("I need something to add")
                        return
                }
            },
        }

        var GonoterCmd = &cobra.Command{
            Use:   "gonoter",
            Short: "gonoter is a go implementation of note.txt specification.",
            Long: `A small, fast and fun implementation of note.txt`,
            Run: func(cmd *cobra.Command, args []string) {
            },
        }

        GonoterCmd.PersistentFlags().StringVarP(&dir, "directory", "", "",
                                     "Location of the note.txt directory.")

        if dir == "" {
                if *flagNotedir == "" {
                        dir = "~/notes"
                } else {
                        dir = *flagNotedir
                }
        }
        conf.ParseAll()

        GonoterCmd.AddCommand(cmdAdd)
        GonoterCmd.Execute()
}
