package main

import  (
        "fmt"
        "github.com/spf13/cobra"
        "github.com/rakyll/globalconf"
        "../go-notetxt"
        "flag"
        "os/user"
        "time"
        "os"
        "os/exec"
        "strings"
        "strconv"
)

func openFileInEditor(file string) {
        editor := os.Getenv("EDITOR")
        if len(editor) == 0 {
                editor = "nano" //FIXME: saner default?
        }

        c := exec.Command(editor, file)

        // nasty hack, see http://stackoverflow.com/a/12089980
        c.Stdin = os.Stdin
        c.Stdout = os.Stdout
        c.Stderr = os.Stderr

        er := c.Run()

        if er != nil {
                fmt.Println(er.Error())
                panic(er)
        }
}

func main() {

        conf, _ := globalconf.New("gonote")

        var flagNotedir = flag.String("dir", "", "Location of the note.txt directory.")
        var dir string

        var today bool

        var cmdAdd = &cobra.Command{
            Use:   "add [title] [tag]",
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

                openFileInEditor(file)
            },
        }
        cmdAdd.Flags().BoolVarP(&today, "today", "T", false,
                                 "Add today's journal entry.")

        var cmdList = &cobra.Command{
            Use:   "list",
            Short: "List notes.",
            Long:  `List all valid note files in the directory.`,
            Run: func(cmd *cobra.Command, args []string) {
                notes, err := notetxt.ParseDir(dir)
                if err != nil {
                    panic(err)
                }

                for i, note := range notes {
                    fmt.Printf("%d %s - %v\n", i, note.Name, note.Tags)
                }
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
        GonoterCmd.Execute()
}
