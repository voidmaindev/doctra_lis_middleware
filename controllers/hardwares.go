package controllers

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/voidmaindev/doctra_lis_middleware/models"
	"github.com/voidmaindev/doctra_lis_middleware/store"
)

func HardwaresGetAll(c *fiber.Ctx) error {
	hardwares := &[]models.Hardware{}
	res := store.DB.Find(hardwares)
	if res.Error != nil {
		c.Status(http.StatusInternalServerError).JSON(&fiber.Map{"success" : false, "message": res.Error.Error()})
		return nil
	}

	c.Status(http.StatusOK).JSON(&fiber.Map{"success" : true, "message": "All hardwares", "hardwares": hardwares})

	return nil
}

func HardwaresGetByID(c *fiber.Ctx) error {
	id := c.Params("id");
	if id == "" {
		c.Status(http.StatusInternalServerError).JSON(&fiber.Map{"success" : false, "message": "No ID entered"})
		return nil
	}
	
	hardware := &models.Hardware{}
	res := store.DB.First(hardware, id)
	if res.Error != nil {
		c.Status(http.StatusInternalServerError).JSON(&fiber.Map{"success" : false, "message": "Record not found"})
		return nil
	}

	c.Status(http.StatusOK).JSON(&fiber.Map{"success" : true, "message": "Hardware", "hardware": hardware})

	return nil
}

func HardwaresCreate(c *fiber.Ctx) error {
	hardware := &models.Hardware{}
	err := c.BodyParser(hardware)
	if err != nil {
		c.Status(http.StatusBadRequest).JSON(&fiber.Map{"success" : false, "message": err.Error()})
		return nil
	}

	res := store.DB.Create(hardware)
	if res.Error != nil {
		c.Status(http.StatusInternalServerError).JSON(&fiber.Map{"success" : false, "message": res.Error.Error()})
		return nil
	}

	c.Status(http.StatusOK).JSON(&fiber.Map{"success" : true, "message": "Created", "hardware": hardware})

	return nil
}

func HardwaresUpdate(c *fiber.Ctx) error {
	idStr := c.Params("id");
	if idStr == "" {
		c.Status(http.StatusInternalServerError).JSON(&fiber.Map{"success" : false, "message": "No ID entered"})
		return nil
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.Status(http.StatusBadRequest).JSON(&fiber.Map{"success" : false, "message": err.Error()})
		return nil
	}
	
	hardware := &models.Hardware{}
	err = c.BodyParser(hardware)
	if err != nil {
		c.Status(http.StatusBadRequest).JSON(&fiber.Map{"success" : false, "message": err.Error()})
		return nil
	}
	hardware.ID = uint(id)

	hardwareOrig := &models.Hardware{}
	res := store.DB.First(hardwareOrig, id)
	if res.Error != nil {
		c.Status(http.StatusInternalServerError).JSON(&fiber.Map{"success" : false, "message": "Record not found"})
		return nil
	}

	res = store.DB.Save(hardware)
	if res.Error != nil {
		c.Status(http.StatusInternalServerError).JSON(&fiber.Map{"success" : false, "message": res.Error.Error()})
		return nil
	}
	
	c.Status(http.StatusOK).JSON(&fiber.Map{"success" : true, "message": "Updated", "hardware": hardware})

	return nil
}

func HardwaresDelete(c *fiber.Ctx) error {
	id := c.Params("id");
	if id == "" {
		c.Status(http.StatusInternalServerError).JSON(&fiber.Map{"success" : false, "message": "No ID entered"})
		return nil
	}
	
	hardwareOrig := &models.Hardware{}
	res := store.DB.First(hardwareOrig, id)
	if res.Error != nil {
		c.Status(http.StatusInternalServerError).JSON(&fiber.Map{"success" : false, "message": "Record not found"})
		return nil
	}
	
	res = store.DB.Delete(&models.Hardware{}, id)
	if res.Error != nil {
		c.Status(http.StatusInternalServerError).JSON(&fiber.Map{"success" : false, "message": res.Error.Error()})
		return nil
	}
	
	c.Status(http.StatusOK).JSON(&fiber.Map{"success" : true, "message": "Deleted", "hardware ID": id})

	return nil
}