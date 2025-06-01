package controllers

import (
	"errors"
	"gbvmis/internals/models"
	"gbvmis/internals/repository"
	"gbvmis/internals/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type SuspectController struct {
	repo repository.SuspectRepository
}

func NewSuspectController(repo repository.SuspectRepository) *SuspectController {
	return &SuspectController{repo: repo}
}

// ==========================

// CreateSuspect godoc
//
//	@Summary		Create a new suspect record with photo and fingerprints upload
//	@Description	Creates a new suspect entry with photo and fingerprints files.
//	@Tags			Suspects
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			first_name		formData	string	true	"First Name"
//	@Param			middle_name		formData	string	false	"Middle Name"
//	@Param			last_name		formData	string	true	"Last Name"
//	@Param			dob				formData	string	false	"DOB in YYYY-MM-DD format"
//	@Param			gender			formData	string	false	"Gender"
//	@Param			phone_number	formData	string	false	"Phone Number"
//	@Param			nin				formData	string	false	"NIN"
//	@Param			nationality		formData	string	false	"Nationality"
//	@Param			address			formData	string	false	"Address"
//	@Param			occupation		formData	string	false	"Occupation"
//	@Param			status			formData	string	false	"Status"
//	@Param			created_by		formData	string	true	"Created By"
//	@Param			updated_by		formData	string	false	"Updated By"
//	@Param			photo			formData	file	false	"Photo file upload"
//	@Param			fingerprints	formData	file	false	"Fingerprints file upload"
//	@Success		201				{object}	fiber.Map	"Successfully created suspect record"
//	@Failure		400				{object}	fiber.Map	"Bad request due to invalid input or file error"
//	@Failure		500				{object}	fiber.Map	"Server error when creating suspect"
//	@Router			/suspect [post]
func (h *SuspectController) CreateSuspect(c *fiber.Ctx) error {
	suspect := models.Suspect{}

	// Parse form values
	suspect.FirstName = c.FormValue("first_name")
	suspect.MiddleName = c.FormValue("middle_name")
	suspect.LastName = c.FormValue("last_name")
	dobStr := c.FormValue("dob")
	if dobStr != "" {
		parsedDob, err := utils.ParseDate(dobStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid date format for dob, expected YYYY-MM-DD",
				"data":    err.Error(),
			})
		}
		suspect.Dob = parsedDob // Store as string if needed, or change Suspect struct to time.Time
	}
	suspect.Gender = c.FormValue("gender")
	suspect.PhoneNumber = c.FormValue("phone_number")
	suspect.Nin = c.FormValue("nin")
	suspect.Nationality = c.FormValue("nationality")
	suspect.Address = c.FormValue("address")
	suspect.Occupation = c.FormValue("occupation")
	suspect.Status = c.FormValue("status")
	suspect.CreatedBy = c.FormValue("created_by")
	suspect.UpdatedBy = c.FormValue("updated_by")

	// Validate required fields
	if suspect.FirstName == "" || suspect.LastName == "" || suspect.CreatedBy == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "first_name, last_name, and created_by fields are required",
		})
	}

	// Handle photo file upload
	photoFile, err := c.FormFile("photo")
	if err == nil && photoFile != nil {
		photoOpened, err := photoFile.Open()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to open photo file",
				"data":    err.Error(),
			})
		}
		defer photoOpened.Close()

		photoBytes := make([]byte, photoFile.Size)
		_, err = photoOpened.Read(photoBytes)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to read photo file",
				"data":    err.Error(),
			})
		}
		suspect.Photo = photoBytes
	}

	// Handle fingerprints file upload
	fpFile, err := c.FormFile("fingerprints")
	if err == nil && fpFile != nil {
		fpOpened, err := fpFile.Open()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to open fingerprints file",
				"data":    err.Error(),
			})
		}
		defer fpOpened.Close()

		fpBytes := make([]byte, fpFile.Size)
		_, err = fpOpened.Read(fpBytes)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to read fingerprints file",
				"data":    err.Error(),
			})
		}
		suspect.Fingerprints = fpBytes
	}

	// Save suspect to database
	if err := h.repo.CreateSuspect(&suspect); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create suspect",
			"data":    err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Suspect created successfully",
		"data":    suspect,
	})
}

