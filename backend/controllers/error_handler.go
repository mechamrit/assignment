package controllers

import (
	"github.com/go-playground/validator/v10"
)

func getErrorMessage(err error) string {
	if ve, ok := err.(validator.ValidationErrors); ok {
		fe := ve[0]
		switch fe.Field() {
		case "Username":
			if fe.Tag() == "required" {
				return "Username is required"
			}
			if fe.Tag() == "alphanum" {
				return "Username must be alphanumeric"
			}
			if fe.Tag() == "min" || fe.Tag() == "max" {
				return "Username must be between 3 and 30 characters"
			}
		case "Password":
			if fe.Tag() == "required" {
				return "Password is required"
			}
			if fe.Tag() == "min" {
				return "Password must be at least 8 characters long"
			}
		case "Role":
			return "Please select a valid role"
		case "Title":
			if fe.Tag() == "required" {
				return "Title is required"
			}
			if fe.Tag() == "min" || fe.Tag() == "max" {
				return "Title must be between 3 and 100 characters"
			}
		case "ProjectID":
			return "Valid Project ID is required"
		}
	}
	return "Invalid input data"
}
