package aof


type OperationType string

var (
	PutOp OperationType = "Put"
	DeleteOp OperationType = "Delete"
)

type Operation struct {
	OperationType OperationType `json:"operationType"`
	Key string `json:"key"`
	Value interface{} `json:"value"`
	Timestamp string `json:"timestamp"`
}