// ==================

// UpdateSuspect godoc
//
//	@Summary		Update a suspect record by ID
//	@Description	Partially updates fields of an existing suspect, including photo and fingerprints.
//	@Tags			Suspects
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			id				path		string	true	"Suspect ID"
//	@Param			first_name		formData	string	false	"First Name"
//	@Param			middle_name		formData	string	false	"Middle Name"
//	@Param			last_name		formData	string	false	"Last Name"
//	@Param			dob				formData	string	false	"DOB in YYYY-MM-DD format"
//	@Param			gender			formData	string	false	"Gender"
//	@Param			phone_number	formData	string	false	"Phone Number"
//	@Param			nin				formData	string	false	"NIN"
//	@Param			nationality		formData	string	false	"Nationality"
//	@Param			address			formData	string	false	"Address"
//	@Param			occupation		formData	string	false	"Occupation"
//	@Param			status			formData	string	false	"Status"
//	@Param			updated_by		formData	string	true	"Updated By"
//	@Param			photo			formData	file		false	"Photo file upload"
//	@Param			fingerprints	formData	file		false	"Fingerprints file upload"
//	@Success		200				{object}	fiber.Map	"Suspect updated successfully"
//	@Failure		400				{object}	fiber.Map	"Bad request due to invalid input or file error"
//	@Failure		404				{object}	fiber.Map	"Suspect not found"
//	@Failure		500				{object}	fiber.Map	"Server error when updating suspect"
//	@Router			/suspect/{id} [put]
func (h *SuspectController) UpdateSuspect(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Suspect ID is required",
		})
	}

	// Verify suspect exists
	_, err := h.repo.GetSuspectByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Suspect not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve suspect",
			"data":    err.Error(),
		})
	}

	updates := make(map[string]interface{})

	// Helper func to add non-empty form value to updates
	addUpdate := func(key string) {
		val := c.FormValue(key)
		if val != "" {
			updates[key] = val
		}
	}

	addUpdate("first_name")
	addUpdate("middle_name")
	addUpdate("last_name")
	addUpdate("gender")
	addUpdate("phone_number")
	addUpdate("nin")
	addUpdate("nationality")
	addUpdate("address")
	addUpdate("occupation")
	addUpdate("status")

	// UpdatedBy is required for audit
	updatedBy := c.FormValue("updated_by")
	if updatedBy == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "updated_by field is required",
		})
	}
	updates["updated_by"] = updatedBy

	// Parse dob separately
	dobStr := c.FormValue("dob")
	if dobStr != "" {
		parsedDob, err := utils.ParseDate(dobStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid date format for dob, expected YYYY-MM-DD",
				"data":    err.Error(),
			})
		}
		updates["dob"] = parsedDob.Format("2006-01-02")
	}

	// Handle photo file upload
	photoFile, err := c.FormFile("photo")
	if err == nil && photoFile != nil {
		photoOpened, err := photoFile.Open()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to open photo file",
				"data":    err.Error(),
			})
		}
		defer photoOpened.Close()

		photoBytes := make([]byte, photoFile.Size)
		_, err = photoOpened.Read(photoBytes)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to read photo file",
				"data":    err.Error(),
			})
		}
		updates["photo"] = photoBytes
	}

	// Handle fingerprints file upload
	fpFile, err := c.FormFile("fingerprints")
	if err == nil && fpFile != nil {
		fpOpened, err := fpFile.Open()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to open fingerprints file",
				"data":    err.Error(),
			})
		}
		defer fpOpened.Close()

		fpBytes := make([]byte, fpFile.Size)
		_, err = fpOpened.Read(fpBytes)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to read fingerprints file",
				"data":    err.Error(),
			})
		}
		updates["fingerprints"] = fpBytes
	}

	// Check if anything to update
	if len(updates) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "No valid fields or files provided for update",
		})
	}

	// Perform update in database via repository
	if err := h.repo.UpdateSuspect(id, updates); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update suspect",
			"data":    err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Suspect updated successfully",
		"data":    updates,
	})
}

