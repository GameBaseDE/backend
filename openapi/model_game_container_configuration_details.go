/*
 * GameBase Communication API
 *
 * This is the REST API used as an communication layer between backend and frontend.
 *
 * API version: 2.1.0
 * Contact: gamebase@gahr.dev
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

// GameContainerConfigurationDetails - General details of server
type GameContainerConfigurationDetails struct {

	// Name of server that will displayed on one's Dashboard
	ServerName string `json:"serverName,omitempty"`

	// UUID of owner whom this server belongs to
	OwnerId string `json:"ownerId,omitempty"`

	// Short description of server which will be displayed on one's Dashboard
	Description string `json:"description,omitempty"`
}
