package task

type URIParam struct {
	ID int `uri:"id" binding:"required"`
}
type CreateRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}
type CreateResponse struct {
	ID int `json:"id"`
}
type UpdateRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	Completed   bool   `json:"completed" binding:"required"`
}
