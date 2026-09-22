package constants

// PetSpecies enumerates animal kinds.
const (
	PetSpeciesDog    = "dog"
	PetSpeciesCat    = "cat"
	PetSpeciesRabbit = "rabbit"
	PetSpeciesOther  = "other"
)

// ValidPetSpecies returns all accepted species.
func ValidPetSpecies() []string {
	return []string{PetSpeciesDog, PetSpeciesCat, PetSpeciesRabbit, PetSpeciesOther}
}

// IsValidPetSpecies reports whether the species is known.
func IsValidPetSpecies(s string) bool {
	for _, v := range ValidPetSpecies() {
		if v == s {
			return true
		}
	}
	return false
}

// PetStatus enumerates pet availability.
const (
	PetStatusAvailable = "available"
	PetStatusPending   = "pending"
	PetStatusAdopted   = "adopted"
)

// ValidPetStatuses returns all accepted pet statuses.
func ValidPetStatuses() []string {
	return []string{PetStatusAvailable, PetStatusPending, PetStatusAdopted}
}

// IsValidPetStatus reports whether the status is known.
func IsValidPetStatus(s string) bool {
	for _, v := range ValidPetStatuses() {
		if v == s {
			return true
		}
	}
	return false
}
