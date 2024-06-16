package main

import "fmt"

// Lista de tareas
type TaskList struct {
	tasks []*Task
}

func (tl *TaskList) appendTask(t *Task) {
	tl.tasks = append(tl.tasks, t)
}

func (tl *TaskList) removeTask(index int) {
	tl.tasks = append(tl.tasks[:index], tl.tasks[index+1:]...)
}

// Tareas
type Task struct {
	name      string
	desc      string
	completed bool
}

func (t *Task) toPrint() {
	fmt.Printf("Nombre: %s\nDescripción: %s\nCompletado: %t\n", t.name, t.desc, t.completed)
}

func (t *Task) markCompleted() {
	t.completed = true
}

func main() {

	t1 := Task{
		name:      "Curso de GO",
		desc:      "Completar curso de GO este mes",
		completed: false,
	}

	t2 := Task{
		name:      "Curso de HTML",
		desc:      "Completar curso de HTML esta semana",
		completed: true,
	}

	lista := TaskList{}
	lista.appendTask(&t1)
	lista.appendTask(&t2)

	fmt.Println(lista)

	lista.removeTask(1)

	t3 := Task{
		name:      "Curso de CSS",
		desc:      "Completar curso de CSS esta semana",
		completed: false,
	}

	lista.appendTask(&t3)
	fmt.Println(lista)

	for i, task := range lista.tasks {
		fmt.Println(i, task.name)
	}

	//t1.toPrint()
	//t2.toPrint()

}
