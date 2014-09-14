package main

import  (
        "fmt"
        "github.com/spf13/cobra"
        "github.com/rakyll/globalconf"
        "../go-notetxt"
        "flag"
        "os/user"
        "time"
        "strings"
        "strconv"
)


func main() {

        conf, _ := globalconf.New("gonote")

        var flagNotedir = flag.String("dir", "", "Location of the note.txt directory.")
        var dir string

        var today bool

        var cmdAdd = &cobra.Command{
            Use:   "add <title> [tag]",
            Short: "Add a note.",
            Long:  `Add a note and tag it.`,
            Run: func(cmd *cobra.Command, args []string) {
                if len(args) < 1 && !today {
                        fmt.Println("I need something to add.")
                        return
                }

                var text string
                t := time.Now().Local()

                if today {
                        text = fmt.Sprintf("Daily journal, date %s", t.Format("02. 01. 2006"))
                } else {
                        text = strings.Join(args, " ")
                }

                file, err := notetxt.CreateNote(text, t.Format("2006/01/"), dir)
                if err != nil {
                        panic(err);
                }

                notetxt.OpenFileInEditor(file)
            },
        }
        cmdAdd.Flags().BoolVarP(&today, "today", "T", false,
                                 "Add today's journal entry.")

        var cmdList = &cobra.Command{
            Use:   "ls <query>",
            Short: "List notes.",
            Long:  `List all valid note files in the directory.`,
            Run: func(cmd *cobra.Command, args []string) {
                notes, err := notetxt.ParseDir(dir)
                if err != nil {
                    panic(err)
                }

                needle := strings.Join(args, " ")

                for i, note := range notes {
                    if note.Matches(needle) {
                        fmt.Printf("%d %s - %v\n", i, note.Name, note.Tags)
                    }
                }
            },
        }

        var cmdEdit = &cobra.Command{
            Use:   "edit <id>|<selector>",
            Short: "Edit notes.",
            Long:  `Edit a note identified by either an ID or a selector.`,
            Run: func(cmd *cobra.Command, args []string) {
                if len(args) < 1 {
                        fmt.Println("Either a note ID or a selector is required.")
                        return
                }

                notes, err := notetxt.ParseDir(dir)
                if err != nil {
                    panic(err)
                }

                noteid, err := strconv.Atoi(args[0])
                if err != nil {
                        needle := strings.Join(args, " ")
                        filtered_notes := notes.FilterBy(needle)

                        if len(filtered_notes) == 1 {
                                notetxt.OpenFileInEditor(filtered_notes[0].Filename)
                        } else if len (filtered_notes) == 0 {
                                fmt.Printf("No notes matched your selector.")
                        } else {
                                fmt.Printf("Notes matching your selector:\n")
                                filtered_notes.Print()
                        }

                        return
                }

                if noteid > len(notes) || noteid < 0 {
                        fmt.Printf("Invalid note ID (%v)\n", noteid)
                        return
                }

                notetxt.OpenFileInEditor(notes[noteid].Filename)

            },
        }

        var cmdTag = &cobra.Command{
            Use:   "tag <noteid> <tag-name>",
            Short: "Attaches a tag to a note.",
            Long:  `Tags a note with a one or more tags.`,
            Run: func(cmd *cobra.Command, args []string) {
                if len(args) < 2 {
                        fmt.Printf("Too few arguments.")
                }

                notes, err := notetxt.ParseDir(dir)
                if err != nil {
                    panic(err)
                }

                noteid, err := strconv.Atoi(args[0])
                if err != nil {
                        fmt.Printf("Do you really consider that a number? %v\n", err)
                        return
                }

                if noteid > len(notes) || noteid < 0 {
                        fmt.Printf("Invalid note ID (%v)\n", noteid)
                        return
                }

                file := notes[noteid].Filename
                tag := args[1]
                err = notetxt.TagNote(file, tag, dir)
                if err != nil {
                        panic(err)
                }

            },
        }

        var GonoterCmd = &cobra.Command{
            Use:   "gonote",
            Short: "gonote is a go implementation of note.txt specification.",
            Long: `A small, fast and fun implementation of note.txt`,
            Run: func(cmd *cobra.Command, args []string) {
                cmdList.Run(cmd, args)
            },
        }

        GonoterCmd.PersistentFlags().StringVarP(&dir, "directory", "", "",
                                     "Location of the note.txt directory.")

        conf.ParseAll()
        if dir == "" {
                if *flagNotedir == "" {
                        usr, err := user.Current()
                        if err != nil {
                                panic(err)
                        }

                        dir = usr.HomeDir + "/notes"
                } else {
                        dir = *flagNotedir
                }
        }

        GonoterCmd.AddCommand(cmdAdd)
        GonoterCmd.AddCommand(cmdList)
        GonoterCmd.AddCommand(cmdTag)
        GonoterCmd.AddCommand(cmdEdit)
        GonoterCmd.Execute()
}
