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
	taskListView := tview.NewList().
		SetSelectedBackgroundColor(tcell.ColorDarkOrange).
		SetHighlightFullLine(true).
		SetMainTextColor(tcell.ColorWhite)

	taskListView.SetBackgroundColor(tcell.ColorBlack)
	taskListView.SetBorder(true)
	taskListView.SetBorderColor(tcell.ColorDarkCyan)
	taskListView.SetTitle("Tasks")
	taskListView.SetTitleColor(tcell.ColorYellow)

	// Add Vim-like keybindings for navigating the list
	taskListView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'j': // Move down
			index := taskListView.GetCurrentItem()
			if index < taskListView.GetItemCount()-1 {
				taskListView.SetCurrentItem(index + 1)
			}
			return nil
		case 'k': // Move up
			index := taskListView.GetCurrentItem()
			if index > 0 {
				taskListView.SetCurrentItem(index - 1)
			}
			return nil
		}
		return event
	})
	// Create the main list for actions
	actionList := tview.NewList().
		SetSelectedBackgroundColor(tcell.ColorDarkCyan).
		SetHighlightFullLine(true).
		SetMainTextColor(tcell.ColorWhite)

	actionList.SetBackgroundColor(tcell.ColorBlack)
	actionList.SetMainTextColor(tcell.ColorWhite.TrueColor())
	actionList.SetBorder(true)
	actionList.SetBorderColor(tcell.ColorYellow)
	actionList.SetTitle("Actions")
	actionList.SetTitleColor(tcell.Color100)
	actionList.ShowSecondaryText(false)

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
		AddItem(actionList, 0, 1, true). // Left side: action list
		AddItem(taskListView, 0, 5, false)

	// Populate the action list with options
	actionList.
		AddItem("Add Task", "Add a new task", 'a', func() {
			// Create a new Flex layout for the form and task list
			form := tview.NewForm()
			form.
				AddInputField("Task", "", 20, nil, nil).
				AddButton("Add", func() {
					task := form.GetFormItemByLabel("Task").(*tview.InputField).GetText()
					if task != "" {
						t.AddTask(task)
						updateTaskList(t, taskListView)
					}
					app.SetRoot(mainLayout, true).SetFocus(actionList)
				}).
				AddButton("Cancel", func() {
					app.SetRoot(mainLayout, true).SetFocus(actionList)
				})

			// Show the form alongside the task list by using a Modal
			modalLayout := tview.NewFlex().
				AddItem(mainLayout, 0, 1, false). // Keep the main layout visible
				AddItem(form, 30, 1, true)        // Display the form on the right

			app.SetRoot(modalLayout, true).SetFocus(form) // Show the form with the layout
		}).
		AddItem("Remove Task", "Remove a task", 'r', func() {
			// Set focus to the taskListView after selecting "Remove Task"
			app.SetRoot(mainLayout, true).SetFocus(taskListView)

			// Capture "d" key press to delete the currently highlighted task
			taskListView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
				switch event.Rune() {
				case 'd': // Press "d" to delete the task
					index := taskListView.GetCurrentItem()
					if index < len(t.Tasks) {
						// Remove the selected task
						t.RemTask(index + 1)            // Tasks are 1-based
						updateTaskList(t, taskListView) // Refresh the task list display
						taskListView.SetCurrentItem(0)  // Reset the task list selection
						app.SetFocus(actionList)        // Return focus to the action list
					}
					return nil
				}
				return event
			})
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

func updateTaskList(t *todo.Todo, taskListView *tview.List) {
	taskListView.Clear() // Clear the current items in the list
	for i, task := range t.Tasks {
		status := "[red]Undone"
		if task.Done {
			status = "[green]Done"
		}
		// Add the task as an item to the List view
		taskListView.AddItem(fmt.Sprintf("%d. %s", i+1, task.Name), status, 0, nil)
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
