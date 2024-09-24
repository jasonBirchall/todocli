package main

import (
	"fmt"

	"github.com/HxX2/todo/pkg/todo"
	"github.com/gdamore/tcell/v2"
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

	// Add Vim-like keybindings for navigating the list
	actionList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'j': // Move down
			index := actionList.GetCurrentItem()
			if index < actionList.GetItemCount()-1 {
				actionList.SetCurrentItem(index + 1)
			}
			return nil
		case 'k': // Move up
			index := actionList.GetCurrentItem()
			if index > 0 {
				actionList.SetCurrentItem(index - 1)
			}
			return nil
		}
		return event
	})
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
			// Create a list of tasks to remove
			removeTaskList := tview.NewList()
			for i, task := range t.Tasks {
				taskStatus := "[red]Undone"
				if task.Done {
					taskStatus = "[green]Done"
				}
				removeTaskList.AddItem(fmt.Sprintf("%d. %s [%s]", i+1, task.Name, taskStatus), "", 0, nil)
			}

			// Set Vim-like navigation for the task removal list
			removeTaskList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
				switch event.Key() {
				case tcell.KeyEnter: // Remove the selected task
					index := removeTaskList.GetCurrentItem()
					t.RemTask(index + 1)            // Tasks are 1-based, but the index is 0-based
					updateTaskList(t, taskListView) // Refresh task list view
					app.SetRoot(mainLayout, true).SetFocus(actionList)
					return nil
				case tcell.KeyRune: // Check for character inputs
					switch event.Rune() {
					case 'j': // Move down
						index := removeTaskList.GetCurrentItem()
						if index < removeTaskList.GetItemCount()-1 {
							removeTaskList.SetCurrentItem(index + 1)
						}
						return nil
					case 'k': // Move up
						index := removeTaskList.GetCurrentItem()
						if index > 0 {
							removeTaskList.SetCurrentItem(index - 1)
						}
						return nil
					case 'q': // Cancel removal
						app.SetRoot(mainLayout, true).SetFocus(actionList)
						return nil
					}
				}
				return event
			})

			app.SetRoot(removeTaskList, true).SetFocus(removeTaskList)

		}).
		AddItem("Toggle Task", "Toggle done for a task", 't', func() {
			// Create a list of tasks to toggle
			toggleTaskList := tview.NewList()
			for i, task := range t.Tasks {
				taskStatus := "[red]Undone"
				if task.Done {
					taskStatus = "[green]Done"
				}
				toggleTaskList.AddItem(fmt.Sprintf("%d. %s [%s]", i+1, task.Name, taskStatus), "", 0, nil)
			}

			// Set Vim-like navigation for the task toggling list
			toggleTaskList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
				switch event.Key() {
				case tcell.KeyEnter: // Toggle the selected task's status
					index := toggleTaskList.GetCurrentItem()
					t.ToggleTask(index + 1)         // Tasks are 1-based, but the index is 0-based
					updateTaskList(t, taskListView) // Refresh task list view
					app.SetRoot(mainLayout, true).SetFocus(actionList)
					return nil
				case tcell.KeyRune: // Check for character inputs
					switch event.Rune() {
					case 'j': // Move down
						index := toggleTaskList.GetCurrentItem()
						if index < toggleTaskList.GetItemCount()-1 {
							toggleTaskList.SetCurrentItem(index + 1)
						}
						return nil
					case 'k': // Move up
						index := toggleTaskList.GetCurrentItem()
						if index > 0 {
							toggleTaskList.SetCurrentItem(index - 1)
						}
						return nil
					case 'q': // Cancel toggling
						app.SetRoot(mainLayout, true).SetFocus(actionList)
						return nil
					}
				}
				return event
			})

			app.SetRoot(toggleTaskList, true).SetFocus(toggleTaskList)
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