// GetAllSuspects godoc
//
//	@Summary		Retrieve a paginated list of suspects
//	@Description	Fetches all suspect records with pagination support.
//	@Tags			Suspects
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	fiber.Map	"Suspects retrieved successfully"
//	@Failure		500		{object}	fiber.Map	"Failed to retrieve suspects"
//	@Router			/suspects [get]
func (h *SuspectController) GetAllSuspects(c *fiber.Ctx) error {
	pagination, suspects, err := h.repo.GetPaginatedSuspects(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve suspects",
			"data":    err.Error(),
		})
	}

	// Return the paginated response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Suspects retrieved successfully",
		"data":    suspects,
		"pagination": fiber.Map{
			"total_items":  pagination.TotalItems,
			"total_pages":  pagination.TotalPages,
			"current_page": pagination.CurrentPage,
			"limit":        pagination.ItemsPerPage,
		},
	})
}

// =========

// GetSingleSuspect godoc
//
//	@Summary		Retrieve a single suspect record by ID
//	@Description	Fetches a suspect record based on the provided ID.
//	@Tags			Suspects
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string	true	"Suspect ID"
//	@Success		200		{object}	fiber.Map	"Suspect retrieved successfully"
//	@Failure		404		{object}	fiber.Map	"Suspect not found"
//	@Failure		500		{object}	fiber.Map	"Server error when retrieving suspect"
//	@Router			/suspect/{id} [get]
func (h *SuspectController) GetSingleSuspect(c *fiber.Ctx) error {
	// Get the Suspect ID from the route parameters
	id := c.Params("id")

	// Fetch the suspect by ID
	suspect, err := h.repo.GetSuspectByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Suspect not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve suspect",
			"data":    err.Error(),
		})
	}

	// Return the response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Suspect and associated data retrieved successfully",
		"data":    suspect,
	})
}

// ==================

// DeleteSuspectByID godoc
//
//	@Summary		Delete a suspect record by ID
//	@Description	Deletes a suspect record based on the provided ID.
//	@Tags			Suspects
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string	true	"Suspect ID"
//	@Success		200		{object}	fiber.Map	"Suspect deleted successfully"
//	@Failure		404		{object}	fiber.Map	"Suspect not found"
//	@Failure		500		{object}	fiber.Map	"Server error when deleting suspect"
//	@Router			/suspect/{id} [delete]
func (h *SuspectController) DeleteSuspectByID(c *fiber.Ctx) error {
	// Get the Suspect ID from the route parameters
	id := c.Params("id")

	// Find the Suspect in the database
	suspect, err := h.repo.GetSuspectByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{
				"status":  "error",
				"message": "Suspect not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to find Suspect",
			"data":    err.Error(),
		})
	}

	// Delete the Suspect
	if err := h.repo.DeleteByID(id); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to delete suspect",
			"data":    err.Error(),
		})
	}

	// Return success response
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Suspect deleted successfully",
		"data":    suspect,
	})
}

// =================

// SearchSuspects godoc
//
//	@Summary		Search for suspects with pagination
//	@Description	Retrieves a paginated list of suspects based on search criteria.
//	@Tags			Suspects
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	fiber.Map	"Suspects retrieved successfully"
//	@Failure		500		{object}	fiber.Map	"Failed to retrieve suspects"
//	@Router			/suspects/search [get]
func (h *SuspectController) SearchSuspects(c *fiber.Ctx) error {
	// Call the repository function to get paginated search results
	pagination, suspects, err := h.repo.SearchPaginatedSuspects(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve suspects",
			"data":    err.Error(),
		})
	}

	// Return the response with pagination details
	return c.Status(200).JSON(fiber.Map{
		"status":     "success",
		"message":    "Suspects retrieved successfully",
		"pagination": pagination,
		"data":       suspects,
	})
}
