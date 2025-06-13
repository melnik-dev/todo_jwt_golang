package task

type URIParam struct {
	ID int `uri:"id" binding:"required,min=1"`
}

type CreateRequest struct {
	Title       string `json:"title" binding:"required,min=1,max=100"`
	Description string `json:"description" binding:"max=500"`
}

type CreateResponse struct {
	ID int `json:"id"`
}

type UpdateRequest struct {
	Title       string `json:"title" binding:"required,min=1,max=100"`
	Description string `json:"description" binding:"max=500"`
	Completed   bool   `json:"completed" binding:"required"`
}
