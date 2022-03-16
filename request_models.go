package anticaptcha

type CreateTaskResponse struct {
	ErrorId          int    `json:"errorId"`
	ErrorCode        string `json:"errorCode"`
	ErrorDescription string `json:"errorDescription"`
	TaskId           int    `json:"taskId"`
}

type GetTaskResultResponse struct {
	ErrorId          int                    `json:"errorId"`
	ErrorCode        string                 `json:"errorCode"`
	ErrorDescription string                 `json:"errorDescription"`
	Status           string                 `json:"status"`
	Solution         map[string]interface{} `json:"solution"`
	Cost             string                 `json:"cost"`
	Ip               string                 `json:"ip"`
	CreateTime       int                    `json:"createTime"`
	EndTime          int                    `json:"endTime"`
	SolveCount       int                    `json:"solveCount"`
}
