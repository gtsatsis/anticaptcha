package anticaptcha

import "errors"

var (
	ErrCheckingTooFast    = errors.New("checks should be 3 seconds apart")
	ErrTaskNotSubmitted   = errors.New("task not yet submitted")
	ErrTaskNotComplete    = errors.New("task not yet complete")
	ErrZeroBalance        = errors.New("account has zero or negative balance")
	ErrApiKeyInvalid      = errors.New("invalid api key")
	ErrNoAvailableWorkers = errors.New("no available workers")
	ErrCaptchaIdExpired   = errors.New("captcha expired")
	ErrCaptchaUnsolvable  = errors.New("captcha could not be solved by 5 different workers")
)
