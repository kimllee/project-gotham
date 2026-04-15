package models

import (
	"time"

	"gorm.io/datatypes"
)

// ID
// CreatedAt
// UpdatedAt
// Name
// Description
// Image
// Skills

// Gorm créé lui-même les tables si elles n'existent pas.
// Type des données qui seront stockées en bdd + des règles spécifiques PrimaryKey, json etc.
// Écriture en Pascal case pour rendre les données publiques et pour que Gorm les interprète et les réécrive en Snake case.
type Project struct {
	ID          uint `gorm:"primaryKey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Name        string `binding:"required"`
	Description string `binding:"required"`
	Image       string
	//	Skills      []string `gorm:"type:json"` // Tableau de chaînes stockées en format json dans la bdd
	// Gorm ne sait pas gérer les tableaux de strings directement donc on passe par datatypes qu'on a installé précédemment.
	Skills   datatypes.JSONSlice[string] `gorm:"type:json" swaggertype:"array,string"` // Tableau de chaînes stockées en format json dans la bdd
	Comments []Comment                   `gorm:"foreignKey:ProjectID"`
	Likes    []User                      `gorm:"many2many:project_likes"`
}

// Payload model - DTO
type ProjectUpdateInput struct {
	Name        *string   `json:"name"`
	Description *string   `json:"description"`
	Skills      *[]string `json:"skills"`
}
