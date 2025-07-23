package constants

const (
	//  Roles
	//  ----------------------------------------------------------------

	Admin = iota + 1
	User

	// Difficulty

	Easy = iota - 1
	Medium
	Hard

	//  Handler Constants
	//  ----------------------------------------------------------------

	AuthorizationHeader = "Authorization"
	HandlerName         = "HANDLER_NAME"
	ServiceName         = "SERVICE_NAME"

	//  JWT Token Constants
	//  ----------------------------------------------------------------

	Role   = "role"
	UserID = "user"
	Name   = "name"
	Email  = "email"
	Type   = "type"
	Exp    = "exp"

	//  Entities
	//  ----------------------------------------------------------------

	TaskEntity     = "task"
	SyllabusEntity = "syllabus"

	//  Operations
	//  ----------------------------------------------------------------

	Get    = "get"
	Create = "create"
	Update = "update"
	Delete = "delete"
)

func DifficultyToInt(difficulty string) int32 {
	switch difficulty {
	case "easy":
		return Easy
	case "medium":
		return Medium
	case "hard":
		return Hard
	default:
		return 0
	}
}

func DifficultyToString(difficulty int32) string {
	switch difficulty {
	case Easy:
		return "easy"
	case Medium:
		return "medium"
	case Hard:
		return "hard"
	default:
		return ""
	}
}
