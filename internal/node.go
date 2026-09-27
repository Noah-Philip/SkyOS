package internal

import "math"
type NodeStatus string


type Position struct {
	X float64
	Y float64
	Z float64
}

const (
	NodeReady NodeStatus = "READY"
	NodeBusy NodeStatus = "BUSY"
	NodeDead NodeStatus = "DEAD"
)

type Resources struct {
	Battery float64
	// CPU_Usage float64
	// Memory_Usage float64 
}

type Node struct {

	ID string
	
	//Position
	Position Position

	//Resource Usage
	Resources Resources

	Status NodeStatus
	
	Tasks []*Task
}

func (n *Node) isEligibleForTask(task *Task) bool {
	if(n.Status != NodeReady || n.Resources.Battery < task.MinRequirements.Battery) {
		return false;
	}

	return true;
}

func (n *Node) calculateScore(task *Task) float64 {
	score := -1 * Distance3D(n.Position, task.Position) / (n.Resources.Battery - task.MinRequirements.Battery) - float64(len(n.Tasks)) * 10
	return score;
}

// Distance3D calculates the distance between two 3D points
func Distance3D(p1, p2 Position) float64 {
	return math.Sqrt(math.Pow(p2.X-p1.X, 2) + math.Pow(p2.Y-p1.Y, 2) + math.Pow(p2.Z-p1.Z, 2))
}