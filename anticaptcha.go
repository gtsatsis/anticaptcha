package anticaptcha

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type Api struct {
	ClientKey  string `json:"clientKey"` // ClientKey is the API key used to access the Anti-Captcha service.
	httpClient *http.Client
}

// NewAntiCaptchaApi creates a new Api with the client key provided.
func NewAntiCaptchaApi(clientKey string) *Api {
	return &Api{
		ClientKey:  clientKey,
		httpClient: &http.Client{},
	}
}

// SubmitTask submits a Task object for solving.
// The Task object is assigned its ID.
func (a *Api) SubmitTask(t *Task) error {
	j, err := t.GetJson()
	if err != nil {
		return err
	}

	j["type"] = string(t.Type)
	rbdm := map[string]interface{}{
		"clientKey": a.ClientKey,
		"task":      j,
	}

	b, err := json.Marshal(rbdm)
	if err != nil {
		return err
	}

	post, err := a.httpClient.Post("https://api.anti-captcha.com/createTask", "application/json", bytes.NewReader(b))
	if err != nil {
		return err
	}
	rbytes, err := io.ReadAll(post.Body)
	if err != nil {
		return err
	}

	rbody := &CreateTaskResponse{}
	err = json.Unmarshal(rbytes, rbody)
	if err != nil {
		return err
	}

	if rbody.ErrorId != 0 {
		return getApiError(rbody.ErrorCode)
	}

	t.ID = rbody.TaskId
	t.Status = TaskStatusProcessing
	t.lastChecked = time.Now().Unix()
	return nil
}

// GetTaskStatus retrieves the status of a Task object.
func (a *Api) GetTaskStatus(t *Task) (map[string]interface{}, error) {
	if t.ID == 0 || t.Status == TaskStatusNotSubmitted {
		return nil, ErrTaskNotSubmitted
	}

	if time.Now().Sub(time.Unix(t.lastChecked, 0)) < 3*time.Second {
		return nil, ErrCheckingTooFast
	}

	b, err := json.Marshal(map[string]interface{}{
		"clientKey": a.ClientKey,
		"taskId":    t.ID,
	})
	if err != nil {
		return nil, err
	}

	post, err := a.httpClient.Post("https://api.anti-captcha.com/getTaskResult", "application/json", bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	rbytes, err := io.ReadAll(post.Body)
	if err != nil {
		return nil, err
	}

	rbody := &GetTaskResultResponse{}
	err = json.Unmarshal(rbytes, rbody)
	if err != nil {
		return nil, err
	}
	if rbody.ErrorId != 0 {
		return nil, getApiError(rbody.ErrorCode)
	}

	if rbody.Status == string(TaskStatusProcessing) {
		t.lastChecked = time.Now().Unix()
		return nil, ErrTaskNotComplete
	}

	if rbody.Status == string(TaskStatusReady) {
		return rbody.Solution, nil
	}

	return nil, nil
}

func getApiError(errorCode string) error {
	switch errorCode {
	case "ERROR_KEY_DOES_NOT_EXIST":
		return ErrApiKeyInvalid
	case "ERROR_NO_SLOT_AVAILABLE":
		return ErrNoAvailableWorkers
	case "ERROR_ZERO_BALANCE":
		return ErrZeroBalance
	case "ERROR_NO_SUCH_CAPCHA_ID":
		return ErrCaptchaIdExpired
	default:
		return nil
	}

}
