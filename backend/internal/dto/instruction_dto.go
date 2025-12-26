package dto

type CreateInstructionRequest struct {
	StepNumber int    `json:"step_number" binding:"required"`
	Text       string `json:"text" binding:"required"`
}

type UpdateInstructionRequest struct {
	StepNumber int    `json:"step_number" binding:"required"`
	Text       string `json:"text" binding:"required"`
}
