package main

import  (
        "fmt"
        "github.com/spf13/cobra"
        "github.com/rakyll/globalconf"
        "github.com/mrshu/go-notetxt"
        "flag"
        "os/user"
        "time"
        "os"
        "io/ioutil"
        "os/exec"
        "strings"
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

                if today {
                        t := time.Now().Local()
                        dir := fmt.Sprintf("%s/%s", *flagNotedir, t.Format("2006/01/"))
                        os.MkdirAll(dir, 755)

                        text := fmt.Sprintf("Daily journal, date %s", t.Format("02. 01. 2006"))
                        spacer := "\n" + strings.Repeat("=", len(text))
                        file := fmt.Sprintf("%s%s.rst", dir, notetxt.TitleToFilename(text))

                        if _, err := os.Stat(file); err == nil {
                                fmt.Println("gonote: Notefile for today already exists. " +
                                                "You can still edit it if you want.")
                                return
                        }

                        e := ioutil.WriteFile(file,
                                                []byte(text + spacer),
                                                0644)
                        if e != nil {
                                panic(e)
                        }

                        openFileInEditor(file)
                }
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
        GonoterCmd.Execute()
}
