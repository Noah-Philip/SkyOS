package internal

type Task struct {
	ID string

	//Position
	Position Position

	//Resource Usage
	MinRequirements Resources
}
