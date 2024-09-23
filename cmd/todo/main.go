package main

import (
	"fmt"
	"strconv"

	"github.com/HxX2/todo/pkg/todo"
	"github.com/rivo/tview"
)

func main() {
	t := todo.Init()
	app := tview.NewApplication()

	// Create a text view to display tasks
	taskListView := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft).
		SetChangedFunc(func() {
			app.Draw()
		})

	// Create the main list for actions
	actionList := tview.NewList()

	// Create a Flex layout to divide the screen
	mainLayout := tview.NewFlex().
		AddItem(actionList, 0, 1, true).   // Left side: action list
		AddItem(taskListView, 0, 2, false) // Right side: task list display

	// Populate the action list with options
	actionList.
		AddItem("Add Task", "Add a new task", 'a', func() {
			form := tview.NewForm()
			form.
				AddInputField("Task", "", 20, nil, nil).
				AddButton("Add", func() {
					task := form.GetFormItemByLabel("Task").(*tview.InputField).GetText()
					if task != "" {
						t.AddTask(task)
						updateTaskList(t, taskListView) // Refresh task list view
					}
					app.SetRoot(mainLayout, true).SetFocus(actionList)
				}).
				AddButton("Cancel", func() {
					app.SetRoot(mainLayout, true).SetFocus(actionList)
				})
			app.SetRoot(form, true).SetFocus(form)
		}).
		AddItem("Remove Task", "Remove a task", 'r', func() {
			form := tview.NewForm()
			form.
				AddInputField("Task Number", "", 20, nil, nil).
				AddButton("Remove", func() {
					taskNumStr := form.GetFormItemByLabel("Task Number").(*tview.InputField).GetText()
					taskNum, err := strconv.Atoi(taskNumStr)
					if err != nil {
						showErrorModal(app, "Invalid task number", mainLayout)
						return
					}
					err = t.RemTask(taskNum)
					if err != nil {
						showErrorModal(app, fmt.Sprintf("Error: %v", err), mainLayout)
					} else {
						updateTaskList(t, taskListView) // Refresh task list view
					}
					app.SetRoot(mainLayout, true).SetFocus(actionList)
				}).
				AddButton("Cancel", func() {
					app.SetRoot(mainLayout, true).SetFocus(actionList)
				})
			app.SetRoot(form, true).SetFocus(form)
		}).
		AddItem("Toggle Task", "Toggle done for a task", 't', func() {
			form := tview.NewForm()
			form.
				AddInputField("Task Number", "", 20, nil, nil).
				AddButton("Toggle", func() {
					taskNumStr := form.GetFormItemByLabel("Task Number").(*tview.InputField).GetText()
					taskNum, err := strconv.Atoi(taskNumStr)
					if err != nil {
						showErrorModal(app, "Invalid task number", mainLayout)
						return
					}
					err = t.ToggleTask(taskNum)
					if err != nil {
						showErrorModal(app, fmt.Sprintf("Error: %v", err), mainLayout)
					} else {
						updateTaskList(t, taskListView) // Refresh task list view
					}
					app.SetRoot(mainLayout, true).SetFocus(actionList)
				}).
				AddButton("Cancel", func() {
					app.SetRoot(mainLayout, true).SetFocus(actionList)
				})
			app.SetRoot(form, true).SetFocus(form)
		}).
		AddItem("Quit", "Quit the application", 'q', func() {
			app.Stop()
		})

	// Initial population of the task list
	updateTaskList(t, taskListView)

	// Start the application
	if err := app.SetRoot(mainLayout, true).Run(); err != nil {
		fmt.Printf("Error running application: %s\n", err)
	}
}

// Function to update the task list display
func updateTaskList(t *todo.Todo, taskListView *tview.TextView) {
	taskListView.Clear()
	fmt.Fprintln(taskListView, "[::b]Tasks:")
	for i, task := range t.Tasks {
		status := "[red]Undone"
		if task.Done {
			status = "[green]Done"
		}
		fmt.Fprintf(taskListView, "%d. %s [%s]\n", i+1, task.Name, status)
	}
}

// Function to show an error message in a modal
func showErrorModal(app *tview.Application, message string, previousRoot tview.Primitive) {
	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			app.SetRoot(previousRoot, true) // Go back to the previous layout
		})
	app.SetRoot(modal, true).SetFocus(modal)
}
