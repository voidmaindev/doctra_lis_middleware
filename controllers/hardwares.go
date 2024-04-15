package controllers

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/voidmaindev/doctra_lis_middleware/models"
	"github.com/voidmaindev/doctra_lis_middleware/store"
)

// @Summary Get all hardwares
// @Description Get all hardwares
// @Tags Hardwares
// @Accept json
// @Produce json
// @Success 200 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /hardwares [get]
func HardwaresGetAll(c *fiber.Ctx) error {
	hardwares := &[]models.Hardware{}
	res := store.DB.Find(hardwares)
	if res.Error != nil {
		c.Status(http.StatusInternalServerError).JSON(&fiber.Map{"success": false, "message": res.Error.Error()})
		return nil
	}

	c.Status(http.StatusOK).JSON(&fiber.Map{"success": true, "message": "All hardwares", "hardwares": hardwares})

	return nil
}

// @Summary Get hardware by ID
// @Description Get hardware by ID
// @Tags Hardwares
// @Accept json
// @Produce json
// @Param id path int true "Hardware ID"
// @Success 200 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /hardwares/{id} [get]
func HardwaresGetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		c.Status(http.StatusInternalServerError).JSON(&fiber.Map{"success": false, "message": "No ID entered"})
		return nil
	}

	hardware := &models.Hardware{}
	res := store.DB.First(hardware, id)
	if res.Error != nil {
		c.Status(http.StatusInternalServerError).JSON(&fiber.Map{"success": false, "message": "Record not found"})
		return nil
	}

	c.Status(http.StatusOK).JSON(&fiber.Map{"success": true, "message": "Hardware", "hardware": hardware})

	return nil
}

// @Summary Create a new hardware
// @Description Create a new hardware
// @Tags Hardwares
// @Accept json
// @Produce json
// @Param hardware body models.Hardware true "Hardware object"
// @Success 200 {object} fiber.Map
// @Failure 400 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /hardwares [post]
func HardwaresCreate(c *fiber.Ctx) error {
	hardware := &models.Hardware{}
	err := c.BodyParser(hardware)
	if err != nil {
		c.Status(http.StatusBadRequest).JSON(&fiber.Map{"success": false, "message": err.Error()})
		return nil
	}

	res := store.DB.Create(hardware)
	if res.Error != nil {
		c.Status(http.StatusInternalServerError).JSON(&fiber.Map{"success": false, "message": res.Error.Error()})
		return nil
	}

	c.Status(http.StatusOK).JSON(&fiber.Map{"success": true, "message": "Created", "hardware": hardware})

	return nil
}

// @Summary Update hardware by ID
// @Description Update hardware by ID
// @Tags Hardwares
// @Accept json
// @Produce json
// @Param id path int true "Hardware ID"
// @Param hardware body models.Hardware true "Updated hardware object"
// @Success 200 {object} fiber.Map
// @Failure 400 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /hardwares/{id} [put]
func HardwaresUpdate(c *fiber.Ctx) error {
	idStr := c.Params("id")
	if idStr == "" {
		c.Status(http.StatusInternalServerError).JSON(&fiber.Map{"success": false, "message": "No ID entered"})
		return nil
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.Status(http.StatusBadRequest).JSON(&fiber.Map{"success": false, "message": err.Error()})
		return nil
	}

	hardware := &models.Hardware{}
	err = c.BodyParser(hardware)
	if err != nil {
		c.Status(http.StatusBadRequest).JSON(&fiber.Map{"success": false, "message": err.Error()})
		return nil
	}
	hardware.ID = uint(id)

	hardwareOrig := &models.Hardware{}
	res := store.DB.First(hardwareOrig, id)
	if res.Error != nil {
		c.Status(http.StatusInternalServerError).JSON(&fiber.Map{"success": false, "message": "Record not found"})
		return nil
	}

	res = store.DB.Save(hardware)
	if res.Error != nil {
		c.Status(http.StatusInternalServerError).JSON(&fiber.Map{"success": false, "message": res.Error.Error()})
		return nil
	}

	c.Status(http.StatusOK).JSON(&fiber.Map{"success": true, "message": "Updated", "hardware": hardware})

	return nil
}

// @Summary Delete hardware by ID
// @Description Delete hardware by ID
// @Tags Hardwares
// @Accept json
// @Produce json
// @Param id path int true "Hardware ID"
// @Success 200 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /hardwares/{id} [delete]
func HardwaresDelete(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		c.Status(http.StatusInternalServerError).JSON(&fiber.Map{"success": false, "message": "No ID entered"})
		return nil
	}

	hardwareOrig := &models.Hardware{}
	res := store.DB.First(hardwareOrig, id)
	if res.Error != nil {
		c.Status(http.StatusInternalServerError).JSON(&fiber.Map{"success": false, "message": "Record not found"})
		return nil
	}

	res = store.DB.Delete(&models.Hardware{}, id)
	if res.Error != nil {
		c.Status(http.StatusInternalServerError).JSON(&fiber.Map{"success": false, "message": res.Error.Error()})
		return nil
	}

	c.Status(http.StatusOK).JSON(&fiber.Map{"success": true, "message": "Deleted", "hardware ID": id})

	return nil
}
