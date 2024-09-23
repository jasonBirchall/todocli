package main

import (
	"fmt"
	"os"

	"github.com/HxX2/todo/pkg/todo"
	"github.com/spf13/cobra"
)

func main() {
	t := todo.Init()

	var newTask string
	var remTaskNum int
	var toggleTaskNum int
	var editFlag bool
	var listDone bool
	var listUndone bool
	var hideProgress bool

	var rootCmd = &cobra.Command{
		Use:   "todo",
		Short: "Todo CLI application",
		Run: func(cmd *cobra.Command, args []string) {
			t.ListDone = !listUndone
			t.ListUndone = !listDone
			t.ShowProgress = !hideProgress

			switch {
			case remTaskNum != 0:
				t.RemTask(remTaskNum)
			case newTask != "":
				t.AddTask(newTask)
			case toggleTaskNum != 0:
				t.ToggleTask(toggleTaskNum)
			case editFlag:
				t.OpenEditor()
			case t.ListUndone:
				t.PrintList()
			case t.ListDone:
				t.PrintList()
			default:
				t.PrintList()
			}
		},
	}

	rootCmd.Flags().StringVarP(&newTask, "add", "a", "", "add a task")
	rootCmd.Flags().IntVarP(&remTaskNum, "remove", "r", 0, "remove a task")
	rootCmd.Flags().IntVarP(&toggleTaskNum, "toggle", "t", 0, "toggle done for a task")
	rootCmd.Flags().BoolVarP(&editFlag, "edit", "e", false, "edit todo file")
	rootCmd.Flags().BoolVar(&listDone, "list-done", false, "list done tasks")
	rootCmd.Flags().BoolVar(&listUndone, "list-undone", false, "list undone tasks")
	rootCmd.Flags().BoolVar(&hideProgress, "hide-progress", false, "hide progress bar")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
