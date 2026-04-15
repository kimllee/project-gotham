package controllers

import (
	"errors"
	"net/http"
	"os"
	"project-gotham/config"
	"project-gotham/models"
	"strconv"

	"github.com/disintegration/imaging"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Gère quelle fonction va s'exécuter en fonction de quel appel à quelle URL
// Cette fonction prend en entrée le contexte de Gin.

// GetProjects godoc
// @Description Récupération de tous les projets
// @Tags Projects
// @Produce json
// @Success 200 {array} models.Project
// @Security BearerAuth
// @Router /projects [get]
func GetProjects(c *gin.Context) {
	// projects := []models.Project{
	// 	{
	// 		Name:        "Projet 1",
	// 		Description: "Description projet 1",
	// 		Image:       "background.webp",
	// 		Skills:      []string{"Go", "mysql", "Docker"},
	// 	},
	// 	{
	// 		Name:        "Projet 2",
	// 		Description: "Description projet 1",
	// 		Image:       "background.webp",
	// 		Skills:      []string{"Go", "mysql"},
	// 	},
	// 	{
	// 		Name:        "Projet 3",
	// 		Description: "Description projet 1",
	// 		Image:       "background.webp",
	// 		Skills:      []string{"Go", "mysql", "JS", "Kubernetes"},
	// 	},
	// }
	var projects []models.Project // On utilise la structure Project créée précédemment

	if err := config.DB.Preload("Comments").Preload("Likes").Find(&projects).Error; err != nil { // Preload charge tous les commentaires liés à chaque projet.
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de récupérer les projets."})
		return // Termine la fonction
	}
	c.JSON(http.StatusOK, projects)
}

func GetProject(c *gin.Context) {
	var project models.Project

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam) // Conversion string to integer

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide."}) // Header http
		return
	}
	if err := config.DB.Preload("Comments").Preload("Likes").First(&project, id).Error; err != nil { // Preload charge tous les commentaires liés au projet sélectionné.
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Projet introuvable."})
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de récupérer le projet."})
	}
	c.JSON(http.StatusOK, project)
}

func PostProject(c *gin.Context) {
	var project models.Project

	// Contrôle de saisie
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides."})
		return
	}

	file, err := c.FormFile("image")

	if err == nil {
		path := "uploads/" + file.Filename
		if err := c.SaveUploadedFile(file, path); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Impossible d'enregistrer l'image."})
			return
		}
		img, _ := imaging.Open(path)                            // L'underscore permet de ne pas utiliser la variable erreur car on l'a déjà traitée juste avant.
		resized := imaging.Resize(img, 800, 0, imaging.Lanczos) // Lanczos format de compression
		if err := imaging.Save(resized, path); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Impossible d'enregistrer l'image redimensionnée."})
			return
		}

		project.Image = path
	}

	// & = adresse mémoire | * = données
	if err := config.DB.Create(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création du projet."})
		return
	}
	c.JSON(http.StatusCreated, project)
}

func PutProject(c *gin.Context) {
	var project models.Project

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)

	// Vérifier si l'ID est au bon format
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID non valide."})
		return
	}

	// Vérifier si le projet existe en base de données
	// Utilise le premier objet de la bdd trouvé et compare les ids. On récupère une potentielle erreur avec le ".Error qui sera envoyée dans "err".
	if err := config.DB.First(&project, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Projet introuvable."})
		return
	}

	// Récupération de tous les inputs envoyés par l'utilisateur
	var input models.ProjectUpdateInput

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format de données incorrect."})
		return
	}

	// Variable intermédiaire
	updates := make(map[string]interface{})

	if input.Name != nil {
		updates["Name"] = *input.Name
	}

	if input.Description != nil {
		updates["Description"] = *input.Description
	}

	file, err := c.FormFile("image")

	if err == nil {
		path := "uploads/" + file.Filename
		if err := c.SaveUploadedFile(file, path); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Impossible d'enregistrer l'image."})
			return
		}
		img, _ := imaging.Open(path)                            // L'underscore permet de ne pas utiliser la variable erreur car on l'a déjà traitée juste avant.
		resized := imaging.Resize(img, 800, 0, imaging.Lanczos) // Lanczos format de compression
		if err := imaging.Save(resized, path); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Impossible d'enregistrer l'image redimensionnée."})
			return
		}

		if project.Image != "" {
			_ = os.Remove(project.Image)
		}
		updates["image"] = path
	}

	if input.Skills != nil {
		updates["Skills"] = datatypes.JSONSlice[string](*input.Skills)
	}

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Aucune mise à jour détectée."})
		return
	}

	// On met en bdd ce qu'il y a dans updates
	if err := config.DB.Model(&project).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la mise à jour du projet."})
		return
	}
	c.JSON(http.StatusOK, project)
}

func DeleteProject(c *gin.Context) {
	var project models.Project

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)

	// Vérifier si l'ID est au bon format
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID non valide."})
		return
	}

	// Vérifier si le projet existe en base de données
	// Utilise le premier objet de la bdd trouvé et compare les ids. On récupère une potentielle erreur avec le ".Error qui sera envoyée dans "err".
	if err := config.DB.First(&project, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Projet introuvable."})
		return
	}

	// Supprimer en bdd un projet
	// &projet -> récupération de l'adresse mémoire et non pas la valeur.
	if err := config.DB.Delete(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la suppression du projet."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Le projet a été supprimé avec succès."})
}

func LikedProjects(c *gin.Context) {
	var project models.Project
	var user models.User

	idParam := c.Param("id")
	projectId, err := strconv.Atoi(idParam)

	// Vérifier si l'ID est au bon format
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID non valide."})
		return

	}
	if err := config.DB.Preload("Likes").First(&project, projectId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Projet introuvable."})
		return
	}

	// Récupération de l'ID de l'utilisateur
	userID, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "L'ID de l'utilisateur est introuvable."})
		return
	}

	userIDInt, ok := userID.(int) // Conversion en int

	if err := config.DB.First(&user, uint(userIDInt)).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Utilisateur introuvable."})
		return
	}

	liked := false
	for _, u := range project.Likes {
		if u.ID == user.ID {
			liked = true
			break
		}
	}

	if liked {
		if err := config.DB.Model(&project).Association("Likes").Delete(&user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de retirer le Like."})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Like retiré."})
	} else {
		if err := config.DB.Model(&project).Association("Likes").Append(&user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible d'ajouter le like."})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Like ajouté."})
	}
}
