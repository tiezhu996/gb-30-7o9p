package dto

// PetCreateRequest publishes a pet.
type PetCreateRequest struct {
	Name         string  `json:"name" binding:"required,max=64"`
	Species      string  `json:"species" binding:"required"`
	Breed        string  `json:"breed" binding:"omitempty,max=64"`
	Age          int     `json:"age"`
	Gender       string  `json:"gender" binding:"omitempty,max=8"`
	Size         string  `json:"size" binding:"omitempty,max=16"`
	City         string  `json:"city" binding:"omitempty,max=64"`
	Description  string  `json:"description"`
	Personality  string  `json:"personality" binding:"omitempty,max=255"`
	HealthStatus string  `json:"health_status" binding:"omitempty,max=128"`
	Neutered     bool    `json:"neutered"`
	Vaccinated   bool    `json:"vaccinated"`
	ImageURLs    string  `json:"image_urls"`
}

// PetStatusRequest changes a pet's status.
type PetStatusRequest struct {
	Status string `json:"status" binding:"required"`
}
