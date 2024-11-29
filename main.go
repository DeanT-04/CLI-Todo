package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"github.com/fatih/color"
)

// Task represents a to-do list task with ID, description, and completion status
type Task struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

// LoadTasks reads tasks from a JSON file
func LoadTasks(filename string) ([]Task, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return []Task{}, nil
		}
		return nil, err
	}
	var tasks []Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

// SaveTasks writes tasks to a JSON file
func SaveTasks(filename string, tasks []Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

// AddTask adds a new task to the list
func AddTask(tasks []Task, description string) []Task {
	maxID := 0
	for _, task := range tasks {
		if task.ID > maxID {
			maxID = task.ID
		}
	}
	
	newTask := Task{
		ID:          maxID + 1,
		Description: description,
		Completed:   false,
	}
	return append(tasks, newTask)
}

// ListTasks prints all tasks with ASCII status indicators
func ListTasks(tasks []Task) {
	if len(tasks) == 0 {
		color.Yellow("📝 No tasks found.")
		return
	}

	title := color.New(color.FgHiCyan, color.Bold).SprintFunc()
	taskID := color.New(color.FgHiYellow).SprintFunc()
	completed := color.New(color.FgHiGreen).SprintFunc()
	pending := color.New(color.FgHiRed).SprintFunc()
	desc := color.New(color.FgHiWhite).SprintFunc()

	fmt.Printf("\n%s\n", title("📋 Your Tasks"))
	fmt.Printf("%s\n", title("============"))
	for _, task := range tasks {
		status := pending("[ ]")
		if task.Completed {
			status = completed("[✓]")
		}
		fmt.Printf("%s %s %s\n", taskID(fmt.Sprintf("%d.", task.ID)), status, desc(task.Description))
	}
	fmt.Println()
}

// MarkTaskComplete marks a task as complete
func MarkTaskComplete(tasks []Task, id int) []Task {
	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i].Completed = true
			fmt.Printf("Marked task %d as complete: %s\n", id, tasks[i].Description)
			return tasks
		}
	}
	fmt.Printf("Task with ID %d not found\n", id)
	return tasks
}

// DeleteTask deletes a task by ID
func DeleteTask(tasks []Task, id int) []Task {
	for i := range tasks {
		if tasks[i].ID == id {
			fmt.Printf("Deleted task: %s\n", tasks[i].Description)
			return append(tasks[:i], tasks[i+1:]...)
		}
	}
	fmt.Printf("Task with ID %d not found\n", id)
	return tasks
}

func showHelp() {
	title := color.New(color.FgHiCyan, color.Bold).SprintFunc()
	cmd := color.New(color.FgHiYellow).SprintFunc()
	desc := color.New(color.FgHiWhite).SprintFunc()

	fmt.Printf("\n%s\n", title("📝 TODO CLI - Task Manager"))
	fmt.Printf("%s\n", title("========================"))
	fmt.Println("\n📌 Commands:")
	fmt.Printf("  %s  %s\n", cmd("add <description>"), desc("- Add a new task"))
	fmt.Printf("  %s          %s\n", cmd("list"), desc("- List all tasks"))
	fmt.Printf("  %s     %s\n", cmd("complete <id>"), desc("- Mark a task as complete"))
	fmt.Printf("  %s       %s\n", cmd("delete <id>"), desc("- Delete a task"))
	fmt.Printf("  %s            %s\n", cmd("help"), desc("- Show this help message"))
	fmt.Printf("  %s            %s\n", cmd("exit"), desc("- Exit the program"))
	fmt.Println()
}

func main() {
	tasks, err := LoadTasks("tasks.json")
	if err != nil {
		color.Red("Error loading tasks: %v\n", err)
		return
	}

	reader := bufio.NewReader(os.Stdin)
	showHelp()

	prompt := color.New(color.FgHiGreen, color.Bold).SprintFunc()
	success := color.New(color.FgHiGreen).SprintFunc()
	errorMsg := color.New(color.FgHiRed).SprintFunc()

	for {
		fmt.Printf("%s ", prompt("todo>"))
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("%s\n", errorMsg(fmt.Sprintf("Error reading input: %v", err)))
			continue
		}

		input = strings.TrimSpace(input)
		args := strings.Fields(input)
		
		if len(args) == 0 {
			continue
		}

		command := args[0]

		switch command {
		case "add":
			if len(args) < 2 {
				fmt.Printf("%s\n", errorMsg("Error: Please provide a task description"))
				continue
			}
			description := strings.Join(args[1:], " ")
			tasks = AddTask(tasks, description)
			fmt.Printf("%s\n", success(fmt.Sprintf("✨ Added task: %s", description)))

		case "list":
			ListTasks(tasks)

		case "complete":
			if len(args) != 2 {
				fmt.Printf("%s\n", errorMsg("Error: Please provide a task ID"))
				continue
			}
			id, err := strconv.Atoi(args[1])
			if err != nil {
				fmt.Printf("%s\n", errorMsg(fmt.Sprintf("Error: Invalid task ID: %s", args[1])))
				continue
			}
			tasks = MarkTaskComplete(tasks, id)

		case "delete":
			if len(args) != 2 {
				fmt.Printf("%s\n", errorMsg("Error: Please provide a task ID"))
				continue
			}
			id, err := strconv.Atoi(args[1])
			if err != nil {
				fmt.Printf("%s\n", errorMsg(fmt.Sprintf("Error: Invalid task ID: %s", args[1])))
				continue
			}
			tasks = DeleteTask(tasks, id)

		case "help":
			showHelp()

		case "exit":
			if err := SaveTasks("tasks.json", tasks); err != nil {
				fmt.Printf("%s\n", errorMsg(fmt.Sprintf("Error saving tasks: %v", err)))
			}
			fmt.Printf("%s\n", success("👋 Goodbye!"))
			return

		default:
			fmt.Printf("%s\n", errorMsg(fmt.Sprintf("Unknown command: %s\nType 'help' for usage information", command)))
		}

		if err := SaveTasks("tasks.json", tasks); err != nil {
			fmt.Printf("%s\n", errorMsg(fmt.Sprintf("Error saving tasks: %v", err)))
		}
	}
}
